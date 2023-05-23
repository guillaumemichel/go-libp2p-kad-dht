package server

import (
	"github.com/libp2p/go-libp2p-kad-dht/events"
	"github.com/libp2p/go-libp2p-kad-dht/network/pb"
	rt "github.com/libp2p/go-libp2p-kad-dht/routingtable"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/protocol"
)

type Server struct {
	RoutingTable rt.RoutingTable
	host         host.Host
	em           *events.EventsManager
}

func HandleRequest(s Server, req *pb.Message, stream network.Stream) error {

	switch req.GetType() {
	case pb.Message_FIND_NODE:
		HandleFindNodeRequest(s, req, stream)
	default:
	}

	// TODO: check if remote peer is in server mode. If yes, add them to the routing table

	return nil
}

func SetStreamHandler(s Server, handler func(network.Stream), proto protocol.ID) {
	s.host.SetStreamHandler(proto, handler)
}
