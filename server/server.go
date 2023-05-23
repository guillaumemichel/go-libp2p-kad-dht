package server

import (
	"github.com/libp2p/go-libp2p-kad-dht/network/pb"
	rt "github.com/libp2p/go-libp2p-kad-dht/routingtable"

	"github.com/libp2p/go-libp2p/core/network"
)

type Server struct {
	RoutingTable rt.RoutingTable
}

func HandleRequest(s Server, req *pb.DhtMessage, stream network.Stream) error {

	// TODO: enhance with https://pkg.go.dev/google.golang.org/protobuf/reflect/protoreflect
	if req.GetFindPeerRequestType() != nil {

	}

	return nil
}
