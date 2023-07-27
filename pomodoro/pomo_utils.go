package pomodoro

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/dglitxh/tyrenz/helpers"
)

type Callback func(Config)

const (
	CatPomodoro = "Pomodoro"
	CatShortBreak = "ShortBreak"
	CatLongBreak = "LongBreak"
)

type Actions interface {
	Create(c Config) (int, error)
	Update(c Config) (error)
	GetById(id int) (Config, error)
	Delete(id int) (string, error)
	GetCompleted() ([]Config)
}

const (
	StateNotStarted = iota
	StateRunning
	StatePaused
	StateDone
	StateCancelled
)

type Config struct {
	ID int
	StartTime time.Time
	Duration time.Duration
	TimeLeft time.Duration
	Category string
	State int
}

type Instance struct {
	Conf Config
	Action  InMemStore
}

var (
	ErrNoIntervals = errors.New("no intervals")
	ErrIntervalNotRunning = errors.New("interval not running")
	ErrIntervalCompleted = errors.New("interval is completed or cancelled")
	ErrInvalidState = errors.New("invalid state") 
	ErrInvalidID = errors.New("invalid id")
)

func Tick (ctx context.Context, id int, instance *Instance, start, periodic, end Callback) error {
	ticker:= time.NewTicker(time.Second)
	defer ticker.Stop()
	i, err := instance.Action.GetById(id)	
	if err != nil {
		helpers.Logger(err.Error(), "tick")
		return err
	}
	helpers.Logger(strconv.Itoa(i.State), "state")
	if i.State == StateNotStarted {
		helpers.Logger("State Defined")
	}
	expire := time.After(instance.Conf.TimeLeft)
	start(i)
	for {
		select {
		case <- ticker.C:
			if i.State == StatePaused {
				return nil
			}
			i.TimeLeft -= time.Second
			if c, err := instance.Action.Update(i); err != nil {
				return nil
			}else {
				i = c
			}
			periodic(i)
		case <- expire:
			i.State = StateDone
			end(i)
			_, err := instance.Action.Update(i); if err != nil {
				return nil
			}

		case <- ctx.Done():
			i.State = StateCancelled
			_, err := instance.Action.Update(i); if err != nil {
				return nil
			}
			
		}
	}
}

func NewInstance(inst *Instance, pomodoro, longbrk, shortbrk int) *Instance {
	i := Config{
		Category: CatPomodoro,
	}
	switch i.Category {
		case CatPomodoro:
			if pomodoro < 1 {
				i.Duration = time.Minute * 25
			}else {
				i.Duration = time.Minute * time.Duration(pomodoro)
			}
		case CatShortBreak:
			if shortbrk  < 1 {
				i.Duration = time.Minute * 5
			}else {
				i.Duration = time.Minute * time.Duration(shortbrk)
			}
		case CatLongBreak:
			if longbrk  < 1 {
				i.Duration = time.Minute * 15
			}else {
				i.Duration = time.Minute * time.Duration(longbrk)
			}
	}
	
	i.Category = CatPomodoro
	i.ID = len(inst.Action.Pomodoros )+1
	inst.Conf = i
	inst.Action.Create(i)
	helpers.Logger(string(inst.Conf.Category))
	return inst
}

func Start(ctx context.Context, i *Instance,
	start, periodic, end Callback) error {
	switch i.Conf.State {
		case StateRunning:
			return nil
		case StateNotStarted:
			i.Conf.StartTime = time.Now()
			fallthrough
		case StatePaused:
			i.Conf.State = StateRunning
			if _, err := i.Action.Update(i.Conf); err != nil {
				return err
			}
			return Tick(ctx, i.Conf.ID, i, start, periodic, end)
		case StateCancelled, StateDone:
				return fmt.Errorf("%w: Cannot start", ErrIntervalCompleted)
		default:
			return fmt.Errorf("%w: %d", ErrInvalidState, i.Conf.State)
		}
}


func Pause(i *Instance) error {
	if i.Conf.State != StateRunning {
		return ErrIntervalNotRunning
	}
	i.Conf.State = StatePaused
	if _, err := i.Action.Update(i.Conf); err != nil {
		return err
	}
	return nil
}