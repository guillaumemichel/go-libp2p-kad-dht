// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.29.1
// 	protoc        v3.21.12
// source: message.proto

package pb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type DhtProvideResponseStatus int32

const (
	DhtProvideResponse_OK    DhtProvideResponseStatus = 0
	DhtProvideResponse_ERROR DhtProvideResponseStatus = 1
)

// Enum value maps for DhtProvideResponseStatus.
var (
	DhtProvideResponseStatus_name = map[int32]string{
		0: "OK",
		1: "ERROR",
	}
	DhtProvideResponseStatus_value = map[string]int32{
		"OK":    0,
		"ERROR": 1,
	}
)

func (x DhtProvideResponseStatus) Enum() *DhtProvideResponseStatus {
	p := new(DhtProvideResponseStatus)
	*p = x
	return p
}

func (x DhtProvideResponseStatus) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (DhtProvideResponseStatus) Descriptor() protoreflect.EnumDescriptor {
	return file_message_proto_enumTypes[0].Descriptor()
}

func (DhtProvideResponseStatus) Type() protoreflect.EnumType {
	return &file_message_proto_enumTypes[0]
}

func (x DhtProvideResponseStatus) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use DhtProvideResponseStatus.Descriptor instead.
func (DhtProvideResponseStatus) EnumDescriptor() ([]byte, []int) {
	return file_message_proto_rawDescGZIP(), []int{2, 0}
}

type DhtLookupResponse_Peer_ConnectionType int32

const (
	// sender does not have a connection to peer, and no extra information
	// (default)
	DhtLookupResponse_Peer_NOT_CONNECTED DhtLookupResponse_Peer_ConnectionType = 0
	// sender has a live connection to peer
	DhtLookupResponse_Peer_CONNECTED DhtLookupResponse_Peer_ConnectionType = 1
	// sender recently connected to peer
	DhtLookupResponse_Peer_CAN_CONNECT DhtLookupResponse_Peer_ConnectionType = 2
	// sender recently tried to connect to peer repeatedly but failed to
	// connect
	// ("try" here is loose, but this should signal "made strong effort,
	// failed")
	DhtLookupResponse_Peer_CANNOT_CONNECT DhtLookupResponse_Peer_ConnectionType = 3
)

// Enum value maps for DhtLookupResponse_Peer_ConnectionType.
var (
	DhtLookupResponse_Peer_ConnectionType_name = map[int32]string{
		0: "NOT_CONNECTED",
		1: "CONNECTED",
		2: "CAN_CONNECT",
		3: "CANNOT_CONNECT",
	}
	DhtLookupResponse_Peer_ConnectionType_value = map[string]int32{
		"NOT_CONNECTED":  0,
		"CONNECTED":      1,
		"CAN_CONNECT":    2,
		"CANNOT_CONNECT": 3,
	}
)

func (x DhtLookupResponse_Peer_ConnectionType) Enum() *DhtLookupResponse_Peer_ConnectionType {
	p := new(DhtLookupResponse_Peer_ConnectionType)
	*p = x
	return p
}

