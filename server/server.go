package server

import (
	"encoding/hex"
	"fmt"
	"io"

	"github.com/libp2p/go-libp2p-kad-dht/dht/protocol"
	dhtnet "github.com/libp2p/go-libp2p-kad-dht/network"

	"github.com/libp2p/go-libp2p-kad-dht/network/pb"

	"github.com/libp2p/go-libp2p/core/network"
	"google.golang.org/protobuf/proto"
)

type DhtServer struct {
	net *dhtnet.DhtNetwork
}

func NewDhtServer(net *dhtnet.DhtNetwork) *DhtServer {
	server := DhtServer{
		net: net,
	}
	// protocol must be defined in options
	server.net.Host.SetStreamHandler(protocol.ProtocolDHT, server.handleNewStream)

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
	for {
		msg, err := dhtnet.ReadMsg(s)
		if err != nil {
			if err == io.EOF {
				return true
			}
			fmt.Println("error reading message:", err)
			return false
		}
		dhtMsg := &pb.DhtMessage{}
		err = proto.Unmarshal(msg, dhtMsg)
		if err != nil {
			fmt.Println("error unmarshaling message:", err)
			return false
		}
		if dhtMsg.GetProvideRequestType() != nil {
			dht.handleProvideRequest(s, dhtMsg.GetProvideRequestType())
		}
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
