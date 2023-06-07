package simpledispatcher

import (
	"context"
	"sort"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/libp2p/go-libp2p-kad-dht/events"
	"github.com/libp2p/go-libp2p-kad-dht/events/scheduler"
	"github.com/libp2p/go-libp2p-kad-dht/internal/util"
	"github.com/libp2p/go-libp2p-kad-dht/network/address"
)

// SimpleDispatcher is a simple implementation of a LoopDispatcher.
type SimpleDispatcher struct {
	clk       *clock.Mock
	peers     map[address.NodeID]scheduler.AwareScheduler
	latencies map[address.NodeID]map[address.NodeID]time.Duration
}

// NewSimpleDispatcher creates a new SimpleDispatcher. The provided mock clock
// must be the same as the one used by the schedulers.
func NewSimpleDispatcher(clk *clock.Mock) *SimpleDispatcher {
	return &SimpleDispatcher{
		clk:       clk,
		peers:     make(map[address.NodeID]scheduler.AwareScheduler),
		latencies: make(map[address.NodeID]map[address.NodeID]time.Duration),
	}
}

// AddPeer adds a peer to the dispatcher. The peer must have an associated
// scheduler, using the same mock clock as the dispatcher.
func (d *SimpleDispatcher) AddPeer(id address.NodeID, s scheduler.AwareScheduler) {
	d.peers[id] = s
}

// RemovePeer removes a peer from the dispatcher.
func (d *SimpleDispatcher) RemovePeer(id address.NodeID) {
	delete(d.peers, id)
	delete(d.latencies, id)
	for _, l := range d.latencies {
		delete(l, id)
	}
}

// DispatchTo immediately dispatches an action to a peer.
func (d *SimpleDispatcher) DispatchTo(ctx context.Context, to address.NodeID, a events.Action) {
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
	// once DispatchLoop is called, the peers, latencies and scheduled actions
	// cannot be mondified anymore
	actionID := 0

	// get the next action time for each peer
	nextActions := make(map[address.NodeID]time.Time)
	for id, s := range d.peers {
		t := s.NextActionTime(ctx)
		if t == util.MaxTime {
			delete(nextActions, id)
		} else {
			nextActions[id] = t
		}
	}
	// TODO: optimize nextActions to be a linked list of actions sorted by time

	for len(nextActions) > 0 {
		// find the time of the next action to be run
		minTime := util.MaxTime
		for _, t := range nextActions {
			if t.Before(minTime) {
				minTime = t
			}
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

		// "wait" minTime for the next action
		d.clk.Set(minTime)

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
				t := d.peers[id].NextActionTime(ctx)
				if t == util.MaxTime {
					// no more actions for this peer
					delete(nextActions, id)
				} else {
					if t == minTime {
						// multiple actions are scheduled at minTime for this peer
						upNext = append(upNext, id)
					} else {
						// peer next action is scheduled at a later time
						nextActions[id] = t
					}
				}
			}
		}
	}
}
