package pomodoro

import (
	"context"
	"errors"
	"time"
)

type Callback func(Pomodoro)

const (
	CatPomodoro = "Pomodoro"
	CatShortBreak = "ShortBreak"
	CatLongBreak = "LongBreak"
)

type Actions interface {
	Create(c Pomodoro) (string, error)
	Update(c Pomodoro) (string, error)
	GetById(id int) (Pomodoro, error)
	Delete(id int) (string, error)
}

const (
	StateNotStarted = iota
	StateRunning
	StatePaused
	StateDone
	StateCancelled
)

type Pomodoro struct {
	ID int
	StartTime time.Time
	Duration time.Duration
	TimeLeft time.Duration
	Category string
	State int
}

type Instance struct {
	repo Actions
}

var (
	ErrNoIntervals = errors.New("no intervals")
	ErrIntervalNotRunning = errors.New("interval not running")
	ErrIntervalCompleted = errors.New("interval is completed or cancelled")
	ErrInvalidState = errors.New("invalid state") 
	ErrInvalidID = errors.New("invalid id")
)

func Tick (ctx context.Context, id int, conf *Pomodoro, start, periodic, end Callback) error {
	
	ticker:= time.NewTicker(time.Second)
	defer ticker.Stop()
	i, err := Instance.repo.Create()
	if err != nil {
		return err
	}
	expire := time.After(conf.TimeLeft)
	start(i)

}