package server

import (
	"github.com/libp2p/go-libp2p-kad-dht/dht/consts"
	"github.com/libp2p/go-libp2p-kad-dht/internal/key"
	"github.com/libp2p/go-libp2p-kad-dht/network"
	"github.com/libp2p/go-libp2p-kad-dht/network/pb"

	"github.com/libp2p/go-libp2p/core/host"
	net "github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
)

func HandleFindNodeRequest(s *Server, req *pb.Message, stream net.Stream) error {

	p := peer.ID(req.GetKey())

	peers := s.RoutingTable.NearestPeers(key.PeerKadID(p), consts.NClosestPeers)

	req.CloserPeers = PeeridsToPbPeers(peers, s.host)

	err := network.WriteMsg(stream, req)
	if err != nil {
		return err
	}
	return stream.Close()
}

func PeeridsToPbPeers(peers []peer.ID, h host.Host) []*pb.Message_Peer {

	pbPeers := make([]*pb.Message_Peer, 0, len(peers))

	for _, p := range peers {
		addrs := h.Peerstore().Addrs(p)
		if len(addrs) == 0 {
			// if no addresses, don't send peer
			continue
		}

		pbAddrs := make([][]byte, len(addrs))
		// convert multiaddresses to bytes
		for i, a := range addrs {
			pbAddrs[i] = a.Bytes()
		}
		pbPeers = append(pbPeers, &pb.Message_Peer{
			Id:         []byte(p),
			Addrs:      pbAddrs,
			Connection: pb.Message_ConnectionType(h.Network().Connectedness(p)),
		})
	}
	return pbPeers
}
