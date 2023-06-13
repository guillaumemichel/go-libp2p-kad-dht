package simpledispatcher

import (
	"context"
	"sort"
	"time"

	"go.opentelemetry.io/otel/trace"

	"github.com/benbjohnson/clock"
	"github.com/libp2p/go-libp2p-kad-dht/events"
	"github.com/libp2p/go-libp2p-kad-dht/events/scheduler"
	"github.com/libp2p/go-libp2p-kad-dht/internal"
	"github.com/libp2p/go-libp2p-kad-dht/internal/util"
	"github.com/libp2p/go-libp2p-kad-dht/network/address"
	"github.com/libp2p/go-libp2p-kad-dht/server/simserver"
	"go.opentelemetry.io/otel/attribute"
)

// SimpleDispatcher is a simple implementation of a LoopDispatcher.
type SimpleDispatcher struct {
	clk       *clock.Mock
	peers     map[address.NodeID]scheduler.AwareScheduler
	servers   map[address.NodeID]*simserver.SimServer
	latencies map[address.NodeID]map[address.NodeID]time.Duration
}

// NewSimpleDispatcher creates a new SimpleDispatcher. The provided mock clock
// must be the same as the one used by the schedulers.
func NewSimpleDispatcher(clk *clock.Mock) *SimpleDispatcher {
	return &SimpleDispatcher{
		clk:       clk,
		peers:     make(map[address.NodeID]scheduler.AwareScheduler),
		servers:   make(map[address.NodeID]*simserver.SimServer),
		latencies: make(map[address.NodeID]map[address.NodeID]time.Duration),
	}
}

// AddPeer adds a peer to the dispatcher. The peer must have an associated
// scheduler.AwareScheduler, using the same mock clock as the dispatcher.
func (d *SimpleDispatcher) AddPeer(id address.NodeID, s scheduler.Scheduler, serv *simserver.SimServer) {
	switch s := s.(type) {
	case scheduler.AwareScheduler:
		d.peers[id] = s
	}
	d.servers[id] = serv
}

// RemovePeer removes a peer from the dispatcher.
func (d *SimpleDispatcher) RemovePeer(id address.NodeID) {
	delete(d.peers, id)
	delete(d.latencies, id)
	for _, l := range d.latencies {
		delete(l, id)
	}
	delete(d.servers, id)
}

// DispatchTo immediately dispatches an action to a peer.
func (d *SimpleDispatcher) DispatchTo(ctx context.Context, to address.NodeID, a events.Action) {
	ctx, span := internal.StartSpan(ctx, "SimpleDispatcher.DispatchTo", trace.WithAttributes(
		attribute.String("NodeID", to.String()),
	))
	defer span.End()

	if s, ok := d.peers[to]; ok {
		s.EnqueueAction(ctx, a)
	}
}

// Dispatch immediately dispatches an action to a peer. If a latency is set
// between the two peers, the action will be scheduled to be dispatched after
// the latency.
func (d *SimpleDispatcher) Dispatch(ctx context.Context, from, to address.NodeID,
	a events.Action) {

	if s, ok := d.peers[to]; ok {
		d.DispatchDelay(ctx, from, to, a, s.Now())
	}
}

// DispatchDelay schedules an action to be dispatched to a peer at a given time.
// If a latency is set between the two peers, the action will be scheduled to be
// dispatched after the latency.
func (d *SimpleDispatcher) DispatchDelay(ctx context.Context, from, to address.NodeID,
	a events.Action, t time.Time) {

	if s, ok := d.peers[to]; ok {

		l := d.GetLatency(from, to)

		trigger := t.Add(l)
		now := s.Now()
		if trigger.Before(now) || trigger == now {
			s.EnqueueAction(ctx, a)
		} else {
			s.ScheduleAction(ctx, trigger, a)
		}
	}
}

// SetLatency sets the latency between two peers. The latency is used to
// schedule actions to be dispatched after the latency. It is used to simulate
// network latencies.
func (d *SimpleDispatcher) SetLatency(from, to address.NodeID, l time.Duration) {
	for _, n := range []address.NodeID{from, to} {
		if _, ok := d.peers[n]; !ok {
			return
		}
	}

	if from == to {
		return
	} else if from.String() > to.String() {
		from, to = to, from
	}

	if _, ok := d.latencies[from]; !ok {
		d.latencies[from] = make(map[address.NodeID]time.Duration)
	}
	d.latencies[from][to] = l
}

// GetLatency returns the latency defined between two peers.
func (d *SimpleDispatcher) GetLatency(from, to address.NodeID) time.Duration {
	if from == to {
		return 0
	} else if from.String() > to.String() {
		from, to = to, from
	}

	if latencies, ok := d.latencies[from]; ok {
		if latency, ok := latencies[to]; ok {
			return latency
		}
	}
	return 0
}

type ctxKey string

const (
	ctxPeerKey     = ctxKey("peer")
	ctxTimeKey     = ctxKey("time")
	ctxActionIdKey = ctxKey("actionID")
)

// DispatchLoop runs the dispatch loop. It will run until all peers have no more
// actions to run.
func (d *SimpleDispatcher) DispatchLoop(ctx context.Context) {
	ctx, span := internal.StartSpan(ctx, "SimpleDispatcher.DispatchLoop")
	defer span.End()

	actionID := 0

	// get the next action time for each peer
	nextActions := make(map[address.NodeID]time.Time)
	for id, s := range d.peers {
		nextActions[id] = s.NextActionTime(ctx)
	}
	// TODO: optimize nextActions to be a linked list of actions sorted by time

	for len(nextActions) > 0 {
		span.AddEvent("DispatchLoop iteration")

		// find the time of the next action to be run
		minTime := util.MaxTime
		for _, t := range nextActions {
			if t.Before(minTime) {
				minTime = t
			}
		}

		if minTime == util.MaxTime {
			// no more actions to run
			break
		}

		ctx = context.WithValue(ctx, ctxTimeKey, minTime)

		upNext := make([]address.NodeID, 0)
		for id, t := range nextActions {
			if t == minTime {
				upNext = append(upNext, id)
			}
		}
		// sort the peers by ID to ensure deterministic behavior, because map
		// iteration order is non-deterministic
		sort.Slice(upNext, func(i, j int) bool {
			return upNext[i].String() < upNext[j].String()
		})

		span.AddEvent("DispatchLoop: sorted!")

		// "wait" minTime for the next action
		d.clk.Set(minTime) // slow to execute (because of the mutex?)

		span.AddEvent("DispatchLoop: new time set!")

		for len(upNext) > 0 {
			ongoing := make([]address.NodeID, len(upNext))
			copy(ongoing, upNext)

			upNext = make([]address.NodeID, 0)
			for _, id := range ongoing {
				// run one action for this peer
				ctx = context.WithValue(ctx, ctxPeerKey, id)
				ctx = context.WithValue(ctx, ctxActionIdKey, actionID)
				actionID++
				d.peers[id].RunOne(ctx)
			}
		}

		for id, s := range d.peers {
			t := s.NextActionTime(ctx)
			if t == minTime {
				upNext = append(upNext, id)
			} else {
				nextActions[id] = t
			}
		}
	}
}

func (d *SimpleDispatcher) Server(n address.NodeID) *simserver.SimServer {
	if serv, ok := d.servers[n]; ok {
		return serv
	}
	return nil
}
