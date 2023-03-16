package records

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/binary"
	"time"

	"github.com/ipfs/go-cid"
	"github.com/libp2p/go-libp2p-kad-dht/internal/hash"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-varint"
)

var (
	NonceSize     = 12 // in bytes
	TimestampSize = 4  // in bytes

	// TODO: replace with multicodec.AesGcm256 once new version gets published
	AesGcmMultiCodec = 0x2000
	codecSize        = varint.UvarintSize(uint64(AesGcmMultiCodec))
)

func CurrentTimestamp() uint32 {
	return uint32(time.Now().Unix() / 60)
}

func GetEncPeerId(c cid.Cid, p peer.ID, privkey crypto.PrivKey) ([]byte, []byte, error) {
	encKey := hash.RecordEncryptionKeyFromCid(c)
	block, err := aes.NewCipher(encKey[:])
	if err != nil {
		return nil, nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, nil, err
	}

	timestamp := CurrentTimestamp()

	nonce := make([]byte, NonceSize)
	binary.BigEndian.AppendUint32(nonce, timestamp)
	n, err := rand.Read(nonce[TimestampSize:])
	if err != nil || n != NonceSize-TimestampSize {
		return nil, nil, err
	}
	payload := gcm.Seal(nil, nonce, []byte(p), nil)

	payloadLen := len(payload)
	payloadLenSize := varint.UvarintSize(uint64(payloadLen))

	// TODO: optimize this
	encPeerId := make([]byte, codecSize+payloadLenSize+NonceSize+payloadLen)
	copy(encPeerId, varint.ToUvarint(uint64(AesGcmMultiCodec)))
	copy(encPeerId[codecSize:], varint.ToUvarint(uint64(payloadLen)))
	copy(encPeerId[codecSize+payloadLenSize:], nonce)
	copy(encPeerId[codecSize+payloadLenSize+NonceSize:], payload)

	binary.BigEndian.AppendUint32(payload, timestamp)
	signature, err := privkey.Sign(payload)
	if err != nil {
		return nil, nil, err
	}

	return encPeerId, signature, nil
}
