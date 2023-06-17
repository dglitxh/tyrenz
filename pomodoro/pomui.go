package pomodoro

import (
	"github.com/mum4k/termdash/widgets/donut"
	"github.com/mum4k/termdash/widgets/segmentdisplay"
	"github.com/mum4k/termdash/widgets/text"
)

type widgets struct {
	donTimer *donut.Donut
	disType *segmentdisplay.SegmentDisplay
	txtInfo *text.Texts
	txtTimer *text.Text
	updateDonTimer chan []int
	updateTxtInfo chan string
	updateTxtTimer chan string
	updateTxtType chan string
}