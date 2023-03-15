package providerstore

import (
	"errors"
	"sync"

	"github.com/libp2p/go-libp2p-kad-dht/internal/hashing"
	"github.com/libp2p/go-libp2p/core/peer"
)

// The ProviderStore is a temporary structure. It isn't optimized and will be replaced by a proper datastore.

type ProviderStore struct {
	// HASH2 -> ServerKey -> Content Provider peer.ID -> [EncPeerID, TS, Signature]
	// TODO: link to Provider Store spec
	// TODO: replace with better data structure
	store map[hashing.KadKey]map[hashing.KadKey]map[peer.ID][]byte
	lock  sync.RWMutex

	cache map[hashing.KadKey]ProviderRecord // lru cache
}

func NewProviderStore() *ProviderStore {
	return &ProviderStore{
		store: make(map[hashing.KadKey]map[hashing.KadKey]map[peer.ID][]byte),
		cache: make(map[hashing.KadKey]ProviderRecord),
	}
}

func (ps *ProviderStore) AddProvider(key hashing.KadKey, record ProviderRecord) error {
	ps.lock.Lock()
	defer ps.lock.Unlock()

	_, ok := ps.store[key]
	if !ok {
		ps.store[key] = make(map[hashing.KadKey]map[peer.ID][]byte)
	}

	_, ok = ps.store[key][record.ServerKey]
	if !ok {
		ps.store[key][record.ServerKey] = make(map[peer.ID][]byte)
	}
	// TODO: add signature and timestamp
	ps.store[key][record.ServerKey][record.Provider] = record.EncPeerID
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
