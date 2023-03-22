package network

import (
	"context"
	"fmt"

	"github.com/libp2p/go-libp2p-kad-dht/dht/protocol"
	"github.com/libp2p/go-libp2p-kad-dht/network/pb"
	"github.com/libp2p/go-libp2p-kad-dht/records"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-msgio"

	"google.golang.org/protobuf/proto"
)

type DhtNetwork struct {
	Host host.Host
}

func NewDhtNetwork(host host.Host) *DhtNetwork {
	return &DhtNetwork{
		Host: host,
	}
}

func (net *DhtNetwork) SendMessage(ctx context.Context, p peer.ID, message []byte) error {
	s, err := net.Host.NewStream(ctx, p, protocol.ProtocolDHT)
	if err != nil {
		return err
	}
	defer s.Close()

	w := msgio.NewWriter(s)
	err = w.WriteMsg(message)
	if err != nil {
		return err
	}
	w.Close()

	return nil
}

func publishRecordToDhtMessage(rec *records.PublishRecord) *pb.DhtMessage {
	encPeerId := pb.EncPeerId{
		EncPeerIdFormatVarint: rec.EncPeerId.EncAlgoVarint,
		Nonce:                 rec.EncPeerId.Nonce[:],
		Payload:               rec.EncPeerId.Payload,
	}
	provMsg := pb.DhtProvideRequest{
		ID:        rec.ID,
		ServerKey: rec.ServerKey[:],
		EncPeerId: &encPeerId,
		Signature: rec.Signature,
	}
	provMsgTyp := pb.DhtMessage_ProvideRequestType{ProvideRequestType: &provMsg}
	return &pb.DhtMessage{MessageType: &provMsgTyp}
}

func (net *DhtNetwork) SendProvide(ctx context.Context, p peer.ID, rec *records.PublishRecord) error {

	msg := publishRecordToDhtMessage(rec)

	s, err := net.Host.NewStream(ctx, p, protocol.ProtocolDHT)
	if err != nil {
		return err
	}
	err = WriteMsg(s, msg)
	if err != nil {
		return err
	}

	bytes, err := ReadMsg(s)
	if err != nil {
		return err
	}
	resp := pb.DhtMessage{}
	err = proto.Unmarshal(bytes, &resp)
	if err != nil {
		return err
	}
	if resp.GetProvideResponseType() == nil {
		fmt.Println("error")
		return nil
	}
	fmt.Println(resp.GetProvideResponseType().Status)

	return nil
}
