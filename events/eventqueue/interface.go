package eventqueue

type EventQueue interface {
	Enqueue(interface{})
	Dequeue() interface{}
	NewsChan() <-chan struct{}

	Size() uint
	Close()
}

func Enqueue(q EventQueue, e interface{}) {
	q.Enqueue(e)
}

func Dequeue(q EventQueue) interface{} {
	return q.Dequeue()
}

func NewsChan(q EventQueue) <-chan struct{} {
	return q.NewsChan()
}

func Empty(q EventQueue) bool {
	switch queue := q.(type) {
	case EventQueueWithEmpty:
		return queue.Empty()
	default:
		return q.Size() == 0
	}
}

func Size(q EventQueue) uint {
	return q.Size()
}

func Close(q EventQueue) {
	q.Close()
}

type EventQueueWithEmpty interface {
	EventQueue
	Empty() bool
}
