package forward

import (
	"fmt"
	"strings"
	"sync"

	"github.com/fatih/color"
	pw "github.com/jedib0t/go-pretty/v6/progress"
	"github.com/mattn/go-runewidth"

	"github.com/iyear/tdl/core/forwarder"
	"github.com/iyear/tdl/pkg/prog"
	"github.com/iyear/tdl/pkg/utils"
)

type progress struct {
	pw       pw.Writer
	trackers *sync.Map // map[tuple]*pw.Tracker
	elemName *sync.Map // map[int64]string
}

type tuple struct {
	from int64
	msg  int
	to   int64
}

func newProgress(p pw.Writer) *progress {
	return &progress{
		pw:       p,
		trackers: &sync.Map{},
		elemName: &sync.Map{},
	}
}

func (p *progress) OnAdd(elem forwarder.Elem) {
	tracker := prog.AppendTracker(p.pw, pw.FormatNumber, p.processMessage(elem, false), 1)
	p.trackers.Store(p.tuple(elem), tracker)
}

func (p *progress) OnClone(elem forwarder.Elem, state forwarder.ProgressState) {
	val, ok := p.trackers.Load(p.tuple(elem))
	if !ok {
		return
	}
	tracker := val.(*pw.Tracker)

	// display re-upload transfer info
	tracker.Units.Formatter = utils.Byte.FormatBinaryBytes
	tracker.UpdateMessage(p.processMessage(elem, true))
	tracker.UpdateTotal(state.Total)
	tracker.SetValue(state.Done)
}

func (p *progress) OnDone(elem forwarder.Elem, err error) {
	val, ok := p.trackers.Load(p.tuple(elem))
	if !ok {
		return
	}
	tracker := val.(*pw.Tracker)

	if err != nil {
		p.pw.Log(color.RedString("%s error: %s", p.metaString(elem), err.Error()))
		tracker.MarkAsErrored()
		return
	}

	if tracker.Total == 1 {
		tracker.Increment(1)
	}
	tracker.MarkAsDone()
}

func (p *progress) tuple(elem forwarder.Elem) tuple {
	return tuple{
		from: elem.From().ID(),
		msg:  elem.Msg().ID,
		to:   elem.To().ID(),
	}
}

func (p *progress) processMessage(elem forwarder.Elem, clone bool) string {
	b := &strings.Builder{}

	b.WriteString(p.metaString(elem))
	if clone {
		b.WriteString(" [clone]")
	}

	return b.String()
}

func (p *progress) metaString(elem forwarder.Elem) string {
	// TODO(iyear): better responsive name
	fromID := elem.From().ID()
	fromName, _ := p.elemName.LoadOrStore(fromID, runewidth.Truncate(elem.From().VisibleName(), 15, "..."))

	toID := elem.To().ID()
	toName, _ := p.elemName.LoadOrStore(toID, runewidth.Truncate(elem.To().VisibleName(), 15, "..."))

	return fmt.Sprintf("%s(%d):%d -> %s(%d)",
		fromName.(string),
		fromID,
		elem.Msg().ID,
		toName.(string),
		toID)
}
