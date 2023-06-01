package events

import (
	"context"
	"fmt"

	eq "github.com/libp2p/go-libp2p-kad-dht/events/eventqueue"
)

func Run(ctx context.Context, sched *Scheduler, queue eq.EventQueue) {
	alarm := RunOverdueActions(ctx, sched)
	timer := sched.Clock.Timer(sched.Clock.Until(alarm))

	newsChan := eq.NewsChan(queue)
	for {
		select {
		case <-ctx.Done():
			return
		case <-timer.C:
			alarm = RunOverdueActions(ctx, sched)
			timer = sched.Clock.Timer(sched.Clock.Until(alarm))
		case <-newsChan:
			event := eq.Dequeue(queue)
			switch e := event.(type) {
			case func(context.Context):
				e(ctx)
			case func():
				e()
			default:
				fmt.Println("Unknown event type")
			}
			// TODO: if new events have been scheduled before the end of timer
			// by the handled event, the timer should ring earlier
		}
	}
}

func EmptyQueue(ctx context.Context, q eq.EventQueue) {
	for !eq.Empty(q) {
		event := q.Dequeue()

		switch e := event.(type) {
		case func(context.Context):
			e(ctx)
		case func():
			e()
		default:
			fmt.Println("Unknown event type")
		}
	}
}
