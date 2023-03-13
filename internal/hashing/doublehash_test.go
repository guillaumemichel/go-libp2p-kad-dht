package hashing

import (
	"crypto/sha256"
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/multiformats/go-multihash"
	mhreg "github.com/multiformats/go-multihash/core"
	"github.com/multiformats/go-varint"
)

func TestSalts(t *testing.T) {
	// verify that all salts are 64-bytes long
	saltLen := 64
	assert.Equal(t, len(DoubleHashSalt), saltLen)
	assert.Equal(t, len(EncryptionKeySalt), saltLen)
	assert.Equal(t, len(ServerKeySalt), saltLen)

}

func TestHasher(t *testing.T) {
	// verify that the keysize is 32 bytes = 256-bit
	keysize := 32
	assert.Equal(t, Keysize, keysize)

	assert.Equal(t, multihash.DefaultLengths[HasherID], Keysize)

	// verify that there are no errors getting the hasher
	hasher, err := mhreg.GetHasher(HasherID)
	assert.NoError(t, err)

	// testing the sha256 hasher
	content := []byte("testing double hash")
	// this magic hex string is the result of `echo -n "testing double hash" | sha256sum`
	hash1Check, err := hex.DecodeString("0fc036f2c8508bd0f86ecefc0d97b124eed6a0626b9cd638a67d8a0362b945be")
	assert.NoError(t, err)
	_, err = hasher.Write(content)
	assert.NoError(t, err)
	hash1 := hasher.Sum(nil)
	// assert that the hash computed with the Hasher is the same
	// as the one computed with the sha256sum command
	assert.Equal(t, hash1Check, hash1)

	// verify that the hasher produces the same result as the hasher from the sha256 package
	hasher2 := sha256.New()
	hasher2.Write(content)
	hash2 := hasher2.Sum(nil)
	assert.Equal(t, hash1, hash2)
}

func TestKadIdFromMultihash(t *testing.T) {
	// hex string taken from TestHasher
	digest, err := hex.DecodeString("0fc036f2c8508bd0f86ecefc0d97b124eed6a0626b9cd638a67d8a0362b945be")
	assert.NoError(t, err)
	// SHA256 varint || length of digest varint || digest
	expectedMh := append(append(varint.ToUvarint(multihash.SHA2_256), varint.ToUvarint(uint64(len(digest)))...), digest...)
	mh, err := multihash.Encode(digest, HasherID)
	assert.NoError(t, err)
	assert.Equal(t, expectedMh, mh)

	// verify that the multihash is correctly converted to a KadId
	digest2 := kadIdFromMultihash(mh)

	// get sha256 hasher
	hasher, err := mhreg.GetHasher(HasherID)
	assert.NoError(t, err)
	// write the salt and the multihash to the hasher
	_, err = hasher.Write(DoubleHashSalt)
	assert.NoError(t, err)
	_, err = hasher.Write(mh)
	assert.NoError(t, err)
	// get the digest
	exprectedDigest2 := hasher.Sum(nil)
	// verify that the digest is the same as the multihash
	assert.Equal(t, exprectedDigest2, digest2)
}

func TestSecondMultihash(t *testing.T) {
	// hex string taken from TestHasher
	digest, err := hex.DecodeString("0fc036f2c8508bd0f86ecefc0d97b124eed6a0626b9cd638a67d8a0362b945be")
	assert.NoError(t, err)
	mh, err := multihash.Encode(digest, HasherID)
	assert.NoError(t, err)

	// verify that the multihash is correctly converted to a KadId
	mh2 := SecondMultihash(mh)

	// get sha256 hasher
	hasher, err := mhreg.GetHasher(HasherID)
	assert.NoError(t, err)
	// write the salt and the multihash to the hasher
	_, err = hasher.Write(DoubleHashSalt)
	assert.NoError(t, err)
	_, err = hasher.Write(mh)
	assert.NoError(t, err)
	// get the digest
	digest2 := hasher.Sum(nil)
	// get the expected second multihash
	buf := append(append(varint.ToUvarint(multihash.DBL_SHA2_256), varint.ToUvarint(uint64(len(digest2)))...), digest2...)
	// cast from bytes to multihash
	expectedMh2, err := multihash.Cast(buf)
	assert.NoError(t, err)
	// verify that the second multihash is the same as the one computed
	assert.Equal(t, expectedMh2, mh2[:])

}
