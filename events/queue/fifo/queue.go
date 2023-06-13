package fifo

import (
	"context"
	"sync"

	"github.com/libp2p/go-libp2p-kad-dht/events"
	"github.com/libp2p/go-libp2p-kad-dht/internal"
)

type element struct {
	next *element
	data events.Action
}

type Queue struct {
	ctx  context.Context
	head *element
	tail *element
	size uint
	lock sync.RWMutex

	newsChan chan struct{}
}

func NewQueue(ctx context.Context) *Queue {
	return &Queue{
		ctx:      ctx,
		newsChan: make(chan struct{}, 1),
	}
}

func (q *Queue) Enqueue(ctx context.Context, e events.Action) {
	q.lock.Lock()
	defer q.lock.Unlock()

	_, span := internal.StartSpan(ctx, "FifoQueue.Enqueue")
	defer span.End()

	elem := &element{data: e}
	if q.size == 0 {
		q.head = elem
		q.tail = elem
	} else {
		q.tail.next = elem
		q.tail = elem
	}
	q.size++

	select {
	case q.newsChan <- struct{}{}:
	default:
	}
}

func (q *Queue) Dequeue(ctx context.Context) events.Action {
	q.lock.Lock()
	defer q.lock.Unlock()

	_, span := internal.StartSpan(ctx, "FifoQueue.Dequeue")
	defer span.End()

	if q.size == 0 {
		return nil
	}
	elem := q.head
	q.head = elem.next
	q.size--
	return elem.data
}

func (q *Queue) Empty() bool {
	q.lock.RLock()
	defer q.lock.RUnlock()

	return q.head == nil
}

func (q *Queue) Size() uint {
	q.lock.RLock()
	defer q.lock.RUnlock()

	return q.size
}

func (q *Queue) NewsChan() <-chan struct{} {
	return q.newsChan
}

func (q *Queue) Close() {
	q.lock.Lock()
	defer q.lock.Unlock()

	q.head = nil
	q.tail = nil
	q.size = 0

	close(q.newsChan)
}
