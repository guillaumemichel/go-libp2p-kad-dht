package network

import (
	"context"
	"fmt"

	"github.com/libp2p/go-libp2p-kad-dht/dht/protocol"
	"github.com/libp2p/go-libp2p-kad-dht/internal/hash"
	"github.com/libp2p/go-libp2p-kad-dht/network/pb"
	"github.com/multiformats/go-multiaddr"
	"google.golang.org/protobuf/proto"

	"github.com/libp2p/go-libp2p/core/peer"
)

func (net *DhtNetwork) SendFindPeer(ctx context.Context, p peer.ID, k hash.KadKey) ([]peer.AddrInfo, error) {
	req := &pb.DhtMessage{
		MessageType: &pb.DhtMessage_FindPeerRequestType{
			FindPeerRequestType: &pb.DhtFindPeerRequest{
				KadId: k[:],
			},
		},
	}

	resp, err := net.SendRequest(ctx, p, req)
	if err != nil {
		return nil, err
	}

	if resp.GetFindPeerResponseType() == nil {
		return nil, fmt.Errorf("expected find peer response, got: %v", resp)
	}

	peers := resp.GetFindPeerResponseType().Peers
	return PBPeerToPeerInfos(peers)
}

func (net *DhtNetwork) SendRequest(ctx context.Context, p peer.ID, req *pb.DhtMessage) (*pb.DhtMessage, error) {
	s, err := net.Host.NewStream(ctx, p, protocol.ProtocolDHT)
	if err != nil {
		return nil, err
	}
	defer s.Close()

	err = WriteMsg(s, req)
	if err != nil {
		return nil, err
	}

	rawResp, err := ReadMsg(s)
	if err != nil {
		return nil, err
	}
	resp := &pb.DhtMessage{}
	err = proto.Unmarshal(rawResp, resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func PeerInfosToPBPeer(pis []peer.AddrInfo) []*pb.Peer {
	peers := make([]*pb.Peer, len(pis))
	for i, pi := range pis {
		peers[i] = &pb.Peer{
			PeerId: pi.ID.String(),
			Addrs:  make([][]byte, len(pi.Addrs)),
		}
		for j, maddr := range pi.Addrs {
			peers[i].Addrs[j] = maddr.Bytes()
		}
	}
	return peers
}

func PBPeerToPeerInfos(peers []*pb.Peer) ([]peer.AddrInfo, error) {
	pis := make([]peer.AddrInfo, len(peers))
	var err error
	for i, p := range peers {
		pis[i] = peer.AddrInfo{
			ID:    peer.ID(p.PeerId),
			Addrs: make([]multiaddr.Multiaddr, len(p.Addrs)),
		}
		for j, maddr := range p.Addrs {
			pis[i].Addrs[j], err = multiaddr.NewMultiaddrBytes(maddr)
			if err != nil {
				return nil, err
			}
		}
	}
	return pis, nil
}
