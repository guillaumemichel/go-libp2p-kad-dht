package ipfssimserver

import (
	"context"
	"time"

	"github.com/libp2p/go-libp2p-kad-dht/network/address"
	"github.com/libp2p/go-libp2p-kad-dht/network/address/peerid"
	"github.com/libp2p/go-libp2p-kad-dht/network/endpoint"
	"github.com/libp2p/go-libp2p-kad-dht/network/message"
	"github.com/libp2p/go-libp2p-kad-dht/network/message/ipfskadv1"
	"github.com/libp2p/go-libp2p-kad-dht/routingtable"
	"github.com/libp2p/go-libp2p-kad-dht/server/simserver"
	"github.com/libp2p/go-libp2p-kad-dht/util"

	"github.com/libp2p/go-libp2p/core/peer"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

var _ simserver.SimServer = (*SimServer)(nil)

type SimServer struct {
	rt       routingtable.RoutingTable
	endpoint endpoint.NetworkedEndpoint

	peerstoreTTL              time.Duration
	numberOfCloserPeersToSend int
}

func NewIpfsSimServer(rt routingtable.RoutingTable, endpoint endpoint.NetworkedEndpoint,
	options ...Option) *SimServer {

	var cfg Config
	if err := cfg.Apply(append([]Option{DefaultConfig}, options...)...); err != nil {
		return nil
	}

	return &SimServer{
		rt:                        rt,
		endpoint:                  endpoint,
		peerstoreTTL:              cfg.PeerstoreTTL,
		numberOfCloserPeersToSend: cfg.NumberOfCloserPeersToSend,
	}
}

func (s *SimServer) HandleFindNodeRequest(ctx context.Context, rpeer address.NetworkAddress,
	msg message.MinKadMessage, replyFn simserver.ReplyFn) {

	req, ok := msg.(*ipfskadv1.Message)
	if !ok {
		// invalid request
		return
	}

	s.endpoint.MaybeAddToPeerstore(ctx, rpeer, s.peerstoreTTL)

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

	peers, err := s.rt.NearestPeers(ctx, pid.Key(), s.numberOfCloserPeersToSend)
	if err != nil {
		span.RecordError(err)
		return
	}

	span.AddEvent("Nearest peers", trace.WithAttributes(
		attribute.Int("count", len(peers)),
		attribute.String("peer", peers[0].String())))

	resp := ipfskadv1.FindPeerResponse(peers, s.endpoint)

	replyFn(resp)
}
