package server

import (
	"encoding/hex"
	"fmt"
	"io"

	"github.com/libp2p/go-libp2p-kad-dht/dht/protocol"
	"github.com/libp2p/go-libp2p-kad-dht/internal/hash"
	dhtnet "github.com/libp2p/go-libp2p-kad-dht/network"
	rt "github.com/libp2p/go-libp2p-kad-dht/routingtable/simplert"

	"github.com/libp2p/go-libp2p-kad-dht/network/pb"

	"github.com/libp2p/go-libp2p/core/network"
)

type DhtServer struct {
	msgEndpoint *dhtnet.MessageEndpoint
	rt          rt.RoutingTable
}

func NewDhtServer(msgEndpoint *dhtnet.MessageEndpoint, rt rt.RoutingTable) *DhtServer {
	server := DhtServer{
		msgEndpoint: msgEndpoint,
		rt:          rt,
	}
	// protocol must be defined in options
	server.msgEndpoint.Host.SetStreamHandler(protocol.ProtocolDHT, server.handleNewStream)

	return &server
}

// handleNewStream implements the network.StreamHandler
func (dht *DhtServer) handleNewStream(s network.Stream) {
	if dht.handleNewMessage(s) {
		// If we exited without error, close gracefully.
		_ = s.Close()
	} else {
		// otherwise, send an error.
		_ = s.Reset()
	}
}

func (dht *DhtServer) handleNewMessage(s network.Stream) bool {
	fmt.Println("handleNewStream", s)

	rPeer := s.Conn().RemotePeer()

	for {
		req := &pb.DhtMessage{}
		err := dhtnet.ReadMsg(s, req)
		if err != nil {
			if err == io.EOF {
				return true
			}
			fmt.Println("error reading message:", err)
			return false
		}

		// TODO: enhance with https://pkg.go.dev/google.golang.org/protobuf/reflect/protoreflect
		if req.GetProvideRequestType() != nil {
			if !dht.handleProvideRequest(s, req.GetProvideRequestType()) {
				return false
			}
		} else if req.GetFindPeerRequestType() != nil {
			if !dht.handleFindPeer(s, req.GetFindPeerRequestType()) {
				return false
			}
		}
		dht.rt.AddPeer(dht.msgEndpoint.Host.Peerstore().PeerInfo(rPeer))
	}
}

func (dht *DhtServer) handleProvideRequest(s network.Stream, dhtMsg *pb.DhtProvideRequest) bool {
	fmt.Println("got a provide request for", hex.EncodeToString(dhtMsg.ID))
	response := pb.DhtProvideResponse{
		Status: pb.DhtProvideResponse_OK,
	}
	msg := &pb.DhtMessage{
		MessageType: &pb.DhtMessage_ProvideResponseType{
			ProvideResponseType: &response,
		},
	}
	err := dhtnet.WriteMsg(s, msg)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

func (dht *DhtServer) handleFindPeer(s network.Stream, dhtMsg *pb.DhtFindPeerRequest) bool {
	kadid := dhtMsg.GetKadId()
	if kadid == nil || len(kadid) != hash.Keysize {
		return false
	}
	fmt.Println("got a find peer request for", hex.EncodeToString(kadid))

	peers := dht.rt.NearestPeers(hash.KadKey(kadid), dht.rt.BucketSize())
	resp := pb.DhtFindPeerResponse{}
	resp.Peers = make([]*pb.Peer, len(peers))
	for i, p := range peers {
		resp.Peers[i] = &pb.Peer{
			PeerId: p.ID.String(),
			Addrs:  make([][]byte, len(p.Addrs)),
		}
		for j, maddr := range p.Addrs {
			resp.Peers[i].Addrs[j] = maddr.Bytes()
		}
	}

	err := dhtnet.WriteMsg(s, &pb.DhtMessage{
		MessageType: &pb.DhtMessage_FindPeerResponseType{
			FindPeerResponseType: &resp,
		},
	})
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}
