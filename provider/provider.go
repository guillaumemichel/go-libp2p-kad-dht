package provider

import (
	"github.com/ipfs/go-cid"
	"github.com/libp2p/go-libp2p-kad-dht/internal/hash"
	"github.com/libp2p/go-libp2p-kad-dht/records"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/peer"
)

// TODO: move interface
type Provider interface {
	StartProvide(cid.Cid) error
	StopProvide(cid.Cid) error
	ProvideList() []cid.Cid
}

type DhtProvider struct {
	trackList []cid.Cid
	peerid    peer.ID
	privkey   crypto.PrivKey
}

func NewDhtProvider(p peer.ID) *DhtProvider {
	return &DhtProvider{
		trackList: make([]cid.Cid, 0),
		peerid:    p,
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

func (prov *DhtProvider) provide(c cid.Cid) error {
	id := hash.SecondMultihashFromCid(c)
	serverKey := hash.ServerKeyFromCid(c)
	encPeerId, signature, err := records.GetEncPeerId(c, prov.peerid, prov.privkey)
	if err != nil {
		return err
	}

	rec := records.PublishRecord{
		ID:        id,
		ServerKey: serverKey,
		EncPeerID: encPeerId,
		Signature: signature,
	}

	_ = rec

	return nil
}
