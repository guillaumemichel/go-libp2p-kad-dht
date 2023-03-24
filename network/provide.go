package network

import (
	"context"
	"errors"
	"fmt"

	"github.com/libp2p/go-libp2p-kad-dht/network/pb"
	"github.com/libp2p/go-libp2p-kad-dht/records"
	"github.com/libp2p/go-libp2p/core/peer"
)

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

func (msgEndpoint *MessageEndpoint) SendProvide(ctx context.Context, p peer.ID, rec *records.PublishRecord) error {

	req := publishRecordToDhtMessage(rec)
	resp, err := msgEndpoint.SendDhtRequest(ctx, p, req)
	if err != nil {
		return err
	}

	if resp.GetProvideResponseType() == nil {
		return errors.New("not a provide response")
	}
	fmt.Println(resp.GetProvideResponseType().Status)

	return nil
}
