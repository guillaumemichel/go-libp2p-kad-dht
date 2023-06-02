package events

import (
	"context"
	"sync"

	eq "github.com/libp2p/go-libp2p-kad-dht/events/eventqueue"
	"github.com/libp2p/go-libp2p-kad-dht/events/eventqueue/fifo"

	"github.com/libp2p/go-libp2p-kad-dht/internal"
)

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
	done   bool

	queue eq.EventQueue
}

func NewEventsManager(ctx context.Context) *EventsManager {
	return &EventsManager{queue: fifo.NewQueue(ctx)}
}

func NewEvent(ctx context.Context, em *EventsManager, e interface{}) {
	em.lock.Lock()
	// check if there is an active thread handling events
	if em.active {
		em.lock.Unlock()

		_, span := internal.StartSpan(ctx, "events.NewEvent queued")
		defer span.End()

		// if a thread is already handling events, enqueue the new event
		em.queue.Enqueue(e)
		return
	}
	// if no thread is handling events, this thread becomes active
	em.active = true
	em.lock.Unlock()

	// handle the new event
	handleEvent(ctx, e)

	// if new events were enqueued while this thread was active, handle them
	em.lock.Lock()
	for !eq.Empty(em.queue) && !em.done {
		em.lock.Unlock()

		e := em.queue.Dequeue()
		handleEvent(ctx, e)

		em.lock.Lock()
	}
	// once all events have been handled, this thread is no longer active
	em.active = false
	em.lock.Unlock()
}

// StopEventManager stops the event manager from handling events
func StopEventManager(ctx context.Context, em *EventsManager) {
	_, span := internal.StartSpan(ctx, "events.StopEventManager")
	defer span.End()

	em.lock.Lock()
	em.done = true
	em.lock.Unlock()
}

func handleEvent(ctx context.Context, e interface{}) {
	switch e := e.(type) {
	case func(ctx context.Context):
		e(ctx)
	case *Event:
	default:
		panic("unknown event type") // TODO: handle this
	}
}
