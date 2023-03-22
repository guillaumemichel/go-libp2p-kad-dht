package records

/*
func TestSerializePublishRecord(t *testing.T) {
	require.Equal(t, varint.MaxLenUvarint63, 9)

	rand0 := make([]byte, 32)
	rand.Read(rand0)

	mh, err := multihash.Encode(rand0, multihash.SHA2_256)
	require.NoError(t, err)

	rec := PublishRecord{
		ID:        mh,
		ServerKey: hash.KadKey(rand0),
		EncPeerID: rand0,
		Timestamp: 64,
		Signature: rand0,
	}

	serialized := SerializePublishRecord(rec)
	l, nLen, err := varint.FromUvarint(serialized[:varint.MaxLenUvarint63])
	require.NoError(t, err)
	require.Equal(t, len(serialized), int(l)+nLen)
}

func TestSerializeDeserialisePublishRecord(t *testing.T) {
	rand0 := make([]byte, 32)
	rand.Read(rand0)

	mh, err := multihash.Encode(rand0, multihash.SHA2_256)
	require.NoError(t, err)

	rec := PublishRecord{
		ID:        mh,
		ServerKey: hash.KadKey(rand0),
		EncPeerID: rand0,
		Timestamp: 64,
		Signature: rand0,
	}

	serialized := SerializePublishRecord(rec)
	deserialized, err := DeserializePublishRecord(serialized)
	require.NoError(t, err)
	require.Equal(t, rec, deserialized)
}

func TestSerializeDeserialisePublishRecord2(t *testing.T) {
	rand0 := make([]byte, 32)
	rand.Read(rand0)

	mh, err := multihash.Encode(rand0, multihash.SHA2_256)
	require.NoError(t, err)

	rec := PublishRecord{
		ID:        mh,
		ServerKey: hash.KadKey(rand0),
		EncPeerID: rand0,
		Timestamp: 64,
		Signature: rand0,
	}

	serialized := SerializePublishRecord2(rec)
	deserialized, err := DeserializePublishRecord2(serialized)
	require.NoError(t, err)
	require.Equal(t, rec, deserialized)
}
*/
