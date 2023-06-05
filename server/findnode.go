package server

import (
	"context"

	"github.com/libp2p/go-libp2p-kad-dht/dht/consts"
	"github.com/libp2p/go-libp2p-kad-dht/internal"
	"github.com/libp2p/go-libp2p-kad-dht/internal/key"
	endpoint "github.com/libp2p/go-libp2p-kad-dht/network/endpoint/libp2pendpoint"
	message "github.com/libp2p/go-libp2p-kad-dht/network/message/ipfskadv1"
	"github.com/libp2p/go-libp2p-kad-dht/network/message/ipfskadv1/pb"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
)

func HandleFindNodeRequest(ctx context.Context, s *Server, req *pb.Message, stream network.Stream) error {

	p := peer.ID("")
	err := p.UnmarshalBinary(req.GetKey())

	if err != nil {
		return err
	}

	_, span := internal.StartSpan(ctx, "server.HandleFindNodeRequest", trace.WithAttributes(
		attribute.Stringer("Requester", stream.Conn().RemotePeer()),
		attribute.Stringer("Target", p)))
	defer span.End()

	peers := s.RoutingTable.NearestPeers(ctx, key.PeerKadID(p), consts.NClosestPeers)

	req.CloserPeers = message.PeeridsToPbPeers(peers, s.host)

	err = endpoint.WriteMsg(stream, req)
	if err != nil {
		return err
	}
	return stream.Close()
}
