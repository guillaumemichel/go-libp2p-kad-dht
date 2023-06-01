package simplerouting

import (
	"sync"
	"time"

	"github.com/libp2p/go-libp2p-kad-dht/events"
	eq "github.com/libp2p/go-libp2p-kad-dht/events/eventqueue"
	"github.com/libp2p/go-libp2p-kad-dht/network"
	"github.com/libp2p/go-libp2p-kad-dht/routingtable"

	"github.com/libp2p/go-libp2p/core/protocol"
)

type SimpleRouting struct {
	msgEndpoint *network.MessageEndpoint
	rt          routingtable.RoutingTable

	eventqueue eq.EventQueue
	sched      events.Scheduler

	queryConcurrency      int
	queryTimeout          time.Duration
	maxConcurrentRequests int
	protocolID            protocol.ID

	lock sync.Mutex
}

func NewSimpleRouting(msgEndpoint *network.MessageEndpoint, rt routingtable.RoutingTable,
	queue eq.EventQueue, sched events.Scheduler, options ...Option) (*SimpleRouting, error) {

	var cfg Config
	if err := cfg.Apply(append([]Option{DefaultConfig}, options...)...); err != nil {
		return nil, err
	}
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return &SimpleRouting{
		msgEndpoint:           msgEndpoint,
		rt:                    rt,
		eventqueue:            queue,
		sched:                 sched,
		queryConcurrency:      cfg.QueryConcurrency,
		queryTimeout:          cfg.QueryTimeout,
		maxConcurrentRequests: cfg.MaxConcurrentRequests,
		protocolID:            cfg.ProtocolID,
	}, nil
}
