package records

import (
	"github.com/libp2p/go-libp2p-kad-dht/internal/hash"
	"github.com/multiformats/go-multihash"
)

type PublishRecord struct {
	ID        multihash.Multihash
	ServerKey hash.KadKey
	EncPeerId
	Signature []byte
}

/*
// SerializePublishRecord serializes a PublishRecord into a byte slice
// format: len(message) + format_varint + id + server_key + len(enc_peer_id) + enc_peer_id + timestamp + len(signature) + signature
func SerializePublishRecord(record PublishRecord) []byte {
	encPeerIdLen := len(record.EncPeerID)
	signatureLen := len(record.Signature)
	l := len(protocol.ProvideReqVarint) + len(record.ID) + len(record.ServerKey) + varint.UvarintSize(uint64(encPeerIdLen)) +
		encPeerIdLen + TimestampSize + varint.UvarintSize(uint64(signatureLen)) + signatureLen
	lVarintLen := varint.UvarintSize(uint64(l))

	message := make([]byte, 0, lVarintLen+l)
	message = append(message, varint.ToUvarint(uint64(l))...)
	message = append(message, protocol.ProvideReqVarint...)
	message = append(message, record.ID...)
	message = append(message, record.ServerKey[:]...)
	message = append(message, varint.ToUvarint(uint64(encPeerIdLen))...)
	message = append(message, record.EncPeerID...)
	message = binary.BigEndian.AppendUint32(message, record.Timestamp)
	message = append(message, varint.ToUvarint(uint64(signatureLen))...)
	message = append(message, record.Signature...)

	return message
}

func DeserializePublishRecord(message []byte) (PublishRecord, error) {
	var record PublishRecord
	count := 0

	l, nLen, err := varint.FromUvarint(message[:varint.MaxLenUvarint63])
	if err != nil {
		return record, err
	}
	if len(message) != int(l)+nLen {
		return record, errors.New("wrong length for publish record 1")
	}
	count += nLen

	formatVarint, nFormat, err := varint.FromUvarint(message[count : count+varint.MaxLenUvarint63])
	if err != nil {
		return record, err
	}

	if formatVarint != uint64(protocol.PROVIDE_REQ_MULTICODEC) {
		return record, errors.New("wrong format for publish record")
	}
	count += nFormat

	record.ID = message[count : count+34]
	count += 34
	record.ServerKey = hash.KadKey(message[count : count+hash.Keysize])
	count += hash.Keysize
	encPeerIdLen, nEncPeerIdLen, err := varint.FromUvarint(message[count : count+varint.MaxLenUvarint63])
	if err != nil {
		return record, err
	}
	count += nEncPeerIdLen
	record.EncPeerID = message[count : count+int(encPeerIdLen)]
	count += int(encPeerIdLen)
	record.Timestamp = binary.BigEndian.Uint32(message[count : count+TimestampSize])
	count += TimestampSize
	signatureLen, nSignatureLen, err := varint.FromUvarint(message[count : count+varint.MaxLenUvarint63])
	if err != nil {
		return record, err
	}
	count += nSignatureLen
	record.Signature = message[count : count+int(signatureLen)]
	count += int(signatureLen)

	if count != len(message) {
		return record, errors.New("wrong length for publish record 2")
	}

	return record, nil
}

func SerializePublishRecord2(record PublishRecord) []byte {
	encPeerIdLen := len(record.EncPeerID)
	signatureLen := len(record.Signature)
	l := len(protocol.ProvideReqVarint) + len(record.ID) + len(record.ServerKey) + varint.UvarintSize(uint64(encPeerIdLen)) +
		encPeerIdLen + TimestampSize + varint.UvarintSize(uint64(signatureLen)) + signatureLen

	message := make([]byte, 0, l)
	message = append(message, protocol.ProvideReqVarint...)
	message = append(message, record.ID...)
	message = append(message, record.ServerKey[:]...)
	message = append(message, varint.ToUvarint(uint64(encPeerIdLen))...)
	message = append(message, record.EncPeerID...)
	message = binary.BigEndian.AppendUint32(message, record.Timestamp)
	message = append(message, varint.ToUvarint(uint64(signatureLen))...)
	message = append(message, record.Signature...)

	return message
}

func DeserializePublishRecord2(message []byte) (PublishRecord, error) {
	var record PublishRecord
	count := 0

	formatVarint, nFormat, err := varint.FromUvarint(message[count : count+varint.MaxLenUvarint63])
	if err != nil {
		return record, err
	}

	if formatVarint != uint64(protocol.PROVIDE_REQ_MULTICODEC) {
		return record, errors.New("wrong format for publish record")
	}
	count += nFormat

	record.ID = message[count : count+34]
	count += 34
	record.ServerKey = hash.KadKey(message[count : count+hash.Keysize])
	count += hash.Keysize
	encPeerIdLen, nEncPeerIdLen, err := varint.FromUvarint(message[count : count+varint.MaxLenUvarint63])
	if err != nil {
		return record, err
	}
	count += nEncPeerIdLen
	record.EncPeerID = message[count : count+int(encPeerIdLen)]
	count += int(encPeerIdLen)
	record.Timestamp = binary.BigEndian.Uint32(message[count : count+TimestampSize])
	count += TimestampSize
	signatureLen, nSignatureLen, err := varint.FromUvarint(message[count : count+varint.MaxLenUvarint63])
	if err != nil {
		return record, err
	}
	count += nSignatureLen
	record.Signature = message[count : count+int(signatureLen)]
	count += int(signatureLen)

	if count != len(message) {
		return record, errors.New("wrong length for publish record 2")
	}

	return record, nil
}
*/
