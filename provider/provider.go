package provider

import (
	"context"
	"fmt"

	"github.com/ipfs/go-cid"
	"github.com/libp2p/go-libp2p-kad-dht/internal/hash"
	dhtnet "github.com/libp2p/go-libp2p-kad-dht/network"
	"github.com/libp2p/go-libp2p-kad-dht/records"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
)

// TODO: move interface
type Provider interface {
	StartProvide(cid.Cid) error
	StopProvide(cid.Cid) error
	ProvideList() []cid.Cid
	Provide(peer.ID, cid.Cid) error
}

type DhtProvider struct {
	trackList   []cid.Cid
	peerid      peer.ID
	privkey     crypto.PrivKey
	host        host.Host
	msgEndpoint *dhtnet.MessageEndpoint
}

func NewDhtProvider(msgEndpoint *dhtnet.MessageEndpoint) *DhtProvider {
	return &DhtProvider{
		trackList:   make([]cid.Cid, 0),
		peerid:      msgEndpoint.Host.ID(),
		privkey:     msgEndpoint.Host.Peerstore().PrivKey(msgEndpoint.Host.ID()),
		host:        msgEndpoint.Host,
		msgEndpoint: msgEndpoint,
	}
}

func (prov *DhtProvider) StartProvide(cid cid.Cid) error {
	prov.trackList = append(prov.trackList, cid)
	return nil
}

func (prov *DhtProvider) StopProvide(cid cid.Cid) error {
	for i, c := range prov.trackList {
		if c == cid {
			prov.trackList = append(prov.trackList[:i], prov.trackList[i+1:]...)
			return nil
		}
	}
	return nil
}

func (prov *DhtProvider) ProvideList() []cid.Cid {
	return prov.trackList
}

func (prov *DhtProvider) Provide(p peer.ID, c cid.Cid) error {
	fmt.Println("start provide")
	id := hash.SecondMultihashFromCid(c)
	serverKey := hash.ServerKeyFromCid(c)
	encPeerId, signature, err := records.GetEncPeerId(c, prov.peerid, prov.privkey)
	if err != nil {
		return err
	}

	rec := records.PublishRecord{
		ID:        id,
		ServerKey: serverKey,
		EncPeerId: encPeerId,
		Signature: signature,
	}
	err = prov.msgEndpoint.SendProvide(context.Background(), p, &rec)
	if err != nil {
		fmt.Println(err)
	}
	return nil
}
