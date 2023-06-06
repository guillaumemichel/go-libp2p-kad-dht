package simpledispatcher

import (
	"context"
	"time"

	"github.com/benbjohnson/clock"

	"github.com/libp2p/go-libp2p-kad-dht/events"
	"github.com/libp2p/go-libp2p-kad-dht/events/scheduler"
	"github.com/libp2p/go-libp2p-kad-dht/network/address"
)

type SimpleDispatcher struct {
	peers     map[address.NodeID]scheduler.Scheduler
	latencies map[address.NodeID]map[address.NodeID]time.Duration
}

func NewSimpleDispatcher(clk clock.Clock) *SimpleDispatcher {
	return &SimpleDispatcher{
		peers:     make(map[address.NodeID]scheduler.Scheduler),
		latencies: make(map[address.NodeID]map[address.NodeID]time.Duration),
	}
}

func (d *SimpleDispatcher) AddPeer(id address.NodeID, s scheduler.Scheduler) {
	d.peers[id] = s
}

func (d *SimpleDispatcher) RemovePeer(id address.NodeID) {
	delete(d.peers, id)
}

func (d *SimpleDispatcher) Dispatch(ctx context.Context, from, to address.NodeID,
	a events.Action) {

	if s, ok := d.peers[to]; ok {

		l := d.GetLatency(from, to)
		if l == 0 {
			s.EnqueueAction(ctx, a)
		} else {
			scheduler.ScheduleActionIn(ctx, s, l, a)
		}
	}
}

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
