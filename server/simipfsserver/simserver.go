package simipfsserver

import (
	"context"

	"github.com/libp2p/go-libp2p-kad-dht/dht/consts"
	"github.com/libp2p/go-libp2p-kad-dht/key"
	"github.com/libp2p/go-libp2p-kad-dht/network/address"
	"github.com/libp2p/go-libp2p-kad-dht/network/endpoint"
	"github.com/libp2p/go-libp2p-kad-dht/network/message"
	"github.com/libp2p/go-libp2p-kad-dht/network/message/ipfskadv1"
	"github.com/libp2p/go-libp2p-kad-dht/routingtable"
	"github.com/libp2p/go-libp2p-kad-dht/util"

	"github.com/libp2p/go-libp2p/core/peer"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type SimServer struct {
	rt       routingtable.RoutingTable
	endpoint endpoint.Endpoint
}

func NewSimServer(rt routingtable.RoutingTable, endpoint endpoint.Endpoint) *SimServer {
	return &SimServer{
		rt:       rt,
		endpoint: endpoint,
	}
}

func (s *SimServer) HandleFindNodeRequest(ctx context.Context, rpeer address.NetworkAddress,
	msg message.MinKadMessage, sendFn endpoint.ResponseHandlerFn) {

	req, ok := msg.(*ipfskadv1.Message)
	if !ok {
		// invalid request
		return
	}

	s.endpoint.MaybeAddToPeerstore(ctx, rpeer, consts.PeerstoreTTL)

	p := peer.ID("")
	if p.UnmarshalBinary(req.GetKey()) != nil {
		// invalid requested key (not a peer.ID)
		return
	}

	_, span := util.StartSpan(ctx, "SimServer.HandleFindNodeRequest", trace.WithAttributes(
		attribute.Stringer("Requester", address.ID(rpeer)),
		attribute.Stringer("Target", p)))
	defer span.End()

	peers := s.rt.NearestPeers(ctx, key.PeerKadID(p), consts.NClosestPeers)

	span.AddEvent("Nearest peers", trace.WithAttributes(
		attribute.Int("count", len(peers)),
		attribute.String("peer", peers[0].String())))

	resp := ipfskadv1.FindPeerResponse(p, peers, s.endpoint)

	sendFn(ctx, resp, nil)
}
