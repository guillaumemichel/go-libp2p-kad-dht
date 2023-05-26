package events

import (
	"context"
	"fmt"
)

func Run(ctx context.Context, sched *Scheduler, newEvents chan interface{}) {
	alarm := RunOverdueActions(ctx, sched)
	timer := sched.Clock.Timer(sched.Clock.Until(alarm))
	for {
		select {
		case <-ctx.Done():
			return
		case <-timer.C:
			alarm = RunOverdueActions(ctx, sched)
			timer = sched.Clock.Timer(sched.Clock.Until(alarm))
		case event := <-newEvents:
			switch e := event.(type) {
			case func(context.Context):
				e(ctx)
			default:
				fmt.Println("Unknown event type")
			}
			// TODO: if new events have been scheduled before the end of timer
			// by the handled event, the timer should ring earlier
		}
	}
}

func EmptyQueue(ctx context.Context, q *Queue) {
	for !q.Empty() {
		event := q.Dequeue()

		switch e := event.(type) {
		case func(context.Context):
			e(ctx)
		default:
			fmt.Println("Unknown event type")
		}
	}
}
