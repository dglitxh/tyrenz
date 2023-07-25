package pomodoro

import (
	"context"
	"fmt"

	"github.com/mum4k/termdash/cell"
	"github.com/mum4k/termdash/widgets/button"
	"github.com/mum4k/termdash/widgets/donut"
	"github.com/mum4k/termdash/widgets/segmentdisplay"
	"github.com/mum4k/termdash/widgets/text"
)

type widgets struct {
	donTimer       *donut.Donut
	disType        *segmentdisplay.SegmentDisplay
	txtInfo        *text.Text
	txtTimer       *text.Text
	updateDonTimer chan []int
	updateTxtInfo  chan string
	updateTxtTimer chan string
	updateTxtType  chan string
}

type Buttons struct {
	start *button.Button
	pause *button.Button
}

func (w *widgets) Update(timer []int, txtType, txtInfo, txtTimer string,
	redrawCh chan<- bool) {
	if txtInfo != "" {
		w.updateTxtInfo <- txtInfo
	}
	if txtType != "" {
		w.updateTxtType <- txtType
	}
	if txtTimer != "" {
		w.updateTxtTimer <- txtTimer
	}
	if len(timer) > 0 {
		w.updateDonTimer <- timer
	}
	redrawCh <- true
}

func NewWidgets(ctx context.Context, errorCh chan<- error) (*widgets, error) {
	w := &widgets{}
	var err error
	w.updateDonTimer = make(chan []int)
	w.updateTxtType = make(chan string)
	w.updateTxtInfo = make(chan string)
	w.updateTxtTimer = make(chan string)
	w.donTimer, err = NewDonut(ctx, w.updateDonTimer, errorCh)
	if err != nil {
		return nil, err
	}
	w.disType, err = NewSegmentDisplay(ctx, w.updateTxtType, errorCh)
	if err != nil {
		return nil, err
	}
	w.txtInfo, err = NewText(ctx, w.updateTxtInfo, errorCh)
	if err != nil {
		return nil, err
	}
	w.txtTimer, err = NewText(ctx, w.updateTxtTimer, errorCh)
	if err != nil {
		return nil, err
	}
	return w, nil
}

func NewText(ctx context.Context, updateText <-chan string,
	errorCh chan<- error) (*text.Text, error) {
	txt, err := text.New()
	if err != nil {
		return nil, err
	}
	go func() {
		for {
			select {
			case t := <-updateText:
				txt.Reset()
				errorCh <- txt.Write(t)
			case <-ctx.Done():
				return
			}
		}
	}()
	return txt, nil
}

func NewSegmentDisplay(ctx context.Context, updateText <-chan string,
	errorCh chan<- error) (*segmentdisplay.SegmentDisplay, error) {
	sd, err := segmentdisplay.New()
	if err != nil {
		return nil, err
	}
	go func() {
		for {
			select {
			case t := <-updateText:
				if t == "" {
					t = " "
				}
				errorCh <- sd.Write([]*segmentdisplay.TextChunk{
					segmentdisplay.NewChunk(t),
				})
			case <-ctx.Done():
				return
			}
		}
	}()
	return sd, nil
}

func NewDonut(ctx context.Context, donUpdater <-chan []int,
	errorCh chan<- error) (*donut.Donut, error) {
	don, err := donut.New(
		donut.Clockwise(),
		donut.CellOpts(cell.FgColor(cell.ColorBlue)),
	)
	if err != nil {
		return nil, err
	}
	go func() {
		for {
			select {
			case d := <-donUpdater:
				if d[0] <= d[1] {
					errorCh <- don.Absolute(d[0], d[1])
				}
			case <-ctx.Done():
				return
			}
		}
	}()
	return don, nil
}

func (inst *Instance) NewButtonSet(ctx context.Context,
	w *widgets, redrawCh chan<- bool, errorCh chan<- error) (*Buttons, error) {
	startInterval := func() {
		_, err := inst.Action.GetById(inst.Conf.ID)
		fmt.Println(inst)
		errorCh <- err
		start := func(c Config) {
			message := "Take a break"
			if c.Category == CatPomodoro {
				message = "Focus on your task"
			}
			w.Update([]int{}, c.Category, message, "", redrawCh)
		}
		end := func(Config) {
			w.Update([]int{}, "", "Nothing running...", "", redrawCh)
		}

		periodic := func(c Config) {
			w.Update(
				[]int{int(c.Duration), int(c.TimeLeft)},
				"", "",
				fmt.Sprint(c.Duration-c.TimeLeft),
				redrawCh,
			)
		}

		errorCh <- inst.Start(ctx, start, periodic, end)

	}

	pauseInterval := func() {
		_, err := inst.Action.GetById(inst.Conf.ID)
		if err != nil {
			errorCh <- err
			return
		}
		if err := inst.Pause(); err != nil {
			if err == ErrIntervalNotRunning {
				return
			}
			errorCh <- err
			return
		}
		w.Update([]int{}, "", "Paused... press start to continue", "", redrawCh)
	}

	btStart, err := button.New("(s)tart", func() error {
		go startInterval()
		return nil
	},
		button.GlobalKey('s'),
		button.WidthFor("(p)ause"),
		button.Height(2),
	)

	btPause, err := button.New("(p)ause", func() error {
		go pauseInterval()
		return nil
	},
		button.FillColor(cell.ColorNumber(220)),
		button.GlobalKey('p'),
		button.Height(2),
	)
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	return &Buttons{btStart, btPause}, nil
}
