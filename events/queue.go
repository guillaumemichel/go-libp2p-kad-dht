package events

type element struct {
	next *element
	data interface{}
}

type Queue struct {
	head *element
	tail *element
	size uint
}

func New() *Queue {
	return &Queue{}
}

func (q *Queue) Enqueue(e interface{}) {
	elem := &element{data: e}
	if q.size == 0 {
		q.head = elem
		q.tail = elem
	} else {
		q.tail.next = elem
		q.tail = elem
	}
	q.size++
}

func (q *Queue) Dequeue() interface{} {
	if q.size == 0 {
		return nil
	}
	elem := q.head
	q.head = elem.next
	q.size--
	return elem.data
}

func (q *Queue) Empty() bool {
	return q.head == nil
}

func (q *Queue) Size() uint {
	return q.size
}