package simipfsserver

import (
	"context"
	"time"

	"github.com/libp2p/go-libp2p-kad-dht/network/address"
	"github.com/libp2p/go-libp2p-kad-dht/network/address/peerid"
	"github.com/libp2p/go-libp2p-kad-dht/network/endpoint"
	"github.com/libp2p/go-libp2p-kad-dht/network/message"
	"github.com/libp2p/go-libp2p-kad-dht/network/message/ipfskadv1"
	"github.com/libp2p/go-libp2p-kad-dht/routingtable"
	"github.com/libp2p/go-libp2p-kad-dht/util"

	"github.com/libp2p/go-libp2p/core/peer"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

const (
	peerstoreTTL  = 10 * time.Minute
	nClosestPeers = 20
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

	s.endpoint.MaybeAddToPeerstore(ctx, rpeer, peerstoreTTL)

	p := peer.ID("")
	if p.UnmarshalBinary(req.GetKey()) != nil {
		// invalid requested key (not a peer.ID)
		return
	}
	pid := peerid.PeerID{ID: p}

	_, span := util.StartSpan(ctx, "SimServer.HandleFindNodeRequest", trace.WithAttributes(
		attribute.Stringer("Requester", rpeer.NodeID()),
		attribute.Stringer("Target", p)))
	defer span.End()

	peers, err := s.rt.NearestPeers(ctx, pid.Key(), nClosestPeers)
	if err != nil {
		span.RecordError(err)
		return
	}

	span.AddEvent("Nearest peers", trace.WithAttributes(
		attribute.Int("count", len(peers)),
		attribute.String("peer", peers[0].String())))

	resp := ipfskadv1.FindPeerResponse(pid, peers, s.endpoint)

	sendFn(ctx, resp, nil)
}
