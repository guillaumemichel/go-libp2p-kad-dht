package queue

import (
	"context"

	"github.com/libp2p/go-libp2p-kad-dht/events"
)

type EventQueue interface {
	Enqueue(context.Context, events.Action)
	Dequeue(context.Context) events.Action
	NewsChan() <-chan struct{}

	Size() uint
	Close()
}

type EventQueueEnqueueMany interface {
	EventQueue
	EnqueueMany(context.Context, []events.Action)
}

func EnqueueMany(ctx context.Context, q EventQueue, actions []events.Action) {
	switch queue := q.(type) {
	case EventQueueEnqueueMany:
		queue.EnqueueMany(ctx, actions)
	default:
		for _, a := range actions {
			q.Enqueue(ctx, a)
		}
	}
}

type EventQueueDequeueMany interface {
	DequeueMany(context.Context, int) []events.Action
}

func DequeueMany(ctx context.Context, q EventQueue, n int) []events.Action {
	switch queue := q.(type) {
	case EventQueueDequeueMany:
		return queue.DequeueMany(ctx, n)
	default:
		actions := make([]events.Action, 0, n)
		for i := 0; i < n; i++ {
			if a := q.Dequeue(ctx); a != nil {
				actions = append(actions, a)
			} else {
				break
			}
		}
		return actions
	}
}

type EventQueueDequeueAll interface {
	DequeueAll(context.Context) []events.Action
}

func DequeueAll(ctx context.Context, q EventQueue) []events.Action {
	switch queue := q.(type) {
	case EventQueueDequeueAll:
		return queue.DequeueAll(ctx)
	default:
		actions := make([]events.Action, 0, q.Size())
		for a := q.Dequeue(ctx); a != nil; a = q.Dequeue(ctx) {
			actions = append(actions, a)
		}
		return actions
	}
}

type EventQueueWithEmpty interface {
	EventQueue
	Empty() bool
}

func Empty(q EventQueue) bool {
	switch queue := q.(type) {
	case EventQueueWithEmpty:
		return queue.Empty()
	default:
		return q.Size() == 0
	}
}
