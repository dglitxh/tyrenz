package pomodoro

import (
	"context"
	"image"
	"time"

	"github.com/dglitxh/tyrenz/helpers"
	"github.com/mum4k/termdash"
	"github.com/mum4k/termdash/terminal/tcell"
	"github.com/mum4k/termdash/terminal/terminalapi"
)


type App struct {
	ctx context.Context
	controller *termdash.Controller
	redrawCh chan bool
	errorCh chan error
	term *tcell.Terminal
	size image.Point
}

func (inst *Instance) New() (*App, error) {
	ctx, cancel := context.WithCancel(context.Background())

	quitter := func(k *terminalapi.Keyboard) {
		if k.Key == 'q' || k.Key == 'Q' {
			cancel()
		}
	}

	

	redrawCh := make(chan bool)
	errorCh := make(chan error)

	w, err := NewWidgets(ctx, errorCh)
		if err != nil {
			helpers.Logger(fn, "Error @ new widget")
			return nil, err
		}
	b, err := inst.NewButtonSet(ctx, w, redrawCh, errorCh)
		if err != nil {
			helpers.Logger(fn, "Error @j newbutton set")
			return nil, err
		}

	term, err := tcell.New()
	if err != nil {
		return nil, err
	}
	c, err := NewGrid(b, w, term)
	if err != nil {
		return nil, err
	}

	controller, err := termdash.NewController(term, c,
	termdash.KeyboardSubscriber(quitter))
	if err != nil {
		return nil, err
	}

	return &App{
		ctx: ctx,
		controller: controller,
		redrawCh: redrawCh,
		errorCh: errorCh,
		term: term,
	}, nil
}

func (a *App) resize() error {
	if a.size.Eq(a.term.Size()) {
		return nil
	}
	a.size = a.term.Size()
	if err := a.term.Clear(); err != nil {
		return err
	}
	return a.controller.Redraw()
}

func (a *App) Run() error {
	defer a.term.Close()
	defer a.controller.Close()
	ticker := time.NewTicker(2*time.Second)
	defer ticker.Stop()
	for {
		select {
			case <-a.redrawCh:
			if err := a.controller.Redraw(); err != nil {
				helpers.Logger(fn, err.Error(), ":  app redraw")
				return err
			}
			case <-ticker.C:
				if err := a.resize(); err != nil {
					return err
				}
			case err := <-a.errorCh:
				if err != nil {
					helpers.Logger(fn, err.Error(), ":  app runner")
					return err
				}
			case <-a.ctx.Done():
				helpers.Logger(fn, "done @ app")
				return nil
			
		}
	}
}