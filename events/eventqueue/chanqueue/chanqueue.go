package chanqueue

// ChanQueue is a trivial queue implementation using a channel
type ChanQueue struct {
	queue    chan interface{}
	newsChan chan struct{}
}

// NewChanQueue creates a new queue
func NewChanQueue(capacity int) *ChanQueue {
	return &ChanQueue{
		queue:    make(chan interface{}, capacity),
		newsChan: make(chan struct{}, 1),
	}
}

// Enqueue adds an element to the queue
func (q *ChanQueue) Enqueue(e interface{}) {
	q.queue <- e

	select {
	case q.newsChan <- struct{}{}:
	default:
	}
}

// Dequeue reads the next element from the queue, note that this operation is blocking
func (q *ChanQueue) Dequeue() interface{} {
	return <-q.queue
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
