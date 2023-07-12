package pomodoro

import (
	"context"
	"image"
	"time"
	"github.com/mum4k/termdash"
	"github.com/mum4k/termdash/terminal/tcell"
	"github.com/mum4k/termdash/terminal/terminalapi"
	"pragprog.com/rggo/interactiveTools/pomo/pomodoro"
)


type App struct {
	ctx context.Context
	controller *termdash.Controller
	redrawCh chan bool
	errorCh chan error
	term *tcell.Terminal
	size image.Point
}

func New(config *pomodoro.IntervalConfig) (*App, error) {
	ctx, cancel := context.WithCancel(context.Background())

	quitter := func(k *terminalapi.Keyboard) {
		if k.Key == 'q' || k.Key == 'Q' {
		cancel()
		}
	}
}