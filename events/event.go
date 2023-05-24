package events

import "sync"

// types of Events:
// 1. Client Request
//   a) FindNode
//   b) GetClosestPeers
//   c) FindProvs
//   d) StartProviding
//   e) StopProviding
//   f) ListProviding
//   g) GetValue (IPNS)
//   h) PutValue (IPNS)
//   i) GetPublicKey ???
// 2. Server Request: the node is in server mode and got a request from a remote peer
// 3. Message Response: a remote peer responded to a request from this node
// 4. Message Timeout: a remote peer failed to respond to a request from this node in time
// 5. IO: the node has to do some IO (e.g. write to Provider Store on disk)
// 6. IO Timeout: the node failed to do some IO in time

type Event struct {
	F func()
}

type EventsManager struct {
	lock   sync.Mutex
	active bool

	queue *Queue
}

func NewEventsManager() *EventsManager {
	return &EventsManager{queue: NewQueue()}
}

func NewEvent(em *EventsManager, e interface{}) {
	em.lock.Lock()
	// check if there is an active thread handling events
	if em.active {
		em.lock.Unlock()
		// if a thread is already handling events, enqueue the new event
		em.queue.Enqueue(e)
		return
	}
	// if no thread is handling events, this thread becomes active
	em.active = true
	em.lock.Unlock()

	// handle the new event
	handleEvent(e)

	// if new events were enqueued while this thread was active, handle them
	em.lock.Lock()
	for !em.queue.Empty() {
		em.lock.Unlock()

		e := em.queue.Dequeue()
		handleEvent(e)

		em.lock.Lock()
	}
	// once all events have been handled, this thread is no longer active
	em.active = false
	em.lock.Unlock()
}

func handleEvent(e interface{}) {
	switch e := e.(type) {
	case func() error:
		e()
	case *Event:
	default:
		panic("unknown event type") // TODO: handle this
	}
}
