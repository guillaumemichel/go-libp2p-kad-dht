package protocol

import (
	"github.com/libp2p/go-libp2p/core/protocol"
	"github.com/multiformats/go-varint"
)

var (
	// ProtocolDHT is the default DHT protocol.
	ProtocolDHT protocol.ID = "/ipfs/kad/1.1.0"
	// DefaultProtocols spoken by the DHT.
	DefaultProtocols = []protocol.ID{ProtocolDHT}
)

const (
	PROVIDE_REQ_MULTICODEC = 0xf403
)

var (
	ProvideReqVarint = varint.ToUvarint(PROVIDE_REQ_MULTICODEC)

	// TODO: replace with multicodec.AesGcm256 once new version gets published
	AesGcmMultiCodec uint64 = 0x2000
)
