package network

import (
	"github.com/libp2p/go-libp2p-kad-dht/network/pb"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
)

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
