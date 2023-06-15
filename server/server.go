package server

import (
	"context"

	"github.com/libp2p/go-libp2p-kad-dht/events/scheduler"
	message "github.com/libp2p/go-libp2p-kad-dht/network/message/ipfskadv1"
	rt "github.com/libp2p/go-libp2p-kad-dht/routingtable"
	"github.com/libp2p/go-libp2p-kad-dht/util"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"
)

type Server struct {
	ctx          context.Context
	RoutingTable rt.RoutingTable
	host         host.Host
	sched        scheduler.Scheduler

	serverProtocols []protocol.ID
}

func NewServer(ctx context.Context, h host.Host, rt rt.RoutingTable,
	sched scheduler.Scheduler, serverProtocols []protocol.ID) *Server {
	return &Server{
		ctx:             ctx,
		RoutingTable:    rt,
		host:            h,
		sched:           sched,
		serverProtocols: serverProtocols,
	}
}

func HandleRequest(ctx context.Context, s *Server, req *message.Message, stream network.Stream) {
	ctx, span := util.StartSpan(ctx, "HandleRequest")
	defer span.End()

	switch req.GetType() {
	case message.Message_FIND_NODE:
		HandleFindNodeRequest(ctx, s, req, stream)
	default:
		span.AddEvent("unknown request type")
		return
	}

	// TODO: check if remote peer is in server mode. If yes, add them to the routing table
	// if they are in client mode, add them to the table
	p := stream.Conn().RemotePeer()

	if validDhtServer(s, p) {
		rt.AddPeer(ctx, s.RoutingTable, p)
	}
}

func SetStreamHandler(s *Server, handler func(network.Stream), proto protocol.ID) {
	s.host.SetStreamHandler(proto, handler)
}

func validDhtServer(s *Server, p peer.ID) bool {
	proto, err := s.host.Peerstore().FirstSupportedProtocol(p, s.serverProtocols...)
	if err != nil {
		// TODO: log error
		return false
	}

	// if proto is empty, then the peer does not support any of the server protocols
	return proto != ""
}
