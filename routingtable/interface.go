package routingtable

import (
	"context"

	"github.com/libp2p/go-libp2p-kad-dht/internal/key"
	"github.com/libp2p/go-libp2p/core/peer"
)

type RoutingTable interface {
	AddPeer(context.Context, peer.ID) bool
	RemovePeer(context.Context, key.KadKey) bool
	NearestPeers(context.Context, key.KadKey, int) []peer.ID
}

func AddPeer(ctx context.Context, rt RoutingTable, p peer.ID) bool {
	return rt.AddPeer(ctx, p)
}

func RemovePeer(ctx context.Context, rt RoutingTable, k key.KadKey) bool {
	return rt.RemovePeer(ctx, k)
}

func RemovePeerID(ctx context.Context, rt RoutingTable, p peer.ID) bool {
	return RemovePeer(ctx, rt, key.PeerKadID(p))
}

func NearestPeers(ctx context.Context, rt RoutingTable, k key.KadKey, n int) []peer.ID {
	return rt.NearestPeers(ctx, k, n)
}
