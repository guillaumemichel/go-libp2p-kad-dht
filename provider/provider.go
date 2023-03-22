package provider

import (
	"context"
	"fmt"

	"github.com/ipfs/go-cid"
	"github.com/libp2p/go-libp2p-kad-dht/internal/hash"
	"github.com/libp2p/go-libp2p-kad-dht/network"
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
	trackList []cid.Cid
	peerid    peer.ID
	privkey   crypto.PrivKey
	host      host.Host
	net       *network.DhtNetwork
}

func NewDhtProvider(net *network.DhtNetwork) *DhtProvider {
	return &DhtProvider{
		trackList: make([]cid.Cid, 0),
		peerid:    net.Host.ID(),
		privkey:   net.Host.Peerstore().PrivKey(net.Host.ID()),
		host:      net.Host,
		net:       net,
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
	/*
		marshalled := records.SerializePublishRecord2(rec)

		// find 20 providers
		err = prov.net.SendMessage(context.Background(), p, marshalled)
		if err != nil {
			fmt.Println(err)
		}
	*/
	err = prov.net.SendProvide(context.Background(), p, &rec)
	if err != nil {
		fmt.Println(err)
	}
	return nil
}
