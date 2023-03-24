package network

import (
	"context"
	"fmt"

	"github.com/libp2p/go-libp2p-kad-dht/internal/hash"
	"github.com/libp2p/go-libp2p-kad-dht/network/pb"

	"github.com/libp2p/go-libp2p/core/peer"
)

func (me *MessageEndpoint) SendFindPeer(ctx context.Context, p peer.ID, k hash.KadKey) ([]peer.AddrInfo, error) {
	req := &pb.DhtMessage{
		MessageType: &pb.DhtMessage_FindPeerRequestType{
			FindPeerRequestType: &pb.DhtFindPeerRequest{
				KadId: k[:],
			},
		},
	}

	resp, err := me.SendDhtRequest(ctx, p, req)
	if err != nil {
		return nil, err
	}

	if resp.GetFindPeerResponseType() == nil {
		return nil, fmt.Errorf("expected find peer response, got: %v", resp)
	}

	peers := resp.GetFindPeerResponseType().Peers
	return PBPeerToPeerInfos(peers)
}
