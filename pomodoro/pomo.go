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
	Update(c Config) (string, error)
	GetById(id int) (string, error)
	Delete(id int) (string, error)
}

const (
	StateNotStarted = iota
	StateRunning
	StatePaused
	StateDone
	StateCancelled
)

type Config struct {
	ID int64
	StartTime time.Time
	Duration time.Duration
	Category string
	State int
}

var (
	ErrNoIntervals = errors.New("no intervals")
	ErrIntervalNotRunning = errors.New("interval not running")
	ErrIntervalCompleted = errors.New("interval is completed or cancelled")
	ErrInvalidState = errors.New("invalid state") 
	ErrInvalidID = errors.New("invalid id")
)

func Tick (ctx context.Context, id int, conf *Config, start, periodic, end Callback) error {
	
	ticker:= time.NewTicker(time.Second)
	defer ticker.Stop()


}