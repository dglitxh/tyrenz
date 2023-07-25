package pomodoro

import (
	"context"
	"errors"
	"fmt"
	"time"
)

type Callback func(Config)

const (
	CatPomodoro   = "Pomodoro"
	CatShortBreak = "ShortBreak"
	CatLongBreak  = "LongBreak"
)

type Actions interface {
	Create(c Config) (string, error)
	Update(c Config) error
	GetById(id int) (Config, error)
	Delete(id int) (string, error)
	GetCompleted(id int) []Config
}

const (
	StateNotStarted = iota
	StateRunning
	StatePaused
	StateDone
	StateCancelled
)

type Config struct {
	ID        int
	StartTime time.Time
	Duration  time.Duration
	TimeLeft  time.Duration
	Category  string
	State     int
}

type Instance struct {
	Conf   Config
	Action InMemStore
}

var (
	ErrNoIntervals        = errors.New("no intervals")
	ErrIntervalNotRunning = errors.New("interval not running")
	ErrIntervalCompleted  = errors.New("interval is completed or cancelled")
	ErrInvalidState       = errors.New("invalid state")
	ErrInvalidID          = errors.New("invalid id")
)

func (instance *Instance) Tick(ctx context.Context, id int, start, periodic, end Callback) error {

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	i := instance.Conf
	expire := time.After(i.TimeLeft)
	start(i)
	for {
		select {
		case <-ticker.C:
			if i.State == StatePaused {
				return nil
			}
			i.TimeLeft -= time.Second
			if err := instance.Action.Update(i); err != nil {
				return nil
			}
			periodic(i)
		case <-expire:
			i.State = StateDone
			end(i)
			return instance.Action.Update(i)
		case <-ctx.Done():
			i.State = StateCancelled
			return instance.Action.Update(i)

		}
	}

}

func (inst *Instance) NewInstance(pomodoro, longbrk, shortbrk int) *Instance {
	i := Config{
		Category: CatPomodoro,
	}
	i.State = StateNotStarted
	switch i.Category {
	case CatPomodoro:
		if pomodoro < 1 {
			i.Duration = time.Minute * 25
		} else {
			i.Duration = time.Minute * time.Duration(pomodoro)
		}
	case CatShortBreak:
		if shortbrk < 1 {
			i.Duration = time.Minute * 5
		} else {
			i.Duration = time.Minute * time.Duration(shortbrk)
		}
	case CatLongBreak:
		if longbrk < 1 {
			i.Duration = time.Minute * 15
		} else {
			i.Duration = time.Minute * time.Duration(longbrk)
		}
	}
	i.TimeLeft = i.Duration
	i.ID = len(inst.Action.Pomodoros) + 1
	inst.Conf = i
	inst.Action.Create(i)

	return inst
}

func (i *Instance) Start(ctx context.Context,
	start, periodic, end Callback) error {
	switch i.Conf.State {
	case StateRunning:
		return nil
	case StateNotStarted:
		i.Conf.StartTime = time.Now()
		i.Conf.State = StateRunning
		i.Tick(ctx, i.Conf.ID, start, periodic, end)
		fallthrough
	case StatePaused:
		i.Conf.State = StateRunning
		if err := i.Action.Update(i.Conf); err != nil {
			return err
		}
		return i.Tick(ctx, i.Conf.ID, start, periodic, end)
	case StateCancelled, StateDone:
		return fmt.Errorf("%w: Cannot start", ErrIntervalCompleted)
	default:
		return fmt.Errorf("%w: %d", ErrInvalidState, i.Conf.State)
	}
}

func (i *Instance) Pause() error {
	if i.Conf.State != StateRunning {
		return ErrIntervalNotRunning
	}
	i.Conf.State = StatePaused
	return i.Action.Update(i.Conf)
}
