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
	TimeElapsed time.Duration
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
	ticker:= time.NewTicker(instance.Conf.Duration)
	defer ticker.Stop()
	i, err := instance.Action.GetById(id)	
	if err != nil {
		helpers.Logger(err.Error(), "tick")
		return err
	}
	helpers.Logger(strconv.Itoa(i.State), "stateddd")
	if i.State == StateNotStarted {
		helpers.Logger("State Defined")
	}
	expire := time.After(i.Duration-i.TimeElapsed)
	start(i)
	for {
		select {
		case <- ticker.C:
			if i.State == StatePaused {
				continue
			}
			i.TimeElapsed += time.Second
			if c, err := instance.Action.Update(i); err != nil {
				return err
			}else {
				i = c
			}
			instance.Conf = i
			helpers.Logger(i.Duration)
			periodic(i)
		case <- expire:
			helpers.Logger("Expiring......")
			i.State = StateDone
			_, err := instance.Action.Update(i); if err != nil {
				return nil
			}
			instance.Conf = i
			end(i)

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
	helpers.Logger(pomodoro, longbrk, shortbrk)
	switch i.Category {
		case CatPomodoro:
			if pomodoro < 1 {
				i.Duration = time.Minute * 25
			}else {
				i.Duration = time.Duration(pomodoro)
			}
		case CatShortBreak:
			if shortbrk  < 1 {
				i.Duration = time.Minute * 5
			}else {
				i.Duration = time.Duration(shortbrk)
			}
		case CatLongBreak:
			if longbrk  < 1 {
				i.Duration = time.Minute * 15
			}else {
				i.Duration = time.Duration(longbrk)
			}
	}
	
	i.Category = CatPomodoro
	i.State = StateNotStarted
	i.ID = len(inst.Action.Pomodoros )+1
	inst.Conf = i
	inst.Action.Create(i)
	helpers.Logger(fmt.Sprint(inst.Conf.Duration), "duration my boiiiii")
	return inst
}

func Start(ctx context.Context, i *Instance,
	start, periodic, end Callback) error {
	switch i.Conf.State {
		case StateRunning:
			return nil
		case StateNotStarted:
			i.Conf.StartTime = time.Now()
			i.Conf.State = StateRunning
			if _, err := i.Action.Update(i.Conf); err != nil {
				return err
			}
			Tick(ctx, i.Conf.ID, i, start, periodic, end)
		case StatePaused:
			i.Conf.State = StateRunning
			if _, err := i.Action.Update(i.Conf); err != nil {
				helpers.Logger("default state for start error please help me.")
				return err
			}
			return Tick(ctx, i.Conf.ID, i, start, periodic, end)
		case StateCancelled, StateDone:
			    helpers.Logger(ErrIntervalCompleted)
				return fmt.Errorf("%w: Cannot start", ErrIntervalCompleted)
		default:
			helpers.Logger("default state for start error please help me.")
			return fmt.Errorf("%w: %d", ErrInvalidState, i.Conf.State)
		}	
		return nil
}


func Pause(i *Instance) error {
	if i.Conf.State != StateRunning {
		helpers.Logger(i.Conf.State, "Pause: not started")
		return ErrIntervalNotRunning
	}
	i.Conf.State = StatePaused
	if _, err := i.Action.Update(i.Conf); err != nil {
		helpers.Logger(err.Error(), "Pause")
		return err
	}
	helpers.Logger("Timer paused")
	return nil
}