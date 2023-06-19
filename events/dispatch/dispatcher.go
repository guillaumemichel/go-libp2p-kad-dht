package dispatch

import (
	"context"
	"time"

	"github.com/libp2p/go-libp2p-kad-dht/events/action"
	"github.com/libp2p/go-libp2p-kad-dht/events/scheduler"
	"github.com/libp2p/go-libp2p-kad-dht/network/address"
	"github.com/libp2p/go-libp2p-kad-dht/server/simserver"
)

type StreamID uint64

// Dispatcher is an interface for dispatching actions to peers' schedulers.
type Dispatcher interface {
	// AddPeer adds a peer to the dispatcher.
	AddPeer(address.NodeID, scheduler.Scheduler, *simserver.SimServer)
	// RemovePeer removes a peer from the dispatcher.
	RemovePeer(id address.NodeID)

	// DispatchTo immediately dispatches an action to a peer.
	DispatchTo(context.Context, address.NodeID, action.Action)
}

// LatencyDispatcher is an interface for dispatching actions to peers' schedulers
// with a defined latency between peers.
type LatencyDispatcher interface {
	Dispatcher
	// Dispatch immediately dispatches an action to a peer. If a latency is set
	// between the two peers, the action will be scheduled to be run after
	// this latency.
	Dispatch(context.Context, address.NodeID, address.NodeID, action.Action)

	// SetLatency sets the latency between two peers.
	SetLatency(address.NodeID, address.NodeID, time.Duration)
	// GetLatency returns the latency between two peers.
	GetLatency(address.NodeID, address.NodeID) time.Duration
}

// DelayLatencyDispatcher is an interface for dispatching actions to peers'
// schedulers. It supports scheduling actions to be run at a given time,
// leveraging the latency between peers.
type DelayLatencyDispatcher interface {
	LatencyDispatcher

	// DispatchDelay schedules an action to be dispatched to a peer at a given
	// time. If a latency is set between the two peers, the action is
	// scheduled to be run after the provided time + latency.
	DispatchDelay(context.Context, address.NodeID, address.NodeID, action.Action, time.Time)
}

// LoopDispatcher is an interface for dispatching actions to peers' schedulers.
// All scheduled actions can be run sequentially using DispatchLoop.
type LoopDispatcher interface {
	DelayLatencyDispatcher

	// DispatchLoop runs a loop that dispatches all scheduled actions to peers'
	// schedulers. It "runs the simulation" of scheduled actions.
	DispatchLoop(context.Context)
}
