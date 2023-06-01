package fifo

import "sync"

type element struct {
	next *element
	data interface{}
}

type Queue struct {
	head *element
	tail *element
	size uint
	lock sync.RWMutex

	newsChan chan struct{}
}

func NewQueue() *Queue {
	return &Queue{
		newsChan: make(chan struct{}, 1),
	}
}

func (q *Queue) Enqueue(e interface{}) {
	q.lock.Lock()
	defer q.lock.Unlock()

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

func (q *Queue) Dequeue() interface{} {
	q.lock.Lock()
	defer q.lock.Unlock()

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
