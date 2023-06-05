package simplerouting

import (
	"time"

	"github.com/libp2p/go-libp2p-kad-dht/events"
	eq "github.com/libp2p/go-libp2p-kad-dht/events/eventqueue"
	"github.com/libp2p/go-libp2p-kad-dht/internal/key"
	"github.com/libp2p/go-libp2p-kad-dht/network/endpoint"
	"github.com/libp2p/go-libp2p-kad-dht/routingtable"

	"github.com/libp2p/go-libp2p/core/protocol"
)

type SimpleRouting struct {
	msgEndpoint endpoint.Endpoint
	rt          routingtable.RoutingTable

	eventQueue   eq.EventQueue
	eventPlanner events.EventPlanner

	queryConcurrency      int
	queryTimeout          time.Duration
	maxConcurrentRequests int
	protocolID            protocol.ID

	//lock *sync.Mutex
	// list of ongoing queries, useful if we want to limit the queries
}

func NewSimpleRouting(self key.KadKey, msgEndpoint endpoint.Endpoint,
	rt routingtable.RoutingTable, queue eq.EventQueue, ep events.EventPlanner,
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
		eventQueue:            queue,
		eventPlanner:          ep,
		queryConcurrency:      cfg.QueryConcurrency,
		queryTimeout:          cfg.QueryTimeout,
		maxConcurrentRequests: cfg.MaxConcurrentRequests,
		protocolID:            cfg.ProtocolID,
	}, nil
}
