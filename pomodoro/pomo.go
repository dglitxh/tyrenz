package pomodoro

import (
	"context"
	"errors"
	"fmt"
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
	ErrNoIntervals = errors.New("No intervals")
	ErrIntervalNotRunning = errors.New("Interval not running")
	ErrIntervalCompleted = errors.New("Interval is completed or cancelled")
	ErrInvalidState = errors.New("Invalid State")
	ErrInvalidID = errors.New("Invalid ID")
)