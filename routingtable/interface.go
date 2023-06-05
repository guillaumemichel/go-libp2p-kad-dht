package routingtable

import (
	"context"

	"github.com/libp2p/go-libp2p-kad-dht/internal/key"
	"github.com/libp2p/go-libp2p-kad-dht/network/address"
)

type RoutingTable interface {
	AddPeer(context.Context, address.NodeID) bool
	RemovePeer(context.Context, key.KadKey) bool
	NearestPeers(context.Context, key.KadKey, int) []address.NodeID
}

func AddPeer(ctx context.Context, rt RoutingTable, p address.NodeID) bool {
	return rt.AddPeer(ctx, p)
}

func RemovePeer(ctx context.Context, rt RoutingTable, k key.KadKey) bool {
	return rt.RemovePeer(ctx, k)
}

func RemovePeerID(ctx context.Context, rt RoutingTable, p address.NodeID) bool {
	return RemovePeer(ctx, rt, address.KadID(p))
}

func NearestPeers(ctx context.Context, rt RoutingTable, k key.KadKey, n int) []address.NodeID {
	return rt.NearestPeers(ctx, k, n)
}
