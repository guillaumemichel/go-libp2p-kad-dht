package server

import (
	"errors"

	"github.com/libp2p/go-libp2p-kad-dht/dht/consts"
	"github.com/libp2p/go-libp2p-kad-dht/internal/key"
	"github.com/libp2p/go-libp2p-kad-dht/network/pb"
	"github.com/libp2p/go-libp2p/core/network"
)

func HandleFindPeerRequest(s Server, req *pb.DhtFindPeerRequest, stream network.Stream) error {

	kadid := req.GetKadId()
	if kadid == nil || len(kadid) != key.Keysize {
		return errors.New("invalid kadid")
	}

	peers := s.RoutingTable.NearestPeers(key.KadKey(kadid), consts.NClosestPeers)

	_ = peers

	return nil
}
