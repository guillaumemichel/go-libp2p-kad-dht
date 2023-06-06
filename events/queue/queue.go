package queue

import "github.com/libp2p/go-libp2p-kad-dht/events"

type EventQueue interface {
	Enqueue(events.Action)
	Dequeue() events.Action
	NewsChan() <-chan struct{}

	Size() uint
	Close()
}

type EventQueueEnqueueMany interface {
	EventQueue
	EnqueueMany([]events.Action)
}

func EnqueueMany(q EventQueue, actions []events.Action) {
	switch queue := q.(type) {
	case EventQueueEnqueueMany:
		queue.EnqueueMany(actions)
	default:
		for _, a := range actions {
			q.Enqueue(a)
		}
	}
}

type EventQueueDequeueMany interface {
	DequeueMany(int) []events.Action
}

func DequeueMany(q EventQueue, n int) []events.Action {
	switch queue := q.(type) {
	case EventQueueDequeueMany:
		return queue.DequeueMany(n)
	default:
		actions := make([]events.Action, 0, n)
		for i := 0; i < n; i++ {
			if a := q.Dequeue(); a != nil {
				actions = append(actions, a)
			} else {
				break
			}
		}
		return actions
	}
}

type EventQueueDequeueAll interface {
	DequeueAll() []events.Action
}

func DequeueAll(q EventQueue) []events.Action {
	switch queue := q.(type) {
	case EventQueueDequeueAll:
		return queue.DequeueAll()
	default:
		actions := make([]events.Action, 0, q.Size())
		for a := q.Dequeue(); a != nil; a = q.Dequeue() {
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
