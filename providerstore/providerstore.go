package providerstore

import (
	"errors"
	"sync"

	"github.com/libp2p/go-libp2p-kad-dht/internal/hashing"
	"github.com/libp2p/go-libp2p/core/peer"
)

type ProviderStore struct {
	// HASH2 -> ServerKey -> Content Provider peer.ID -> [EncPeerID, TS, Signature]
	// TODO: link to Provider Store spec
	store map[hashing.KadKey]map[hashing.KadKey]map[peer.ID][]byte
	lock  sync.RWMutex
}

func NewProviderStore() *ProviderStore {
	return &ProviderStore{
		//store: make(map[hashing.KadKey]map[hashing.KadKey]map[peer.ID][]byte),
	}
}

func (ps *ProviderStore) AddProvider(k hashing.KadKey, val []byte) error {
	ps.lock.Lock()
	defer ps.lock.Unlock()

	_, ok := ps.store[k]
	if !ok {
		ps.store[k] = make(map[hashing.KadKey]map[peer.ID][]byte)
	}

	_, ok = ps.store[k][k]
	if !ok {
		ps.store[k][k] = make(map[peer.ID][]byte)
	}

	ps.store[k][k]["peerID"] = val
	return nil
}

func (ps *ProviderStore) GetProviders(k hashing.KadKey) ([]byte, error) {
	ps.lock.RLock()
	defer ps.lock.RUnlock()

	if _, ok := ps.store[k]; !ok {
		return nil, errors.New("not found")
	}

	return ps.store[k][k]["peerID"], nil
}

func (ps *ProviderStore) RemoveProvider(k hashing.KadKey) error {
	ps.lock.Lock()
	defer ps.lock.Unlock()

	if _, ok := ps.store[k]; !ok {
		return nil
	}

	delete(ps.store[k][k], "peerID")
	return nil
}
