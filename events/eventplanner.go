package events

import (
	"context"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/libp2p/go-libp2p-kad-dht/internal"
)

const (
	// MAGIC: the EventPlanner will wake up at most once per hour, even if there
	// is nothing to do.
	UnlimitedSleepDuration = time.Hour
)

type EventPlanner struct {
	Clock clock.Clock

	NextAction *timedAction
}

type timedAction struct {
	action interface{}
	time   time.Time
	next   *timedAction
}

func NewEventPlanner(clk clock.Clock) *EventPlanner {
	return &EventPlanner{
		Clock: clk,
	}
}

func ScheduleAction(ctx context.Context, ep *EventPlanner, d time.Duration, a interface{}) {
	t := ep.Clock.Now().Add(d)
	if ep.NextAction == nil {
		ep.NextAction = &timedAction{action: a, time: t}
		return
	}

	curr := ep.NextAction
	if t.Before(curr.time) {
		ep.NextAction = &timedAction{action: a, time: t, next: curr}
		return
	}
	for curr.next != nil && t.After(curr.next.time) {
		curr = curr.next
	}
	curr.next = &timedAction{action: a, time: t, next: curr.next}
}

func RemoveAction(ctx context.Context, ep *EventPlanner, a interface{}) {
	if ep.NextAction == nil {
		return
	}

	curr := ep.NextAction
	if curr.action == a {
		ep.NextAction = curr.next
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

func RunOverdueActions(ctx context.Context, ep *EventPlanner) time.Time {
	ctx, span := internal.StartSpan(ctx, "events.RunOverdueActions")
	defer span.End()

	now := ep.Clock.Now()
	for ep.NextAction != nil && ep.NextAction.time.Before(now) {
		switch e := ep.NextAction.action.(type) {
		case func(context.Context):
			e(ctx)
		default:
		}

		ep.NextAction = ep.NextAction.next

		now = ep.Clock.Now()
	}

	if ep.NextAction == nil {
		span.AddEvent("no further actions to run")
		return time.Time{}
	}
	return ep.NextAction.time
}

func TimeUntilNextWakeUp(ep *EventPlanner) time.Duration {
	if ep.NextAction == nil {
		return UnlimitedSleepDuration
	}

	now := ep.Clock.Now()
	if ep.NextAction.time.Before(now) {
		return 0
	}
	return ep.NextAction.time.Sub(now)
}
