package routingtable

import (
	"context"

	"github.com/libp2p/go-libp2p-kad-dht/key"
	"github.com/libp2p/go-libp2p-kad-dht/network/address"
)

type RoutingTable interface {
	Self() key.KadKey
	AddPeer(context.Context, address.NodeID) (bool, error)
	RemoveKey(context.Context, key.KadKey) (bool, error)
	NearestPeers(context.Context, key.KadKey, int) ([]address.NodeID, error)
}

func RemovePeer(ctx context.Context, rt RoutingTable, k address.NodeID) (bool, error) {
	return rt.RemoveKey(ctx, k.Key())
}
