package routing

import (
	"context"

	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/routing"
)

type PeerRouting interface {
	FindPeer(context.Context, peer.ID) (peer.AddrInfo, error)
}

type ContentRouting interface {
}

type DhtRouting struct {
	host host.Host
}

func NewDhtRouting(host host.Host) *DhtRouting {
	return &DhtRouting{
		host: host,
	}
}

func (r *DhtRouting) FindPeer(ctx context.Context, p peer.ID) (peer.AddrInfo, error) {
	// Test is provided peer.ID is valid
	if err := p.Validate(); err != nil {
		return peer.AddrInfo{}, err
	}

	// Check if we are already connected to them
	if addrInfo := r.FindLocal(p); addrInfo.ID != "" {
		return addrInfo, nil
	}

	return peer.AddrInfo{}, routing.ErrNotFound
}

// FindLocal looks for a peer with a given ID connected to this dht and returns the peer and the table it was found in.
func (r *DhtRouting) FindLocal(id peer.ID) peer.AddrInfo {
	switch r.host.Network().Connectedness(id) {
	case network.Connected, network.CanConnect:
		return r.host.Peerstore().PeerInfo(id)
	default:
		return peer.AddrInfo{}
	}
}
