package scheduler

import (
	"context"
	"time"

	"github.com/libp2p/go-libp2p-kad-dht/events"
)

type Scheduler interface {
	Now() time.Time

	EnqueueAction(context.Context, events.Action)
	ScheduleAction(context.Context, time.Time, events.Action)

	RunOne(context.Context)
}

func ScheduleActionIn(ctx context.Context, s Scheduler, d time.Duration, a events.Action) {
	if d <= 0 {
		s.EnqueueAction(ctx, a)
	} else {
		s.ScheduleAction(ctx, s.Now().Add(d), a)
	}
}
