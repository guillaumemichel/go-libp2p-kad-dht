package chanqueue

import (
	"context"
	"errors"

	"github.com/libp2p/go-libp2p-kad-dht/events"
	"github.com/libp2p/go-libp2p-kad-dht/internal"
)

// ChanQueue is a trivial queue implementation using a channel
type ChanQueue struct {
	queue chan events.Action
}

// NewChanQueue creates a new queue
func NewChanQueue(ctx context.Context, capacity int) *ChanQueue {
	_, span := internal.StartSpan(ctx, "NewChanQueue")
	defer span.End()

	return &ChanQueue{
		queue: make(chan events.Action, capacity),
	}
}

// Enqueue adds an element to the queue
func (q *ChanQueue) Enqueue(ctx context.Context, e events.Action) {
	_, span := internal.StartSpan(ctx, "ChanQueue.Enqueue")
	defer span.End()

	select {
	case q.queue <- e:
	default:
		span.RecordError(errors.New("cannot write to queue"))
	}
}

// Dequeue reads the next element from the queue, note that this operation is blocking
func (q *ChanQueue) Dequeue(ctx context.Context) events.Action {
	_, span := internal.StartSpan(ctx, "ChanQueue.Dequeue")
	defer span.End()

	if q.Empty() {
		span.AddEvent("empty queue")
		return nil
	}

	return <-q.queue
}

// Empty returns true if the queue is empty
func (q *ChanQueue) Empty() bool {
	return len(q.queue) == 0
}

func (q *ChanQueue) Size() uint {
	return uint(len(q.queue))
}

func (q *ChanQueue) Close() {
	close(q.queue)
}
