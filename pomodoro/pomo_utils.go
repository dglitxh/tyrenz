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
	ticker:= time.NewTicker(time.Second)
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
			i, err := instance.Action.GetById(i.ID); if err != nil {
				helpers.Logger(err.Error())
				return err
			}
			if i.State == StatePaused {
				return nil
			}
			i.TimeElapsed += time.Second
			if c, err := instance.Action.Update(i); err != nil {
				return err
			}else {
				i = c
			}
			instance.Conf = i
			periodic(i)
		case <- expire:
			helpers.Logger("Expiring......")
			i.State = StateDone
			i, err := instance.Action.Update(i); if err != nil {
				return nil
			}
			instance.Conf = i
			end(i)
			return nil
		case <- ctx.Done():
			i.State = StateCancelled
			_, err := instance.Action.Update(i); if err != nil {
				helpers.Logger(err.Error(), "tick ctx done err")
				return err
			}
			
			return nil
		}
	}
}

func NewInstance(inst *Instance, cat string, pomodoro, longbrk, shortbrk int) *Instance {
	i := Config{
		Category: cat,
	}
	helpers.Logger(pomodoro, longbrk, shortbrk)
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
	
	i.State = StateNotStarted
	i.ID = len(inst.Action.Pomodoros )+1
	inst.Conf = i
	inst.Action.Create(inst.Conf)
	helpers.Logger("New Instance created.")
	helpers.Logger(fmt.Sprint(inst.Action.Pomodoros), inst.Conf)
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
			fallthrough
		case StatePaused:
			i.Conf.State = StateRunning
			if _, err := i.Action.Update(i.Conf); err != nil {
				return err
			}
			return Tick(ctx, i.Conf.ID, i, start, periodic, end)
		case StateDone:
			lastind := len(i.Action.Pomodoros)-1
			if i.Action.Pomodoros[lastind].Category != CatPomodoro {
				NewInstance(i, CatPomodoro, 0, 0, 0)
			}else {
				if i.Action.GetBreaks() > 2 {
					NewInstance(i, CatLongBreak, 0, 0, 0)
				}else{
					NewInstance(i, CatShortBreak, 0, 0, 0)
				}
			}
		    Start(ctx, i, start, periodic, end)
			return nil
		case StateCancelled:
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
		helpers.Logger(err.Error(), "Pause")
		return err
	}
	helpers.Logger("Timer paused")
	return nil
}