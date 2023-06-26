package pomodoro

import (
	"context"

	"github.com/mum4k/termdash/cell"
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

func (w *widgets) update(timer []int, txtType, txtInfo, txtTimer string,
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