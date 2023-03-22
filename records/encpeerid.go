package records

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/binary"
	"time"

	"github.com/ipfs/go-cid"
	"github.com/libp2p/go-libp2p-kad-dht/dht/protocol"
	"github.com/libp2p/go-libp2p-kad-dht/internal/hash"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-varint"
)

const (
	NonceSize     = 12 // in bytes
	TimestampSize = 4  // in bytes
)

var (
	codecSize = varint.UvarintSize(uint64(protocol.AesGcmMultiCodec))
)

type EncPeerId struct {
	EncAlgoVarint []byte
	Nonce         [NonceSize]byte
	Payload       []byte
}

func (encPeerId *EncPeerId) Timestamp() uint32 {
	return binary.BigEndian.Uint32(encPeerId.Nonce[:TimestampSize])
}

func GetEncPeerId(c cid.Cid, p peer.ID, privkey crypto.PrivKey) (EncPeerId, []byte, error) {
	encKey := hash.RecordEncryptionKeyFromCid(c)
	block, err := aes.NewCipher(encKey[:])
	if err != nil {
		return EncPeerId{}, nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return EncPeerId{}, nil, err
	}

	timestamp := CurrentTimestamp()

	var nonce [NonceSize]byte
	copy(nonce[:], binary.BigEndian.AppendUint32(nil, timestamp))
	n, err := rand.Read(nonce[TimestampSize:])
	if err != nil || n != NonceSize-TimestampSize {
		return EncPeerId{}, nil, err
	}
	payload := gcm.Seal(nil, nonce[:], []byte(p), nil)

	encPeerId := EncPeerId{
		EncAlgoVarint: varint.ToUvarint(protocol.AesGcmMultiCodec),
		Nonce:         nonce,
		Payload:       payload,
	}

	toSign := append(payload, nonce[:TimestampSize]...)
	signature, err := privkey.Sign(toSign)
	if err != nil {
		return EncPeerId{}, nil, err
	}
	return encPeerId, signature, nil
}

func CurrentTimestamp() uint32 {
	return uint32(time.Now().Unix() / 60)
}

func GetEncPeerId1(c cid.Cid, p peer.ID, privkey crypto.PrivKey) ([]byte, []byte, error) {
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
	copy(encPeerId, varint.ToUvarint(uint64(protocol.AesGcmMultiCodec)))
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
