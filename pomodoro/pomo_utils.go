package pomodoro

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/dglitxh/tyrenz/helpers"
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
	Duration time.Duration
	TimeElapsed time.Duration
	Category string
	State int
}

var fn *os.File = helpers.CreateLogFile(".pomologs")

type UserSpecs struct {
	LongBreak time.Duration
	ShortBreak time.Duration
	Interval time.Duration
}

type Instance struct {

	Conf Config
	Action  InMemStore
	Specs   UserSpecs
}

var (
	ErrNoIntervals        = errors.New("no intervals")
	ErrIntervalNotRunning = errors.New("interval not running")
	ErrIntervalCompleted  = errors.New("interval is completed or cancelled")
	ErrInvalidState       = errors.New("invalid state")
	ErrInvalidID          = errors.New("invalid id")
)

func Tick (ctx context.Context, id int, instance *Instance, start, periodic, end Callback) error {
	ticker:= time.NewTicker(time.Second)
	defer ticker.Stop()
	i, err := instance.Action.GetById(id)	
	if err != nil {
		helpers.Logger(fn, err.Error(), "tick")
		return err
	}
	if i.State == StateNotStarted {
		helpers.Logger(fn, "State Defined")
	}
	expire := time.After(i.Duration-i.TimeElapsed)
	start(i)
	for {
		select {
		case <- ticker.C:
			i, err := instance.Action.GetById(i.ID); if err != nil {
				helpers.Logger(fn, err.Error())
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
			helpers.Logger(fn, "Interval done.")
			i.State = StateDone
			i, err := instance.Action.Update(i); if err != nil {
				return err
			}
			instance.Conf = i
			end(i)
			return nil
		case <- ctx.Done():
			i.State = StateCancelled
			_, err := instance.Action.Update(i); if err != nil {
				helpers.Logger(fn, err.Error(), "tick ctx done err")
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
	helpers.Logger(fn, "pomo: "+fmt.Sprint(pomodoro),
	 		"long: "+fmt.Sprint(longbrk),
	  		"short:"+fmt.Sprint(shortbrk))
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
	helpers.Logger(fn, "New Instance created.")
	helpers.Logger(fn, fmt.Sprint(inst.Action.Pomodoros), inst.Conf)
	return inst
}

func (i *Instance) Start(ctx context.Context,
	start, periodic, end Callback) error {
	spc := i.Specs
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
				NewInstance(i, CatPomodoro, int(spc.Interval), int(spc.LongBreak), int(spc.LongBreak))
			}else {
				if i.Action.GetBreaks() > 2 {
					NewInstance(i, CatLongBreak, int(spc.Interval), int(spc.LongBreak), int(spc.LongBreak))
				}else{
					NewInstance(i, CatShortBreak, int(spc.Interval), int(spc.LongBreak), int(spc.LongBreak))
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
		helpers.Logger(fn, err.Error(), "Pause")
		return err
	}
	helpers.Logger(fn, "Timer paused")
	return nil
}

