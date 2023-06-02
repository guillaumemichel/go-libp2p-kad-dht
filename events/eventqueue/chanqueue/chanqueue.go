package chanqueue

import (
	"context"
	"errors"

	"github.com/libp2p/go-libp2p-kad-dht/internal"
)

// ChanQueue is a trivial queue implementation using a channel
type ChanQueue struct {
	ctx      context.Context
	queue    chan interface{}
	newsChan chan struct{}
}

// NewChanQueue creates a new queue
func NewChanQueue(ctx context.Context, capacity int) *ChanQueue {
	ctx, span := internal.StartSpan(ctx, "NewChanQueue")
	defer span.End()

	return &ChanQueue{
		ctx:      ctx,
		queue:    make(chan interface{}, capacity),
		newsChan: make(chan struct{}, 1),
	}
}

// Enqueue adds an element to the queue
func (q *ChanQueue) Enqueue(e interface{}) {
	_, span := internal.StartSpan(q.ctx, "ChanQueue.Enqueue")
	defer span.End()

	select {
	case <-q.ctx.Done():
		return
	case q.queue <- e:
	default:
		span.RecordError(errors.New("cannot write to queue"))
	}

	select {
	case q.newsChan <- struct{}{}:
	default:
	}
}

// Dequeue reads the next element from the queue, note that this operation is blocking
func (q *ChanQueue) Dequeue() interface{} {
	_, span := internal.StartSpan(q.ctx, "ChanQueue.Dequeue")
	defer span.End()

	if q.Empty() {
		span.AddEvent("empty queue")
		return nil
	}

	select {
	case <-q.ctx.Done():
		span.RecordError(q.ctx.Err())
		return nil
	case e := <-q.queue:
		return e
	}
}

// Empty returns true if the queue is empty
func (q *ChanQueue) Empty() bool {
	return len(q.queue) == 0
}

func (q *ChanQueue) Size() uint {
	return uint(len(q.queue))
}

func (q *ChanQueue) NewsChan() <-chan struct{} {
	return q.newsChan
}

func (q *ChanQueue) Close() {
	close(q.queue)
	close(q.newsChan)
}
