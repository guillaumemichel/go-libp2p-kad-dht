package kadsimserver

import (
	"context"
	"time"

	"github.com/libp2p/go-libp2p-kad-dht/network/address"
	"github.com/libp2p/go-libp2p-kad-dht/network/endpoint"
	"github.com/libp2p/go-libp2p-kad-dht/network/message"
	"github.com/libp2p/go-libp2p-kad-dht/network/message/simmessage"
	"github.com/libp2p/go-libp2p-kad-dht/routingtable"
	"github.com/libp2p/go-libp2p-kad-dht/server/simserver"
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
	msg message.MinKadMessage, replyFn simserver.ReplyFn) {

	req, ok := msg.(*simmessage.SimMessage)
	if !ok {
		// invalid request, don't reply
		return
	}

	target := req.Target()
	if target == nil {
		// invalid request, don't reply
		return
	}

	s.endpoint.MaybeAddToPeerstore(ctx, rpeer, peerstoreTTL)

	_, span := util.StartSpan(ctx, "SimServer.HandleFindNodeRequest", trace.WithAttributes(
		attribute.Stringer("Requester", rpeer.NodeID()),
		attribute.Stringer("Target", target)))
	defer span.End()

	peers, err := s.rt.NearestPeers(ctx, *target, nClosestPeers)
	if err != nil {
		span.RecordError(err)
		replyFn(nil)
		return
	}

	span.AddEvent("Nearest peers", trace.WithAttributes(
		attribute.Int("count", len(peers)),
	))

	resp := simmessage.NewSimResponse(peers)

	replyFn(resp)
}
