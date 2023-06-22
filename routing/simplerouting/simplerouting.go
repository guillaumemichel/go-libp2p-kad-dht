package simplerouting

import (
	"time"

	"github.com/libp2p/go-libp2p-kad-dht/events/scheduler"
	"github.com/libp2p/go-libp2p-kad-dht/key"
	"github.com/libp2p/go-libp2p-kad-dht/network/address"
	"github.com/libp2p/go-libp2p-kad-dht/network/endpoint/libp2pendpoint"
	"github.com/libp2p/go-libp2p-kad-dht/routingtable"
)

type SimpleRouting struct {
	msgEndpoint *libp2pendpoint.Libp2pEndpoint
	rt          routingtable.RoutingTable
	protoID     address.ProtocolID

	sched scheduler.Scheduler

	queryConcurrency      int
	queryTimeout          time.Duration
	maxConcurrentRequests int

	//lock *sync.Mutex
	// list of ongoing queries, useful if we want to limit the queries
}

func NewSimpleRouting(self key.KadKey, msgEndpoint *libp2pendpoint.Libp2pEndpoint,
	rt routingtable.RoutingTable, sched scheduler.Scheduler, proto address.ProtocolID,
	options ...Option) (*SimpleRouting, error) {

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
		protoID:               proto,
		sched:                 sched,
		queryConcurrency:      cfg.QueryConcurrency,
		queryTimeout:          cfg.QueryTimeout,
		maxConcurrentRequests: cfg.MaxConcurrentRequests,
	}, nil
}
