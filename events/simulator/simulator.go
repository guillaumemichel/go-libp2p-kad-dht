package simulator

import (
	"context"

	"github.com/libp2p/go-libp2p-kad-dht/events/scheduler"
)

type Simulator interface {
	AddPeer(scheduler.AwareScheduler)
	RemovePeer(scheduler.AwareScheduler)
	Run(context.Context)
}

func AddPeers(s Simulator, schedulers ...scheduler.AwareScheduler) {
	for _, sched := range schedulers {
		s.AddPeer(sched)
	}
}

func RemovePeers(s Simulator, schedulers ...scheduler.AwareScheduler) {
	for _, sched := range schedulers {
		s.RemovePeer(sched)
	}
}
