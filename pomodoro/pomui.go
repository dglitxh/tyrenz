package pomodoro

import (
	"context"
	"fmt"
	"time"

	"github.com/dglitxh/tyrenz/helpers"
	"github.com/mum4k/termdash/cell"
	"github.com/mum4k/termdash/widgets/button"
	"github.com/mum4k/termdash/widgets/donut"
	"github.com/mum4k/termdash/widgets/segmentdisplay"
	"github.com/mum4k/termdash/widgets/text"
)

type widgets struct {
	donTimer *donut.Donut
	disType *segmentdisplay.SegmentDisplay
	txtInfo *text.Text
	txtTimer *text.Text
	updateDonTimer chan []int
	updateTxtInfo chan string
	updateTxtTimer chan string
	updateTxtType chan string
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
	dftext :=  "Press start to initiate next event..."
	var err error
	w.updateDonTimer = make(chan []int)
	w.updateTxtType = make(chan string)
	w.updateTxtInfo = make(chan string)
	w.updateTxtTimer = make(chan string)
	w.donTimer, err = NewDonut(ctx, w.updateDonTimer, errorCh)
	if err != nil {
		helpers.Logger("Donut creation error")
		return nil, err
	}
	w.disType, err = NewSegmentDisplay(ctx, w.updateTxtType, errorCh)
	if err != nil {
		helpers.Logger("Seg disp err, w.distype ui")
		return nil, err
	}
	w.txtInfo, err = NewText(ctx, dftext, w.updateTxtInfo, errorCh)
	if err != nil {
		helpers.Logger("new text error")
		return nil, err
	}
	w.txtTimer, err = NewText(ctx, " ", w.updateTxtTimer, errorCh)
	if err != nil {
		helpers.Logger("new text error")
		return nil, err
	}
	return w, nil
}

func NewText(ctx context.Context, dftext string, updateText <-chan string,
	errorCh chan<- error) (*text.Text, error) {
	txt, err := text.New()
	if err := txt.Write(dftext); err != nil {
		return nil, err
	}
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
	sd.Write([]*segmentdisplay.TextChunk{
		segmentdisplay.NewChunk("Welcome"), })
	if err != nil {
		helpers.Logger("Seg Disp error")
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
		donut.ShowTextProgress(),
		donut.CellOpts(cell.FgColor(cell.ColorGreen)),
		)
		if err != nil {
			helpers.Logger("donut error : NewDonut")
			return nil, err
		}
		go func() {
			for {
				select {
					case d := <-donUpdater:
					if d[0] >= d[1] {
						err := don.Absolute(d[1], d[0])
						if err != nil {
							helpers.Logger(err.Error(), "donut err still")
							errorCh <- err
						}
						
					}
					case <-ctx.Done():
						helpers.Logger(" Donut done still!")
						return
					}
				}	
		}()
		return don, nil
}

func NewButtonSet(ctx context.Context, config *Instance,
	w *widgets, redrawCh chan<- bool, errorCh chan<- error) (*Buttons, error) {
	startInterval := func() {
		_, err := config.Action.GetById(config.Conf.ID)
		if err != nil {
			errorCh <- err
			return
		}
		start := func(i Config) {
		i.State = StateRunning
		config.Conf.State = i.State
		message := "Take a break"
		if i.Category == CatPomodoro {
			message = "Focus on your task"
		}
		helpers.Logger("Timer started...")
	w.Update([]int{}, i.Category, message, "", redrawCh)
	}
	end := func(i Config) {
		w.Update([]int{}, "Idle", "Press start button to intiate next event", "", redrawCh)
	}

	periodic := func(i Config) {
		w.Update(
		[]int{int(i.Duration/time.Minute), int(i.TimeElapsed/time.Minute)},
		"", "",
		fmt.Sprint(i.Duration-i.TimeElapsed)+" left",
		redrawCh,
		)
	}
		errorCh <- Start(ctx, config, start, periodic, end)
	}

	

	pauseInterval := func() {
		_, err := config.Action.GetById(config.Conf.ID)
		if err != nil {
			errorCh <- err
			return
		}
		if err := Pause(config); err != nil {
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
	if err != nil {
			return nil, err
	}

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
	return &Buttons{btStart, btPause}, nil
}
