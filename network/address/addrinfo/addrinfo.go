package addrinfo

import (
	"github.com/libp2p/go-libp2p-kad-dht/network/address"
	"github.com/libp2p/go-libp2p-kad-dht/network/address/peerid"
	"github.com/libp2p/go-libp2p/core/peer"
)

type AddrInfo struct {
	peer.AddrInfo
}

var _ address.NetworkAddress = (*AddrInfo)(nil)

func NewAddrInfo(ai peer.AddrInfo) *AddrInfo {
	return &AddrInfo{ai}
}

func (ai AddrInfo) NodeID() address.NodeID {
	return &peerid.PeerID{ID: ai.ID}
}