func (x DhtLookupResponse_Peer_ConnectionType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (DhtLookupResponse_Peer_ConnectionType) Descriptor() protoreflect.EnumDescriptor {
	return file_message_proto_enumTypes[1].Descriptor()
}

func (DhtLookupResponse_Peer_ConnectionType) Type() protoreflect.EnumType {
	return &file_message_proto_enumTypes[1]
}

func (x DhtLookupResponse_Peer_ConnectionType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use DhtLookupResponse_Peer_ConnectionType.Descriptor instead.
func (DhtLookupResponse_Peer_ConnectionType) EnumDescriptor() ([]byte, []int) {
	return file_message_proto_rawDescGZIP(), []int{4, 0, 0}
}

type EncPeerId struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	EncPeerIdFormatVarint []byte `protobuf:"bytes,1,opt,name=EncPeerIdFormatVarint,proto3" json:"EncPeerIdFormatVarint,omitempty"`
	Nonce                 []byte `protobuf:"bytes,2,opt,name=Nonce,proto3" json:"Nonce,omitempty"`
	Payload               []byte `protobuf:"bytes,3,opt,name=Payload,proto3" json:"Payload,omitempty"`
}

func (x *EncPeerId) Reset() {
	*x = EncPeerId{}
	if protoimpl.UnsafeEnabled {
		mi := &file_message_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EncPeerId) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EncPeerId) ProtoMessage() {}

func (x *EncPeerId) ProtoReflect() protoreflect.Message {
	mi := &file_message_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EncPeerId.ProtoReflect.Descriptor instead.
func (*EncPeerId) Descriptor() ([]byte, []int) {
	return file_message_proto_rawDescGZIP(), []int{0}
}

func (x *EncPeerId) GetEncPeerIdFormatVarint() []byte {
	if x != nil {
		return x.EncPeerIdFormatVarint
	}
	return nil
}

func (x *EncPeerId) GetNonce() []byte {
	if x != nil {
		return x.Nonce
	}
	return nil
}

func (x *EncPeerId) GetPayload() []byte {
	if x != nil {
		return x.Payload
	}
	return nil
}

type DhtProvideRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ID        []byte     `protobuf:"bytes,1,opt,name=ID,proto3" json:"ID,omitempty"`
	ServerKey []byte     `protobuf:"bytes,2,opt,name=ServerKey,proto3" json:"ServerKey,omitempty"`
	EncPeerId *EncPeerId `protobuf:"bytes,3,opt,name=EncPeerId,proto3" json:"EncPeerId,omitempty"`
	Signature []byte     `protobuf:"bytes,4,opt,name=Signature,proto3" json:"Signature,omitempty"`
}

func (x *DhtProvideRequest) Reset() {
	*x = DhtProvideRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_message_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DhtProvideRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DhtProvideRequest) ProtoMessage() {}

func (x *DhtProvideRequest) ProtoReflect() protoreflect.Message {
	mi := &file_message_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DhtProvideRequest.ProtoReflect.Descriptor instead.
func (*DhtProvideRequest) Descriptor() ([]byte, []int) {
	return file_message_proto_rawDescGZIP(), []int{1}
}

func (x *DhtProvideRequest) GetID() []byte {
	if x != nil {
		return x.ID
	}
	return nil
}

func (x *DhtProvideRequest) GetServerKey() []byte {
	if x != nil {
		return x.ServerKey
	}
	return nil
}

func (x *DhtProvideRequest) GetEncPeerId() *EncPeerId {
	if x != nil {
		return x.EncPeerId
	}
	return nil
}

func (x *DhtProvideRequest) GetSignature() []byte {
	if x != nil {
		return x.Signature
	}
	return nil
}

type DhtProvideResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Status DhtProvideResponseStatus `protobuf:"varint,1,opt,name=Status,proto3,enum=DhtProvideResponseStatus" json:"Status,omitempty"`
}

func (x *DhtProvideResponse) Reset() {
	*x = DhtProvideResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_message_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DhtProvideResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DhtProvideResponse) ProtoMessage() {}

func (x *DhtProvideResponse) ProtoReflect() protoreflect.Message {
	mi := &file_message_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DhtProvideResponse.ProtoReflect.Descriptor instead.
func (*DhtProvideResponse) Descriptor() ([]byte, []int) {
	return file_message_proto_rawDescGZIP(), []int{2}
}

func (x *DhtProvideResponse) GetStatus() DhtProvideResponseStatus {
	if x != nil {
		return x.Status
	}
	return DhtProvideResponse_OK
}

type DhtLookupRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Prefix []byte `protobuf:"bytes,1,opt,name=Prefix,proto3" json:"Prefix,omitempty"`
}

func (x *DhtLookupRequest) Reset() {
	*x = DhtLookupRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_message_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DhtLookupRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DhtLookupRequest) ProtoMessage() {}

func (x *DhtLookupRequest) ProtoReflect() protoreflect.Message {
	mi := &file_message_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DhtLookupRequest.ProtoReflect.Descriptor instead.
func (*DhtLookupRequest) Descriptor() ([]byte, []int) {
	return file_message_proto_rawDescGZIP(), []int{3}
}

func (x *DhtLookupRequest) GetPrefix() []byte {
	if x != nil {
		return x.Prefix
	}
	return nil
}

type DhtLookupResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	DhtLookupResponseFormatVarint []byte                                         `protobuf:"bytes,1,opt,name=DhtLookupResponseFormatVarint,proto3" json:"DhtLookupResponseFormatVarint,omitempty"`
	Flag                          []byte                                         `protobuf:"bytes,2,opt,name=Flag,proto3" json:"Flag,omitempty"`
	Peers                         []*DhtLookupResponse_Peer                      `protobuf:"bytes,3,rep,name=Peers,proto3" json:"Peers,omitempty"`
	ProviderRecords               []*DhtLookupResponse_AggregatedProviderRecords `protobuf:"bytes,4,rep,name=ProviderRecords,proto3" json:"ProviderRecords,omitempty"`
}

func (x *DhtLookupResponse) Reset() {
	*x = DhtLookupResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_message_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DhtLookupResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DhtLookupResponse) ProtoMessage() {}

func (x *DhtLookupResponse) ProtoReflect() protoreflect.Message {
	mi := &file_message_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DhtLookupResponse.ProtoReflect.Descriptor instead.
func (*DhtLookupResponse) Descriptor() ([]byte, []int) {
	return file_message_proto_rawDescGZIP(), []int{4}
}

func (x *DhtLookupResponse) GetDhtLookupResponseFormatVarint() []byte {
	if x != nil {
		return x.DhtLookupResponseFormatVarint
	}
	return nil
}

func (x *DhtLookupResponse) GetFlag() []byte {
	if x != nil {
		return x.Flag
	}
	return nil
}

func (x *DhtLookupResponse) GetPeers() []*DhtLookupResponse_Peer {
	if x != nil {
		return x.Peers
	}
	return nil
}

func (x *DhtLookupResponse) GetProviderRecords() []*DhtLookupResponse_AggregatedProviderRecords {
	if x != nil {
		return x.ProviderRecords
	}
	return nil
}

type DhtMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to MessageType:
	//
	//	*DhtMessage_ProvideRequestType
	//	*DhtMessage_ProvideResponseType
	//	*DhtMessage_LookupRequestType
	//	*DhtMessage_LookupResponseType
	MessageType isDhtMessage_MessageType `protobuf_oneof:"MessageType"`
}

func (x *DhtMessage) Reset() {
	*x = DhtMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_message_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DhtMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DhtMessage) ProtoMessage() {}

func (x *DhtMessage) ProtoReflect() protoreflect.Message {
	mi := &file_message_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DhtMessage.ProtoReflect.Descriptor instead.
func (*DhtMessage) Descriptor() ([]byte, []int) {
	return file_message_proto_rawDescGZIP(), []int{5}
}

func (m *DhtMessage) GetMessageType() isDhtMessage_MessageType {
	if m != nil {
		return m.MessageType
	}
	return nil
}

func (x *DhtMessage) GetProvideRequestType() *DhtProvideRequest {
	if x, ok := x.GetMessageType().(*DhtMessage_ProvideRequestType); ok {
		return x.ProvideRequestType
	}
	return nil
}

func (x *DhtMessage) GetProvideResponseType() *DhtProvideResponse {
	if x, ok := x.GetMessageType().(*DhtMessage_ProvideResponseType); ok {
		return x.ProvideResponseType
	}
	return nil
}

func (x *DhtMessage) GetLookupRequestType() *DhtLookupRequest {
	if x, ok := x.GetMessageType().(*DhtMessage_LookupRequestType); ok {
		return x.LookupRequestType
	}
	return nil
}

func (x *DhtMessage) GetLookupResponseType() *DhtLookupResponse {
	if x, ok := x.GetMessageType().(*DhtMessage_LookupResponseType); ok {
		return x.LookupResponseType
	}
	return nil
}

type isDhtMessage_MessageType interface {
	isDhtMessage_MessageType()
}

type DhtMessage_ProvideRequestType struct {
	ProvideRequestType *DhtProvideRequest `protobuf:"bytes,1,opt,name=ProvideRequestType,proto3,oneof"`
}

type DhtMessage_ProvideResponseType struct {
	ProvideResponseType *DhtProvideResponse `protobuf:"bytes,2,opt,name=ProvideResponseType,proto3,oneof"`
}

type DhtMessage_LookupRequestType struct {
	LookupRequestType *DhtLookupRequest `protobuf:"bytes,3,opt,name=LookupRequestType,proto3,oneof"`
}

type DhtMessage_LookupResponseType struct {
	LookupResponseType *DhtLookupResponse `protobuf:"bytes,4,opt,name=LookupResponseType,proto3,oneof"`
}

func (*DhtMessage_ProvideRequestType) isDhtMessage_MessageType() {}

func (*DhtMessage_ProvideResponseType) isDhtMessage_MessageType() {}

func (*DhtMessage_LookupRequestType) isDhtMessage_MessageType() {}

func (*DhtMessage_LookupResponseType) isDhtMessage_MessageType() {}

type DhtLookupResponse_Peer struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id             []byte                                `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Addrs          [][]byte                              `protobuf:"bytes,2,rep,name=addrs,proto3" json:"addrs,omitempty"`
	ConnectionType DhtLookupResponse_Peer_ConnectionType `protobuf:"varint,3,opt,name=connectionType,proto3,enum=DhtLookupResponse_Peer_ConnectionType" json:"connectionType,omitempty"`
}

func (x *DhtLookupResponse_Peer) Reset() {
	*x = DhtLookupResponse_Peer{}
	if protoimpl.UnsafeEnabled {
		mi := &file_message_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DhtLookupResponse_Peer) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DhtLookupResponse_Peer) ProtoMessage() {}

func (x *DhtLookupResponse_Peer) ProtoReflect() protoreflect.Message {
	mi := &file_message_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DhtLookupResponse_Peer.ProtoReflect.Descriptor instead.
func (*DhtLookupResponse_Peer) Descriptor() ([]byte, []int) {
	return file_message_proto_rawDescGZIP(), []int{4, 0}
}

func (x *DhtLookupResponse_Peer) GetId() []byte {
	if x != nil {
		return x.Id
	}
	return nil
}

func (x *DhtLookupResponse_Peer) GetAddrs() [][]byte {
	if x != nil {
		return x.Addrs
	}
	return nil
}

func (x *DhtLookupResponse_Peer) GetConnectionType() DhtLookupResponse_Peer_ConnectionType {
	if x != nil {
		return x.ConnectionType
	}
	return DhtLookupResponse_Peer_NOT_CONNECTED
}

type DhtLookupResponse_ProviderRecord struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	FormatVarint []byte     `protobuf:"bytes,1,opt,name=FormatVarint,proto3" json:"FormatVarint,omitempty"`
	EncPeerId    *EncPeerId `protobuf:"bytes,2,opt,name=EncPeerId,proto3" json:"EncPeerId,omitempty"`
	ServerNonce  []byte     `protobuf:"bytes,3,opt,name=ServerNonce,proto3" json:"ServerNonce,omitempty"`
	EncMetadata  []byte     `protobuf:"bytes,4,opt,name=EncMetadata,proto3" json:"EncMetadata,omitempty"`
}

func (x *DhtLookupResponse_ProviderRecord) Reset() {
	*x = DhtLookupResponse_ProviderRecord{}
	if protoimpl.UnsafeEnabled {
		mi := &file_message_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DhtLookupResponse_ProviderRecord) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DhtLookupResponse_ProviderRecord) ProtoMessage() {}

func (x *DhtLookupResponse_ProviderRecord) ProtoReflect() protoreflect.Message {
	mi := &file_message_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DhtLookupResponse_ProviderRecord.ProtoReflect.Descriptor instead.
func (*DhtLookupResponse_ProviderRecord) Descriptor() ([]byte, []int) {
	return file_message_proto_rawDescGZIP(), []int{4, 1}
}

func (x *DhtLookupResponse_ProviderRecord) GetFormatVarint() []byte {
	if x != nil {
		return x.FormatVarint
	}
	return nil
}

func (x *DhtLookupResponse_ProviderRecord) GetEncPeerId() *EncPeerId {
	if x != nil {
		return x.EncPeerId
	}
	return nil
}

func (x *DhtLookupResponse_ProviderRecord) GetServerNonce() []byte {
	if x != nil {
		return x.ServerNonce
	}
	return nil
}

func (x *DhtLookupResponse_ProviderRecord) GetEncMetadata() []byte {
	if x != nil {
		return x.EncMetadata
	}
	return nil
}

type DhtLookupResponse_AggregatedProviderRecords struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id              []byte                              `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	ProviderRecords []*DhtLookupResponse_ProviderRecord `protobuf:"bytes,2,rep,name=ProviderRecords,proto3" json:"ProviderRecords,omitempty"`
}

func (x *DhtLookupResponse_AggregatedProviderRecords) Reset() {
	*x = DhtLookupResponse_AggregatedProviderRecords{}
	if protoimpl.UnsafeEnabled {
		mi := &file_message_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DhtLookupResponse_AggregatedProviderRecords) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DhtLookupResponse_AggregatedProviderRecords) ProtoMessage() {}

func (x *DhtLookupResponse_AggregatedProviderRecords) ProtoReflect() protoreflect.Message {
	mi := &file_message_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DhtLookupResponse_AggregatedProviderRecords.ProtoReflect.Descriptor instead.
func (*DhtLookupResponse_AggregatedProviderRecords) Descriptor() ([]byte, []int) {
	return file_message_proto_rawDescGZIP(), []int{4, 2}
}

func (x *DhtLookupResponse_AggregatedProviderRecords) GetId() []byte {
	if x != nil {
		return x.Id
	}
	return nil
}

func (x *DhtLookupResponse_AggregatedProviderRecords) GetProviderRecords() []*DhtLookupResponse_ProviderRecord {
	if x != nil {
		return x.ProviderRecords
	}
	return nil
}

var File_message_proto protoreflect.FileDescriptor

var file_message_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22,
	0x71, 0x0a, 0x09, 0x45, 0x6e, 0x63, 0x50, 0x65, 0x65, 0x72, 0x49, 0x64, 0x12, 0x34, 0x0a, 0x15,
	0x45, 0x6e, 0x63, 0x50, 0x65, 0x65, 0x72, 0x49, 0x64, 0x46, 0x6f, 0x72, 0x6d, 0x61, 0x74, 0x56,
	0x61, 0x72, 0x69, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x15, 0x45, 0x6e, 0x63,
	0x50, 0x65, 0x65, 0x72, 0x49, 0x64, 0x46, 0x6f, 0x72, 0x6d, 0x61, 0x74, 0x56, 0x61, 0x72, 0x69,
	0x6e, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x4e, 0x6f, 0x6e, 0x63, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x0c, 0x52, 0x05, 0x4e, 0x6f, 0x6e, 0x63, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x50, 0x61, 0x79, 0x6c,
	0x6f, 0x61, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x07, 0x50, 0x61, 0x79, 0x6c, 0x6f,
	0x61, 0x64, 0x22, 0x89, 0x01, 0x0a, 0x11, 0x44, 0x68, 0x74, 0x50, 0x72, 0x6f, 0x76, 0x69, 0x64,
	0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x49, 0x44, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0c, 0x52, 0x02, 0x49, 0x44, 0x12, 0x1c, 0x0a, 0x09, 0x53, 0x65, 0x72, 0x76,
	0x65, 0x72, 0x4b, 0x65, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x09, 0x53, 0x65, 0x72,
	0x76, 0x65, 0x72, 0x4b, 0x65, 0x79, 0x12, 0x28, 0x0a, 0x09, 0x45, 0x6e, 0x63, 0x50, 0x65, 0x65,
	0x72, 0x49, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0a, 0x2e, 0x45, 0x6e, 0x63, 0x50,
	0x65, 0x65, 0x72, 0x49, 0x64, 0x52, 0x09, 0x45, 0x6e, 0x63, 0x50, 0x65, 0x65, 0x72, 0x49, 0x64,
	0x12, 0x1c, 0x0a, 0x09, 0x53, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x0c, 0x52, 0x09, 0x53, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x22, 0x65,
	0x0a, 0x12, 0x44, 0x68, 0x74, 0x50, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x32, 0x0a, 0x06, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0e, 0x32, 0x1a, 0x2e, 0x44, 0x68, 0x74, 0x50, 0x72, 0x6f, 0x76, 0x69, 0x64,
	0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2e, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73,
	0x52, 0x06, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x22, 0x1b, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74,
	0x75, 0x73, 0x12, 0x06, 0x0a, 0x02, 0x4f, 0x4b, 0x10, 0x00, 0x12, 0x09, 0x0a, 0x05, 0x45, 0x52,
	0x52, 0x4f, 0x52, 0x10, 0x01, 0x22, 0x2a, 0x0a, 0x10, 0x44, 0x68, 0x74, 0x4c, 0x6f, 0x6f, 0x6b,
	0x75, 0x70, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x50, 0x72, 0x65,
	0x66, 0x69, 0x78, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x06, 0x50, 0x72, 0x65, 0x66, 0x69,
	0x78, 0x22, 0xeb, 0x05, 0x0a, 0x11, 0x44, 0x68, 0x74, 0x4c, 0x6f, 0x6f, 0x6b, 0x75, 0x70, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x44, 0x0a, 0x1d, 0x44, 0x68, 0x74, 0x4c, 0x6f,
	0x6f, 0x6b, 0x75, 0x70, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x46, 0x6f, 0x72, 0x6d,
	0x61, 0x74, 0x56, 0x61, 0x72, 0x69, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x1d,
	0x44, 0x68, 0x74, 0x4c, 0x6f, 0x6f, 0x6b, 0x75, 0x70, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x46, 0x6f, 0x72, 0x6d, 0x61, 0x74, 0x56, 0x61, 0x72, 0x69, 0x6e, 0x74, 0x12, 0x12, 0x0a,
	0x04, 0x46, 0x6c, 0x61, 0x67, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x46, 0x6c, 0x61,
	0x67, 0x12, 0x2d, 0x0a, 0x05, 0x50, 0x65, 0x65, 0x72, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x17, 0x2e, 0x44, 0x68, 0x74, 0x4c, 0x6f, 0x6f, 0x6b, 0x75, 0x70, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x2e, 0x50, 0x65, 0x65, 0x72, 0x52, 0x05, 0x50, 0x65, 0x65, 0x72, 0x73,
	0x12, 0x56, 0x0a, 0x0f, 0x50, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x52, 0x65, 0x63, 0x6f,
	0x72, 0x64, 0x73, 0x18, 0x04, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x2c, 0x2e, 0x44, 0x68, 0x74, 0x4c,
	0x6f, 0x6f, 0x6b, 0x75, 0x70, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2e, 0x41, 0x67,
	0x67, 0x72, 0x65, 0x67, 0x61, 0x74, 0x65, 0x64, 0x50, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72,
	0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x73, 0x52, 0x0f, 0x50, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65,
	0x72, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x73, 0x1a, 0xd5, 0x01, 0x0a, 0x04, 0x50, 0x65, 0x65,
	0x72, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x02, 0x69,
	0x64, 0x12, 0x14, 0x0a, 0x05, 0x61, 0x64, 0x64, 0x72, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0c,
	0x52, 0x05, 0x61, 0x64, 0x64, 0x72, 0x73, 0x12, 0x4e, 0x0a, 0x0e, 0x63, 0x6f, 0x6e, 0x6e, 0x65,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x54, 0x79, 0x70, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0e, 0x32,
	0x26, 0x2e, 0x44, 0x68, 0x74, 0x4c, 0x6f, 0x6f, 0x6b, 0x75, 0x70, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x2e, 0x50, 0x65, 0x65, 0x72, 0x2e, 0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74,
	0x69, 0x6f, 0x6e, 0x54, 0x79, 0x70, 0x65, 0x52, 0x0e, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74,
	0x69, 0x6f, 0x6e, 0x54, 0x79, 0x70, 0x65, 0x22, 0x57, 0x0a, 0x0e, 0x43, 0x6f, 0x6e, 0x6e, 0x65,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x54, 0x79, 0x70, 0x65, 0x12, 0x11, 0x0a, 0x0d, 0x4e, 0x4f, 0x54,
	0x5f, 0x43, 0x4f, 0x4e, 0x4e, 0x45, 0x43, 0x54, 0x45, 0x44, 0x10, 0x00, 0x12, 0x0d, 0x0a, 0x09,
	0x43, 0x4f, 0x4e, 0x4e, 0x45, 0x43, 0x54, 0x45, 0x44, 0x10, 0x01, 0x12, 0x0f, 0x0a, 0x0b, 0x43,
	0x41, 0x4e, 0x5f, 0x43, 0x4f, 0x4e, 0x4e, 0x45, 0x43, 0x54, 0x10, 0x02, 0x12, 0x12, 0x0a, 0x0e,
	0x43, 0x41, 0x4e, 0x4e, 0x4f, 0x54, 0x5f, 0x43, 0x4f, 0x4e, 0x4e, 0x45, 0x43, 0x54, 0x10, 0x03,
	0x1a, 0xa2, 0x01, 0x0a, 0x0e, 0x50, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x52, 0x65, 0x63,
	0x6f, 0x72, 0x64, 0x12, 0x22, 0x0a, 0x0c, 0x46, 0x6f, 0x72, 0x6d, 0x61, 0x74, 0x56, 0x61, 0x72,
	0x69, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0c, 0x46, 0x6f, 0x72, 0x6d, 0x61,
	0x74, 0x56, 0x61, 0x72, 0x69, 0x6e, 0x74, 0x12, 0x28, 0x0a, 0x09, 0x45, 0x6e, 0x63, 0x50, 0x65,
	0x65, 0x72, 0x49, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0a, 0x2e, 0x45, 0x6e, 0x63,
	0x50, 0x65, 0x65, 0x72, 0x49, 0x64, 0x52, 0x09, 0x45, 0x6e, 0x63, 0x50, 0x65, 0x65, 0x72, 0x49,
	0x64, 0x12, 0x20, 0x0a, 0x0b, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x4e, 0x6f, 0x6e, 0x63, 0x65,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0b, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x4e, 0x6f,
	0x6e, 0x63, 0x65, 0x12, 0x20, 0x0a, 0x0b, 0x45, 0x6e, 0x63, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61,
	0x74, 0x61, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0b, 0x45, 0x6e, 0x63, 0x4d, 0x65, 0x74,
	0x61, 0x64, 0x61, 0x74, 0x61, 0x1a, 0x78, 0x0a, 0x19, 0x41, 0x67, 0x67, 0x72, 0x65, 0x67, 0x61,
	0x74, 0x65, 0x64, 0x50, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x52, 0x65, 0x63, 0x6f, 0x72,
	0x64, 0x73, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x02,
	0x69, 0x64, 0x12, 0x4b, 0x0a, 0x0f, 0x50, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x52, 0x65,
	0x63, 0x6f, 0x72, 0x64, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x21, 0x2e, 0x44, 0x68,
	0x74, 0x4c, 0x6f, 0x6f, 0x6b, 0x75, 0x70, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2e,
	0x50, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x52, 0x0f,
	0x50, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x73, 0x22,
	0xb3, 0x02, 0x0a, 0x0a, 0x44, 0x68, 0x74, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x44,
	0x0a, 0x12, 0x50, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x54, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x44, 0x68, 0x74,
	0x50, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x48, 0x00,
	0x52, 0x12, 0x50, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x54, 0x79, 0x70, 0x65, 0x12, 0x47, 0x0a, 0x13, 0x50, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x54, 0x79, 0x70, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x13, 0x2e, 0x44, 0x68, 0x74, 0x50, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x48, 0x00, 0x52, 0x13, 0x50, 0x72, 0x6f, 0x76, 0x69, 0x64,
	0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x54, 0x79, 0x70, 0x65, 0x12, 0x41, 0x0a,
	0x11, 0x4c, 0x6f, 0x6f, 0x6b, 0x75, 0x70, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x54, 0x79,
	0x70, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x44, 0x68, 0x74, 0x4c, 0x6f,
	0x6f, 0x6b, 0x75, 0x70, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x48, 0x00, 0x52, 0x11, 0x4c,
	0x6f, 0x6f, 0x6b, 0x75, 0x70, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x54, 0x79, 0x70, 0x65,
	0x12, 0x44, 0x0a, 0x12, 0x4c, 0x6f, 0x6f, 0x6b, 0x75, 0x70, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x54, 0x79, 0x70, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x44,
	0x68, 0x74, 0x4c, 0x6f, 0x6f, 0x6b, 0x75, 0x70, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x48, 0x00, 0x52, 0x12, 0x4c, 0x6f, 0x6f, 0x6b, 0x75, 0x70, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x54, 0x79, 0x70, 0x65, 0x42, 0x0d, 0x0a, 0x0b, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x54, 0x79, 0x70, 0x65, 0x42, 0x06, 0x5a, 0x04, 0x2e, 0x2f, 0x70, 0x62, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_message_proto_rawDescOnce sync.Once
	file_message_proto_rawDescData = file_message_proto_rawDesc
)

func file_message_proto_rawDescGZIP() []byte {
	file_message_proto_rawDescOnce.Do(func() {
		file_message_proto_rawDescData = protoimpl.X.CompressGZIP(file_message_proto_rawDescData)
	})
	return file_message_proto_rawDescData
}

var file_message_proto_enumTypes = make([]protoimpl.EnumInfo, 2)
var file_message_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_message_proto_goTypes = []interface{}{
	(DhtProvideResponseStatus)(0),                       // 0: DhtProvideResponse.status
	(DhtLookupResponse_Peer_ConnectionType)(0),          // 1: DhtLookupResponse.Peer.ConnectionType
	(*EncPeerId)(nil),                                   // 2: EncPeerId
	(*DhtProvideRequest)(nil),                           // 3: DhtProvideRequest
	(*DhtProvideResponse)(nil),                          // 4: DhtProvideResponse
	(*DhtLookupRequest)(nil),                            // 5: DhtLookupRequest
	(*DhtLookupResponse)(nil),                           // 6: DhtLookupResponse
	(*DhtMessage)(nil),                                  // 7: DhtMessage
	(*DhtLookupResponse_Peer)(nil),                      // 8: DhtLookupResponse.Peer
	(*DhtLookupResponse_ProviderRecord)(nil),            // 9: DhtLookupResponse.ProviderRecord
	(*DhtLookupResponse_AggregatedProviderRecords)(nil), // 10: DhtLookupResponse.AggregatedProviderRecords
}
var file_message_proto_depIdxs = []int32{
	2,  // 0: DhtProvideRequest.EncPeerId:type_name -> EncPeerId
	0,  // 1: DhtProvideResponse.Status:type_name -> DhtProvideResponse.status
	8,  // 2: DhtLookupResponse.Peers:type_name -> DhtLookupResponse.Peer
	10, // 3: DhtLookupResponse.ProviderRecords:type_name -> DhtLookupResponse.AggregatedProviderRecords
	3,  // 4: DhtMessage.ProvideRequestType:type_name -> DhtProvideRequest
	4,  // 5: DhtMessage.ProvideResponseType:type_name -> DhtProvideResponse
	5,  // 6: DhtMessage.LookupRequestType:type_name -> DhtLookupRequest
	6,  // 7: DhtMessage.LookupResponseType:type_name -> DhtLookupResponse
	1,  // 8: DhtLookupResponse.Peer.connectionType:type_name -> DhtLookupResponse.Peer.ConnectionType
	2,  // 9: DhtLookupResponse.ProviderRecord.EncPeerId:type_name -> EncPeerId
	9,  // 10: DhtLookupResponse.AggregatedProviderRecords.ProviderRecords:type_name -> DhtLookupResponse.ProviderRecord
	11, // [11:11] is the sub-list for method output_type
	11, // [11:11] is the sub-list for method input_type
	11, // [11:11] is the sub-list for extension type_name
	11, // [11:11] is the sub-list for extension extendee
	0,  // [0:11] is the sub-list for field type_name
}

func init() { file_message_proto_init() }
func file_message_proto_init() {
	if File_message_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_message_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EncPeerId); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_message_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DhtProvideRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_message_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DhtProvideResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_message_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DhtLookupRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_message_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DhtLookupResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_message_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DhtMessage); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_message_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DhtLookupResponse_Peer); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_message_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DhtLookupResponse_ProviderRecord); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_message_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DhtLookupResponse_AggregatedProviderRecords); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	file_message_proto_msgTypes[5].OneofWrappers = []interface{}{
		(*DhtMessage_ProvideRequestType)(nil),
		(*DhtMessage_ProvideResponseType)(nil),
		(*DhtMessage_LookupRequestType)(nil),
		(*DhtMessage_LookupResponseType)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_message_proto_rawDesc,
			NumEnums:      2,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_message_proto_goTypes,
		DependencyIndexes: file_message_proto_depIdxs,
		EnumInfos:         file_message_proto_enumTypes,
		MessageInfos:      file_message_proto_msgTypes,
	}.Build()
	File_message_proto = out.File
	file_message_proto_rawDesc = nil
	file_message_proto_goTypes = nil
	file_message_proto_depIdxs = nil
}
