package network

import (
	"encoding/binary"
	"errors"

	"github.com/libp2p/go-libp2p-kad-dht/dht/protocol"
	"github.com/libp2p/go-libp2p-kad-dht/internal/hash"
	"github.com/libp2p/go-libp2p-kad-dht/records"
	"github.com/multiformats/go-multicodec"
	"github.com/multiformats/go-varint"
)

/*
Provider is responsible for sending out provide messages.
*/

const (
	// TODO: replace by actual varint once we have them
	DhtProvideFormatV0              uint16 = 0xf403
	DhtPrefixLookupResponseFormatv0 uint16 = 0xf503
)

func SendProvideMessage(record records.PublishRecord)

func marshallProvideMessage(record records.PublishRecord) []byte {
	// TODO: optimize allocation
	marshalled := make([]byte, 0, 100)
	binary.BigEndian.AppendUint16(marshalled, DhtProvideFormatV0)
	marshalled = append(marshalled, record.ID[:]...)
	marshalled = append(marshalled, record.ServerKey[:]...)
	//marshalled = append(marshalled, record.EncPeerID[:]...)
	marshalled = append(marshalled, record.Signature[:]...)

	return marshalled
}

func unmarsallMessage(message []byte) (records.PublishRecord, error) {
	format, nForamat, err := varint.FromUvarint(message)
	if err != nil {
		return records.PublishRecord{}, err
	}
	if format != uint64(DhtProvideFormatV0) {
		return records.PublishRecord{}, errors.New("unexpected message")
	}

	keyType, nKeyType, err := varint.FromUvarint(message[nForamat:])
	if err != nil {
		return records.PublishRecord{}, err
	}
	if keyType != uint64(multicodec.DblSha2_256) {
		return records.PublishRecord{}, errors.New("unexpected key type")
	}

	id := message[nForamat+nKeyType : nForamat+nKeyType+hash.Keysize]
	serverKey := message[nForamat+nKeyType+hash.Keysize : nForamat+nKeyType+2*hash.Keysize]

	encAlgo, nEnc, err := varint.FromUvarint(message[nForamat+nKeyType+2*hash.Keysize:])
	if err != nil {
		return records.PublishRecord{}, err
	}
	if encAlgo != uint64(protocol.AesGcmMultiCodec) {
		return records.PublishRecord{}, errors.New("unexpected encryption algorithm")
	}
	len, nLen, err := varint.FromUvarint(message[nForamat+nKeyType+2*hash.Keysize+nEnc:])
	if err != nil {
		return records.PublishRecord{}, err
	}
	//encPeerId := message[nForamat+nKeyType+2*hash.Keysize+nEnc+nLen+records.NonceSize : nForamat+nKeyType+2*hash.Keysize+nEnc+nLen+records.NonceSize+int(len)]
	signature := message[nForamat+nKeyType+2*hash.Keysize+nEnc+nLen+records.NonceSize+int(len):]
	return records.PublishRecord{
		ID:        id,
		ServerKey: hash.KadKey(serverKey),
		//EncPeerID: encPeerId,
		Signature: signature,
	}, nil

}
