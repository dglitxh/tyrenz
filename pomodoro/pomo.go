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