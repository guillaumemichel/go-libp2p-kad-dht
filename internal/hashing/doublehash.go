package hashing

import (
	"github.com/multiformats/go-multihash"
	mhreg "github.com/multiformats/go-multihash/core"
)

var (
	// DoubleHashSalt is a prefix prepened to a mulithash when deriving the keyspace location of a provider record
	DoubleHashSalt = []byte("CR_DOUBLEHASH\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00")
	// EncryptionKeySalt is a prefix prepened to a mulithash when deriving the encryption key for a provider record
	EncryptionKeySalt = []byte("CR_ENCRYPTIONKEY\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00")
	// ServerKeySalt is a prefix prepened to a mulithash when deriving the server encryption key for a provider record
	ServerKeySalt = []byte("CR_SERVERKEY\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00")
)

func kadIdFromMultihash(mh multihash.Multihash) []byte {
	// hasher is the hash function used to derive the second hash identifiers
	hasher, _ := mhreg.GetHasher(HasherID)
	// hash the doublehash salt and the multihash
	hasher.Write(append(DoubleHashSalt, mh...))
	return hasher.Sum(nil)
}

func SecondMultihash(mh multihash.Multihash) multihash.Multihash {
	kadid := kadIdFromMultihash(mh)
	// encode the second hash as a multihash using the DBL_SHA2_256 format
	mh2, _ := multihash.Encode(kadid, multihash.DBL_SHA2_256)
	return mh2
}
