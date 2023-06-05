package libp2p

import (
	"github.com/libp2p/go-libp2p-kad-dht/network/address"
	"github.com/libp2p/go-libp2p/core/peer"
)

type Libp2pAddr peer.AddrInfo

func (a Libp2pAddr) NodeID() address.NodeID {
	return a.ID
}
