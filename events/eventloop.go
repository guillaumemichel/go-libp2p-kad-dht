package events

import (
	"context"
	"fmt"

	eq "github.com/libp2p/go-libp2p-kad-dht/events/eventqueue"
	"github.com/libp2p/go-libp2p-kad-dht/internal"
)

func RunLoop(ctx context.Context, ep *EventPlanner, queue eq.EventQueue) {
	alarm := RunOverdueActions(ctx, ep)
	timer := ep.Clock.Timer(ep.Clock.Until(alarm))

	newsChan := eq.NewsChan(queue)
	for {
		EmptyQueue(ctx, queue)

		select {
		case <-ctx.Done():
			return
		case <-timer.C:
			alarm = RunOverdueActions(ctx, ep)
			timer = ep.Clock.Timer(ep.Clock.Until(alarm))
		case <-newsChan:
			event := eq.Dequeue(queue)
			RunEvent(ctx, event)
			// TODO: if new events have been scheduled before the end of timer
			// by the handled event, the timer should ring earlier
		}
	}
}

func EmptyQueue(ctx context.Context, q eq.EventQueue) {
	for !eq.Empty(q) {
		event := q.Dequeue()
		RunEvent(ctx, event)
	}
}

func RunEvent(ctx context.Context, event interface{}) {
	ctx, span := internal.StartSpan(ctx, "RunEvent")
	defer span.End()

	switch e := event.(type) {
	case func(context.Context):
		e(ctx)
	case func():
		e()
	case nil:
		// TODO: ignoring nil events (can be generated after ctx cancellation) ?
	default:
		fmt.Printf("Unknown event type: %T\n", event)
	}

}
