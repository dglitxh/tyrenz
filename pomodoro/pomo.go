package pomodoro

import (
	"context"
	"errors"
	"time"
)

type Callback func(Config)

const (
	CatPomodoro = "Pomodoro"
	CatShortBreak = "ShortBreak"
	CatLongBreak = "LongBreak"
)

type Actions interface {
	Create(c Config) (string, error)
	Update(c Config) (error)
	GetById(id int) (Config, error)
	Delete(id int) (string, error)
	GetCompleted(id int) ([]Config)
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
	conf Config
	action   Actions
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
	i, err := instance.action.GetById(id)
	if err != nil {
		return err
	}
	expire := time.After(instance.conf.TimeLeft)
	start(i)
	for {
		select {
		case <- ticker.C:
			if i.State == StatePaused {
				return nil
			}
			i.TimeLeft -= time.Second
			if err := instance.action.Update(i); err != nil {
				return nil
			}
			periodic(i)
		case <- expire:
			i.State = StateDone
			end(i)
			return instance.action.Update(i)
		case <- ctx.Done():
			i.State = StateCancelled
			return instance.action.Update(i)
			
		}
	}

}

func NewInstance(inst Instance, pomodoro, longbrk, shortbrk int) *Instance {
	i := &Instance{
		action: inst.action,
		conf: inst.conf,
	}
	switch i.conf.Category {
		case CatPomodoro:
			if pomodoro < 1 {
				i.conf.Duration = time.Minute * 25
			}else {
				i.conf.Duration = time.Minute * time.Duration(pomodoro)
			}
		case CatShortBreak:
			if shortbrk  < 1 {
				i.conf.Duration = time.Minute * 5
			}else {
				i.conf.Duration = time.Minute * time.Duration(shortbrk)
			}
		case CatLongBreak:
			if longbrk  < 1 {
				i.conf.Duration = time.Minute * 15
			}else {
				i.conf.Duration = time.Minute * time.Duration(longbrk)
			}
	}

	return i
}