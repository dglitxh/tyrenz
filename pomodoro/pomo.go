package pomodoro

import (
	"errors"
	"time"
)

const (
	CatPomodoro = "Pomodoro"
	CatShortBreak = "ShortBreak"
	CatLongBreak = "LongBreak"
)

type Actions interface {
	Create() (string, error)
	Update() (string, error)
	GetById() (string, error)
	Delete() (string, error)
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
	ErrInvalidState = errors.New("invalid itate") 
	ErrInvalidID = errors.New("invalid id")
)

func Tick () {
	
}