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