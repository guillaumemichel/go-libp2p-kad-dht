package simpleplanner

import (
	"context"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/libp2p/go-libp2p-kad-dht/events/action"
	"github.com/libp2p/go-libp2p-kad-dht/events/planner"
	"github.com/libp2p/go-libp2p-kad-dht/util"
)

type SimplePlanner struct {
	Clock clock.Clock

	NextAction *simpleTimedAction
}

type simpleTimedAction struct {
	action action.Action
	time   time.Time
	next   *simpleTimedAction
}

func NewSimplePlanner(clk clock.Clock) *SimplePlanner {
	return &SimplePlanner{
		Clock: clk,
	}
}

func (p *SimplePlanner) ScheduleAction(ctx context.Context, t time.Time, a action.Action) planner.PlannedAction {
	if p.NextAction == nil {
		p.NextAction = &simpleTimedAction{action: a, time: t}
		return p.NextAction
	}

	curr := p.NextAction
	if t.Before(curr.time) {
		p.NextAction = &simpleTimedAction{action: a, time: t, next: curr}
		return p.NextAction
	}
	for curr.next != nil && t.After(curr.next.time) {
		curr = curr.next
	}
	curr.next = &simpleTimedAction{action: a, time: t, next: curr.next}
	return curr.next
}

func (p *SimplePlanner) RemoveAction(ctx context.Context, pa planner.PlannedAction) {
	a, ok := pa.(*simpleTimedAction)
	if !ok {
		return
	}

	curr := p.NextAction
	if curr == nil {
		return
	}

	if curr == a {
		p.NextAction = curr.next
		return
	}
	for curr.next != nil {
		if curr.next == a {
			curr.next = curr.next.next
			return
		}
		curr = curr.next
	}
}

func (p *SimplePlanner) PopOverdueActions(ctx context.Context) []action.Action {
	var overdue []action.Action
	now := p.Clock.Now()
	curr := p.NextAction
	for curr != nil && (curr.time.Before(now) || curr.time == now) {
		overdue = append(overdue, curr.action)
		curr = curr.next
	}
	p.NextAction = curr
	return overdue
}

func (p *SimplePlanner) NextActionTime(context.Context) time.Time {
	if p.NextAction == nil {
		return util.MaxTime
	}
	return p.NextAction.time
}
