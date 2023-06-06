package simpleplanner

import (
	"context"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/libp2p/go-libp2p-kad-dht/events"
	"github.com/libp2p/go-libp2p-kad-dht/internal"
)

const (
	// MAGIC: the SimplePlanner will wake up at most once per hour, even if there
	// is nothing to do.
	UnlimitedSleepDuration = time.Hour
)

type SimplePlanner struct {
	Clock clock.Clock

	NextAction *timedAction
}

type timedAction struct {
	action events.Action
	time   time.Time
	next   *timedAction
}

func NewSimplePlanner(clk clock.Clock) *SimplePlanner {
	return &SimplePlanner{
		Clock: clk,
	}
}

func (p *SimplePlanner) ScheduleAction(ctx context.Context, t time.Time, a events.Action) {

	if p.NextAction == nil {
		p.NextAction = &timedAction{action: a, time: t}
		return
	}

	curr := p.NextAction
	if t.Before(curr.time) {
		p.NextAction = &timedAction{action: a, time: t, next: curr}
		return
	}
	for curr.next != nil && t.After(curr.next.time) {
		curr = curr.next
	}
	curr.next = &timedAction{action: a, time: t, next: curr.next}
}

func (p *SimplePlanner) RemoveAction(ctx context.Context, a events.Action) {
	if p.NextAction == nil {
		return
	}

	curr := p.NextAction
	if curr.action == a {
		p.NextAction = curr.next
		return
	}
	for curr.next != nil {
		if curr.next.action == a {
			curr.next = curr.next.next
			return
		}
		curr = curr.next
	}
}

func (p *SimplePlanner) PopOverdueActions(ctx context.Context) []events.Action {
	_, span := internal.StartSpan(ctx, "events.OverdueActions")
	defer span.End()

	var overdue []events.Action
	now := p.Clock.Now()
	curr := p.NextAction
	for curr != nil && curr.time.Before(now) {
		overdue = append(overdue, curr.action)
		curr = curr.next
	}
	p.NextAction = curr
	return overdue
}

/*
func RunOverdueActions(ctx context.Context, p *SimplePlanner) time.Time {
	ctx, span := internal.StartSpan(ctx, "events.RunOverdueActions")
	defer span.End()

	now := p.Clock.Now()
	for p.NextAction != nil && p.NextAction.time.Before(now) {
		switch e := p.NextAction.action.(type) {
		case func(context.Context):
			e(ctx)
		default:
		}

		p.NextAction = p.NextAction.next

		now = p.Clock.Now()
	}

	if p.NextAction == nil {
		span.AddEvent("no further actions to run")
		return time.Time{}
	}
	return p.NextAction.time
}

func TimeUntilNextWakeUp(p *SimplePlanner) time.Duration {
	if p.NextAction == nil {
		return UnlimitedSleepDuration
	}

	now := p.Clock.Now()
	if p.NextAction.time.Before(now) {
		return 0
	}
	return p.NextAction.time.Sub(now)
}
*/
