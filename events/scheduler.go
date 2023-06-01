package events

import (
	"context"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/libp2p/go-libp2p-kad-dht/internal"
)

const (
	// MAGIC: the scheduler will wake up at most once per hour, even if there
	// is nothing to do.
	UnlimitedSleepDuration = time.Hour
)

type Scheduler struct {
	Clock clock.Clock

	NextAction *timedAction
}

type timedAction struct {
	action interface{}
	time   time.Time
	next   *timedAction
}

func NewScheduler(clk clock.Clock) *Scheduler {
	return &Scheduler{
		Clock: clk,
	}
}

func ScheduleAction(ctx context.Context, sched *Scheduler, d time.Duration, a interface{}) {
	t := sched.Clock.Now().Add(d)
	if sched.NextAction == nil {
		sched.NextAction = &timedAction{action: a, time: t}
		return
	}

	curr := sched.NextAction
	if t.Before(curr.time) {
		sched.NextAction = &timedAction{action: a, time: t, next: curr}
		return
	}
	for curr.next != nil && t.After(curr.next.time) {
		curr = curr.next
	}
	curr.next = &timedAction{action: a, time: t, next: curr.next}
}

func RemoveAction(ctx context.Context, sched *Scheduler, a interface{}) {
	if sched.NextAction == nil {
		return
	}

	curr := sched.NextAction
	if curr.action == a {
		sched.NextAction = curr.next
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

func RunOverdueActions(ctx context.Context, sched *Scheduler) time.Time {
	ctx, span := internal.StartSpan(ctx, "events.RunOverdueActions")
	defer span.End()

	now := sched.Clock.Now()
	for sched.NextAction != nil && sched.NextAction.time.Before(now) {
		switch e := sched.NextAction.action.(type) {
		case func(context.Context):
			e(ctx)
		default:
		}

		sched.NextAction = sched.NextAction.next

		now = sched.Clock.Now()
	}

	if sched.NextAction == nil {
		span.AddEvent("no further actions to run")
		return time.Time{}
	}
	return sched.NextAction.time
}

func TimeUntilNextWakeUp(sched *Scheduler) time.Duration {
	if sched.NextAction == nil {
		return UnlimitedSleepDuration
	}

	now := sched.Clock.Now()
	if sched.NextAction.time.Before(now) {
		return 0
	}
	return sched.NextAction.time.Sub(now)
}
