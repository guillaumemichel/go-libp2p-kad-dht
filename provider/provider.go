package provider

import (
	"time"

	"github.com/ipfs/go-cid"
	"github.com/libp2p/go-libp2p-kad-dht/internal/hashing"
	"golang.org/x/net/context"
)

type ContentProvider interface {
	StartProviding(cid.Cid) error
	StopProviding(cid.Cid) error
	ProvideList() []cid.Cid

	// Content Router
}

type DhtProvider struct {
	trackList map[cid.Cid]time.Time

	provideFnc func(hashing.KadKey)
}

func NewDhtProvider(ctx context.Context, provideFnc func(hashing.KadKey)) *DhtProvider {
	prov := DhtProvider{
		trackList:  make(map[cid.Cid]time.Time),
		provideFnc: provideFnc,
	}
	prov.run(ctx)
	return &prov
}

func (prov *DhtProvider) StartProviding(c cid.Cid) error {
	prov.trackList[c] = time.Now()
	return nil
}

func (prov *DhtProvider) StopProviding(c cid.Cid) error {
	delete(prov.trackList, c)
	return nil
}

func (prov *DhtProvider) ProvideList() []cid.Cid {
	cids := make([]cid.Cid, 0, len(prov.trackList))

	// We only need the keys
	for cid := range prov.trackList {
		cids = append(cids, cid)
	}
	return cids
}

func (prov *DhtProvider) run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			time.Sleep(time.Second * 10)
			for c := range prov.trackList {
				mh := hashing.SecondMultihash(c.Hash())
				digest := hashing.KadKey(mh[2:])
				prov.provideFnc(digest)
			}
		}
	}
}
