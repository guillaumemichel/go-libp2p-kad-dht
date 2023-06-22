package kadsimserver

import (
	"context"
	"time"

	"github.com/libp2p/go-libp2p-kad-dht/key"
	"github.com/libp2p/go-libp2p-kad-dht/network/address"
	"github.com/libp2p/go-libp2p-kad-dht/network/address/peerid"
	"github.com/libp2p/go-libp2p-kad-dht/network/endpoint"
	"github.com/libp2p/go-libp2p-kad-dht/network/message"
	"github.com/libp2p/go-libp2p-kad-dht/network/message/ipfskadv1"
	"github.com/libp2p/go-libp2p-kad-dht/network/message/simmessage"
	"github.com/libp2p/go-libp2p-kad-dht/routingtable"
	"github.com/libp2p/go-libp2p-kad-dht/server/simserver"
	"github.com/libp2p/go-libp2p-kad-dht/util"
	"github.com/libp2p/go-libp2p/core/peer"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type SimServer struct {
	rt       routingtable.RoutingTable
	endpoint endpoint.Endpoint

	peerstoreTTL              time.Duration
	numberOfCloserPeersToSend int
}

var _ simserver.SimServer = (*SimServer)(nil)

func NewKadSimServer(rt routingtable.RoutingTable, endpoint endpoint.Endpoint, options ...Option) *SimServer {

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

	var target key.KadKey

	switch msg := msg.(type) {
	case *simmessage.SimMessage:
		t := msg.Target()
		if t == nil {
			// invalid request (nil target), don't reply
			return
		}
		target = *t
	case *ipfskadv1.Message:
		p := peer.ID("")
		if p.UnmarshalBinary(msg.GetKey()) != nil {
			// invalid requested key (not a peer.ID)
			return
		}
		t := peerid.NewPeerID(p)
		target = t.Key()
	default:
		// invalid request (unknown message format), don't reply
		return
	}

	s.endpoint.MaybeAddToPeerstore(ctx, rpeer, s.peerstoreTTL)

	_, span := util.StartSpan(ctx, "SimServer.HandleFindNodeRequest", trace.WithAttributes(
		attribute.Stringer("Requester", rpeer.NodeID()),
		attribute.Stringer("Target", target)))
	defer span.End()

	peers, err := s.rt.NearestPeers(ctx, target, s.numberOfCloserPeersToSend)
	if err != nil {
		span.RecordError(err)
		// invalid request, don't reply
		return
	}

	span.AddEvent("Nearest peers", trace.WithAttributes(
		attribute.Int("count", len(peers)),
	))

	resp := simmessage.NewSimResponse(peers)

	replyFn(resp)
}
