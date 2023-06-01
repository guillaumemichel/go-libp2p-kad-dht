package network

import (
	"context"
	"errors"

	"github.com/libp2p/go-libp2p-kad-dht/internal"
	"github.com/libp2p/go-libp2p-kad-dht/network/pb"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
)

var (
	ErrNoValidAddresses = errors.New("no valid addresses")
)

func PBPeerToPeerInfo(ctx context.Context, pbp *pb.Message_Peer) (peer.AddrInfo, error) {
	_, span := internal.StartSpan(ctx, "network.PBPeerToPeerInfo")
	defer span.End()

	addrs := make([]multiaddr.Multiaddr, 0, len(pbp.Addrs))
	for _, a := range pbp.Addrs {
		addr, err := multiaddr.NewMultiaddrBytes(a)
		if err == nil {
			addrs = append(addrs, addr)
		} else {
			span.RecordError(err)
		}
	}
	if len(addrs) == 0 {
		span.RecordError(ErrNoValidAddresses)
		return peer.AddrInfo{}, ErrNoValidAddresses
	}

	return peer.AddrInfo{
		ID:    peer.ID(pbp.Id),
		Addrs: addrs,
	}, nil
}

func ParsePeers(ctx context.Context, pbps []*pb.Message_Peer) []peer.AddrInfo {

	peers := make([]peer.AddrInfo, 0, len(pbps))
	for _, p := range pbps {
		pi, err := PBPeerToPeerInfo(ctx, p)
		if err == nil {
			peers = append(peers, pi)
		}
	}
	return peers
}
