package kadsimserver

import (
	"context"

	"github.com/libp2p/go-libp2p-kad-dht/dht/consts"
	"github.com/libp2p/go-libp2p-kad-dht/key"
	"github.com/libp2p/go-libp2p-kad-dht/network/address"
	"github.com/libp2p/go-libp2p-kad-dht/network/endpoint"
	"github.com/libp2p/go-libp2p-kad-dht/network/message"
	"github.com/libp2p/go-libp2p-kad-dht/network/message/simmessage"
	"github.com/libp2p/go-libp2p-kad-dht/routingtable"
	"github.com/libp2p/go-libp2p-kad-dht/util"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
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
		return
	}

	s.endpoint.MaybeAddToPeerstore(ctx, rpeer, consts.PeerstoreTTL)

	_, span := util.StartSpan(ctx, "SimServer.HandleFindNodeRequest", trace.WithAttributes(
		attribute.Stringer("Requester", address.ID(rpeer)),
		attribute.Stringer("Target", req.Target())))
	defer span.End()

	peers := s.rt.NearestPeers(ctx, req.Target(), consts.NClosestPeers)

	span.AddEvent("Nearest peers", trace.WithAttributes(
		attribute.Int("count", len(peers)),
	))

	kadPeers := make([]key.KadKey, len(peers))
	for i, p := range peers {
		kadPeers[i] = p.(key.KadKey)
	}
	resp := simmessage.NewSimResponse(kadPeers)

	sendFn(ctx, resp, nil)
}
