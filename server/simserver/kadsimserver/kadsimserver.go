package kadsimserver

import (
	"context"
	"time"

	"github.com/libp2p/go-libp2p-kad-dht/network/address"
	"github.com/libp2p/go-libp2p-kad-dht/network/endpoint"
	"github.com/libp2p/go-libp2p-kad-dht/network/message"
	"github.com/libp2p/go-libp2p-kad-dht/network/message/simmessage"
	"github.com/libp2p/go-libp2p-kad-dht/routingtable"
	"github.com/libp2p/go-libp2p-kad-dht/util"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

const (
	peerstoreTTL  = 10 * time.Minute
	nClosestPeers = 20
)

type KadSimServer struct {
	rt       routingtable.RoutingTable
	endpoint endpoint.Endpoint
}

func NewKadSimServer(rt routingtable.RoutingTable, endpoint endpoint.Endpoint) *KadSimServer {
	return &KadSimServer{
		rt:       rt,
		endpoint: endpoint,
	}
}

func (s *KadSimServer) HandleFindNodeRequest(ctx context.Context, rpeer address.NetworkAddress,
	msg message.MinKadMessage, sendFn endpoint.ResponseHandlerFn) {

	req, ok := msg.(*simmessage.SimMessage)
	if !ok {
		// invalid request
		sendFn(ctx, nil, message.ErrInvalidRequest)
		return
	}

	s.endpoint.MaybeAddToPeerstore(ctx, rpeer, peerstoreTTL)

	_, span := util.StartSpan(ctx, "SimServer.HandleFindNodeRequest", trace.WithAttributes(
		attribute.Stringer("Requester", rpeer.NodeID()),
		attribute.Stringer("Target", req.Target())))
	defer span.End()

	peers, err := s.rt.NearestPeers(ctx, req.Target(), nClosestPeers)
	if err != nil {
		span.RecordError(err)
		sendFn(ctx, nil, err)
		return
	}

	span.AddEvent("Nearest peers", trace.WithAttributes(
		attribute.Int("count", len(peers)),
	))

	resp := simmessage.NewSimResponse(peers)

	sendFn(ctx, resp, nil)
}
