package routing

import (
	"context"

	"github.com/libp2p/go-libp2p-kad-dht/internal/hash"
	"github.com/libp2p/go-libp2p-kad-dht/network/pb"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/routing"
)

func (r *DhtRouting) FindPeer(ctx context.Context, p peer.ID) (peer.AddrInfo, error) {
	// Test is provided peer.ID is valid
	if err := p.Validate(); err != nil {
		return peer.AddrInfo{}, err
	}

	// Check if we are already connected to them
	if addrInfo := r.FindLocal(ctx, p); addrInfo.ID != "" {
		return addrInfo, nil
	}

	kadId := hash.PeerKadID(p)
	req := &pb.DhtMessage{
		MessageType: &pb.DhtMessage_FindPeerRequestType{
			FindPeerRequestType: &pb.DhtFindPeerRequest{KadId: kadId[:]},
		},
	}
	r.qManager.Query(ctx, kadId, req)

	return peer.AddrInfo{}, routing.ErrNotFound
}

// FindLocal looks for a peer with a given ID connected to this dht and returns
// the peer and the table it was found in.
func (r *DhtRouting) FindLocal(ctx context.Context, id peer.ID) peer.AddrInfo {
	switch r.me.Host.Network().Connectedness(id) {
	case network.Connected, network.CanConnect:
		return r.me.Host.Peerstore().PeerInfo(id)
	default:
		return peer.AddrInfo{}
	}
}
