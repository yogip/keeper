// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v3.6.1
// source: internal/proto/keeper.proto

package proto

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

type SecretType int32

const (
	SecretType_PASSWORD SecretType = 0
	SecretType_CARD     SecretType = 1
	SecretType_NOTE     SecretType = 2
	SecretType_FILE     SecretType = 3
)

// Enum value maps for SecretType.
var (
	SecretType_name = map[int32]string{
		0: "PASSWORD",
		1: "CARD",
		2: "NOTE",
		3: "FILE",
	}
	SecretType_value = map[string]int32{
		"PASSWORD": 0,
		"CARD":     1,
		"NOTE":     2,
		"FILE":     3,
	}
)

func (x SecretType) Enum() *SecretType {
	p := new(SecretType)
	*p = x
	return p
}

func (x SecretType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (SecretType) Descriptor() protoreflect.EnumDescriptor {
	return file_internal_proto_keeper_proto_enumTypes[0].Descriptor()
}

func (SecretType) Type() protoreflect.EnumType {
	return &file_internal_proto_keeper_proto_enumTypes[0]
}

func (x SecretType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use SecretType.Descriptor instead.
func (SecretType) EnumDescriptor() ([]byte, []int) {
	return file_internal_proto_keeper_proto_rawDescGZIP(), []int{0}
}

// IAM section
type LoginRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Login    string `protobuf:"bytes,1,opt,name=login,proto3" json:"login,omitempty"`
	Password string `protobuf:"bytes,2,opt,name=password,proto3" json:"password,omitempty"`
}

func (x *LoginRequest) Reset() {
	*x = LoginRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_proto_keeper_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LoginRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LoginRequest) ProtoMessage() {}

func (x *LoginRequest) ProtoReflect() protoreflect.Message {
	mi := &file_internal_proto_keeper_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LoginRequest.ProtoReflect.Descriptor instead.
func (*LoginRequest) Descriptor() ([]byte, []int) {
	return file_internal_proto_keeper_proto_rawDescGZIP(), []int{0}
}

func (x *LoginRequest) GetLogin() string {
	if x != nil {
		return x.Login
	}
	return ""
}

func (x *LoginRequest) GetPassword() string {
	if x != nil {
		return x.Password
	}
	return ""
}

type SignUpRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Login    string `protobuf:"bytes,1,opt,name=login,proto3" json:"login,omitempty"`
	Password string `protobuf:"bytes,2,opt,name=password,proto3" json:"password,omitempty"`
}

func (x *SignUpRequest) Reset() {
	*x = SignUpRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_proto_keeper_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SignUpRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SignUpRequest) ProtoMessage() {}

func (x *SignUpRequest) ProtoReflect() protoreflect.Message {
	mi := &file_internal_proto_keeper_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SignUpRequest.ProtoReflect.Descriptor instead.
func (*SignUpRequest) Descriptor() ([]byte, []int) {
	return file_internal_proto_keeper_proto_rawDescGZIP(), []int{1}
}

func (x *SignUpRequest) GetLogin() string {
	if x != nil {
		return x.Login
	}
	return ""
}

func (x *SignUpRequest) GetPassword() string {
	if x != nil {
		return x.Password
	}
	return ""
}

type Token struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Token string `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
}

func (x *Token) Reset() {
	*x = Token{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_proto_keeper_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Token) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Token) ProtoMessage() {}

func (x *Token) ProtoReflect() protoreflect.Message {
	mi := &file_internal_proto_keeper_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Token.ProtoReflect.Descriptor instead.
func (*Token) Descriptor() ([]byte, []int) {
	return file_internal_proto_keeper_proto_rawDescGZIP(), []int{2}
}

func (x *Token) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

// List Secrets
type ListRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
}

func (x *ListRequest) Reset() {
	*x = ListRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_proto_keeper_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListRequest) ProtoMessage() {}

func (x *ListRequest) ProtoReflect() protoreflect.Message {
	mi := &file_internal_proto_keeper_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListRequest.ProtoReflect.Descriptor instead.
func (*ListRequest) Descriptor() ([]byte, []int) {
	return file_internal_proto_keeper_proto_rawDescGZIP(), []int{3}
}

func (x *ListRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

type SecretMeta struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id   int64      `protobuf:"zigzag64,1,opt,name=id,proto3" json:"id,omitempty"`
	Name string     `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Note string     `protobuf:"bytes,3,opt,name=note,proto3" json:"note,omitempty"`
	Type SecretType `protobuf:"varint,4,opt,name=type,proto3,enum=proto.SecretType" json:"type,omitempty"`
}

func (x *SecretMeta) Reset() {
	*x = SecretMeta{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_proto_keeper_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SecretMeta) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SecretMeta) ProtoMessage() {}

func (x *SecretMeta) ProtoReflect() protoreflect.Message {
	mi := &file_internal_proto_keeper_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SecretMeta.ProtoReflect.Descriptor instead.
func (*SecretMeta) Descriptor() ([]byte, []int) {
	return file_internal_proto_keeper_proto_rawDescGZIP(), []int{4}
}

func (x *SecretMeta) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *SecretMeta) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *SecretMeta) GetNote() string {
	if x != nil {
		return x.Note
	}
	return ""
}

func (x *SecretMeta) GetType() SecretType {
	if x != nil {
		return x.Type
	}
	return SecretType_PASSWORD
}

type ListResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Secrets []*SecretMeta `protobuf:"bytes,1,rep,name=secrets,proto3" json:"secrets,omitempty"`
}

func (x *ListResponse) Reset() {
	*x = ListResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_proto_keeper_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListResponse) ProtoMessage() {}

func (x *ListResponse) ProtoReflect() protoreflect.Message {
	mi := &file_internal_proto_keeper_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListResponse.ProtoReflect.Descriptor instead.
func (*ListResponse) Descriptor() ([]byte, []int) {
	return file_internal_proto_keeper_proto_rawDescGZIP(), []int{5}
}

func (x *ListResponse) GetSecrets() []*SecretMeta {
	if x != nil {
		return x.Secrets
	}
	return nil
}

// Secret section
type Secret struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id      int64      `protobuf:"zigzag64,1,opt,name=id,proto3" json:"id,omitempty"`
	Name    string     `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Note    string     `protobuf:"bytes,3,opt,name=note,proto3" json:"note,omitempty"`
	Type    SecretType `protobuf:"varint,4,opt,name=type,proto3,enum=proto.SecretType" json:"type,omitempty"`
	Payload []byte     `protobuf:"bytes,5,opt,name=payload,proto3" json:"payload,omitempty"`
}

func (x *Secret) Reset() {
	*x = Secret{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_proto_keeper_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Secret) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Secret) ProtoMessage() {}

func (x *Secret) ProtoReflect() protoreflect.Message {
	mi := &file_internal_proto_keeper_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Secret.ProtoReflect.Descriptor instead.
func (*Secret) Descriptor() ([]byte, []int) {
	return file_internal_proto_keeper_proto_rawDescGZIP(), []int{6}
}

func (x *Secret) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Secret) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Secret) GetNote() string {
	if x != nil {
		return x.Note
	}
	return ""
}

func (x *Secret) GetType() SecretType {
	if x != nil {
		return x.Type
	}
	return SecretType_PASSWORD
}

func (x *Secret) GetPayload() []byte {
	if x != nil {
		return x.Payload
	}
	return nil
}

type SecretRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id   int64      `protobuf:"zigzag64,1,opt,name=id,proto3" json:"id,omitempty"`
	Type SecretType `protobuf:"varint,2,opt,name=type,proto3,enum=proto.SecretType" json:"type,omitempty"`
}

func (x *SecretRequest) Reset() {
	*x = SecretRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_proto_keeper_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SecretRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SecretRequest) ProtoMessage() {}

func (x *SecretRequest) ProtoReflect() protoreflect.Message {
	mi := &file_internal_proto_keeper_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SecretRequest.ProtoReflect.Descriptor instead.
func (*SecretRequest) Descriptor() ([]byte, []int) {
	return file_internal_proto_keeper_proto_rawDescGZIP(), []int{7}
}

func (x *SecretRequest) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *SecretRequest) GetType() SecretType {
	if x != nil {
		return x.Type
	}
	return SecretType_PASSWORD
}

type SecretCreateRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name    string     `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Type    SecretType `protobuf:"varint,2,opt,name=type,proto3,enum=proto.SecretType" json:"type,omitempty"`
	Payload []byte     `protobuf:"bytes,3,opt,name=payload,proto3" json:"payload,omitempty"`
	Note    string     `protobuf:"bytes,4,opt,name=note,proto3" json:"note,omitempty"`
}

func (x *SecretCreateRequest) Reset() {
	*x = SecretCreateRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_proto_keeper_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SecretCreateRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SecretCreateRequest) ProtoMessage() {}

func (x *SecretCreateRequest) ProtoReflect() protoreflect.Message {
	mi := &file_internal_proto_keeper_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SecretCreateRequest.ProtoReflect.Descriptor instead.
func (*SecretCreateRequest) Descriptor() ([]byte, []int) {
	return file_internal_proto_keeper_proto_rawDescGZIP(), []int{8}
}

func (x *SecretCreateRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *SecretCreateRequest) GetType() SecretType {
	if x != nil {
		return x.Type
	}
	return SecretType_PASSWORD
}

func (x *SecretCreateRequest) GetPayload() []byte {
	if x != nil {
		return x.Payload
	}
	return nil
}

func (x *SecretCreateRequest) GetNote() string {
	if x != nil {
		return x.Note
	}
	return ""
}

type SecretUpdateRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id      int64      `protobuf:"zigzag64,1,opt,name=id,proto3" json:"id,omitempty"`
	Name    string     `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Note    string     `protobuf:"bytes,3,opt,name=note,proto3" json:"note,omitempty"`
	Type    SecretType `protobuf:"varint,4,opt,name=type,proto3,enum=proto.SecretType" json:"type,omitempty"`
	Payload []byte     `protobuf:"bytes,5,opt,name=payload,proto3" json:"payload,omitempty"`
}

func (x *SecretUpdateRequest) Reset() {
	*x = SecretUpdateRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_proto_keeper_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SecretUpdateRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SecretUpdateRequest) ProtoMessage() {}

func (x *SecretUpdateRequest) ProtoReflect() protoreflect.Message {
	mi := &file_internal_proto_keeper_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SecretUpdateRequest.ProtoReflect.Descriptor instead.
func (*SecretUpdateRequest) Descriptor() ([]byte, []int) {
	return file_internal_proto_keeper_proto_rawDescGZIP(), []int{9}
}

func (x *SecretUpdateRequest) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *SecretUpdateRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *SecretUpdateRequest) GetNote() string {
	if x != nil {
		return x.Note
	}
	return ""
}

func (x *SecretUpdateRequest) GetType() SecretType {
	if x != nil {
		return x.Type
	}
	return SecretType_PASSWORD
}

func (x *SecretUpdateRequest) GetPayload() []byte {
	if x != nil {
		return x.Payload
	}
	return nil
}

type SecretFileCreateRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name     string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	FileName string `protobuf:"bytes,2,opt,name=file_name,json=fileName,proto3" json:"file_name,omitempty"`
	Payload  []byte `protobuf:"bytes,3,opt,name=payload,proto3" json:"payload,omitempty"`
	Note     string `protobuf:"bytes,4,opt,name=note,proto3" json:"note,omitempty"`
}

func (x *SecretFileCreateRequest) Reset() {
	*x = SecretFileCreateRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_proto_keeper_proto_msgTypes[10]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SecretFileCreateRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SecretFileCreateRequest) ProtoMessage() {}

func (x *SecretFileCreateRequest) ProtoReflect() protoreflect.Message {
	mi := &file_internal_proto_keeper_proto_msgTypes[10]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SecretFileCreateRequest.ProtoReflect.Descriptor instead.
func (*SecretFileCreateRequest) Descriptor() ([]byte, []int) {
	return file_internal_proto_keeper_proto_rawDescGZIP(), []int{10}
}

func (x *SecretFileCreateRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *SecretFileCreateRequest) GetFileName() string {
	if x != nil {
		return x.FileName
	}
	return ""
}

func (x *SecretFileCreateRequest) GetPayload() []byte {
	if x != nil {
		return x.Payload
	}
	return nil
}

func (x *SecretFileCreateRequest) GetNote() string {
	if x != nil {
		return x.Note
	}
	return ""
}

type SecretFileCreateResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id       int64  `protobuf:"zigzag64,1,opt,name=id,proto3" json:"id,omitempty"`
	Name     string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	FileName string `protobuf:"bytes,3,opt,name=file_name,json=fileName,proto3" json:"file_name,omitempty"`
	Note     string `protobuf:"bytes,4,opt,name=note,proto3" json:"note,omitempty"`
}

func (x *SecretFileCreateResponse) Reset() {
	*x = SecretFileCreateResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_proto_keeper_proto_msgTypes[11]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SecretFileCreateResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SecretFileCreateResponse) ProtoMessage() {}

func (x *SecretFileCreateResponse) ProtoReflect() protoreflect.Message {
	mi := &file_internal_proto_keeper_proto_msgTypes[11]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SecretFileCreateResponse.ProtoReflect.Descriptor instead.
func (*SecretFileCreateResponse) Descriptor() ([]byte, []int) {
	return file_internal_proto_keeper_proto_rawDescGZIP(), []int{11}
}

func (x *SecretFileCreateResponse) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *SecretFileCreateResponse) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *SecretFileCreateResponse) GetFileName() string {
	if x != nil {
		return x.FileName
	}
	return ""
}

func (x *SecretFileCreateResponse) GetNote() string {
	if x != nil {
		return x.Note
	}
	return ""
}

var File_internal_proto_keeper_proto protoreflect.FileDescriptor

var file_internal_proto_keeper_proto_rawDesc = []byte{
	0x0a, 0x1b, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2f, 0x6b, 0x65, 0x65, 0x70, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x22, 0x40, 0x0a, 0x0c, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x6c, 0x6f, 0x67, 0x69, 0x6e, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x05, 0x6c, 0x6f, 0x67, 0x69, 0x6e, 0x12, 0x1a, 0x0a, 0x08, 0x70, 0x61,
	0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x70, 0x61,
	0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x22, 0x41, 0x0a, 0x0d, 0x53, 0x69, 0x67, 0x6e, 0x55, 0x70,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x6c, 0x6f, 0x67, 0x69, 0x6e,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x6c, 0x6f, 0x67, 0x69, 0x6e, 0x12, 0x1a, 0x0a,
	0x08, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x08, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x22, 0x1d, 0x0a, 0x05, 0x54, 0x6f, 0x6b,
	0x65, 0x6e, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x22, 0x21, 0x0a, 0x0b, 0x4c, 0x69, 0x73, 0x74,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x22, 0x6b, 0x0a, 0x0a, 0x53,
	0x65, 0x63, 0x72, 0x65, 0x74, 0x4d, 0x65, 0x74, 0x61, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x12, 0x52, 0x02, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x12, 0x0a,
	0x04, 0x6e, 0x6f, 0x74, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x6f, 0x74,
	0x65, 0x12, 0x25, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0e, 0x32,
	0x11, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x53, 0x65, 0x63, 0x72, 0x65, 0x74, 0x54, 0x79,
	0x70, 0x65, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x22, 0x3b, 0x0a, 0x0c, 0x4c, 0x69, 0x73, 0x74,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2b, 0x0a, 0x07, 0x73, 0x65, 0x63, 0x72,
	0x65, 0x74, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2e, 0x53, 0x65, 0x63, 0x72, 0x65, 0x74, 0x4d, 0x65, 0x74, 0x61, 0x52, 0x07, 0x73, 0x65,
	0x63, 0x72, 0x65, 0x74, 0x73, 0x22, 0x81, 0x01, 0x0a, 0x06, 0x53, 0x65, 0x63, 0x72, 0x65, 0x74,
	0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x12, 0x52, 0x02, 0x69, 0x64,
	0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04,
	0x6e, 0x61, 0x6d, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x6f, 0x74, 0x65, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x04, 0x6e, 0x6f, 0x74, 0x65, 0x12, 0x25, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x11, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x53,
	0x65, 0x63, 0x72, 0x65, 0x74, 0x54, 0x79, 0x70, 0x65, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12,
	0x18, 0x0a, 0x07, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0c,
	0x52, 0x07, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x22, 0x46, 0x0a, 0x0d, 0x53, 0x65, 0x63,
	0x72, 0x65, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x12, 0x52, 0x02, 0x69, 0x64, 0x12, 0x25, 0x0a, 0x04, 0x74, 0x79,
	0x70, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x11, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2e, 0x53, 0x65, 0x63, 0x72, 0x65, 0x74, 0x54, 0x79, 0x70, 0x65, 0x52, 0x04, 0x74, 0x79, 0x70,
	0x65, 0x22, 0x7e, 0x0a, 0x13, 0x53, 0x65, 0x63, 0x72, 0x65, 0x74, 0x43, 0x72, 0x65, 0x61, 0x74,
	0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x25, 0x0a, 0x04,
	0x74, 0x79, 0x70, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x11, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2e, 0x53, 0x65, 0x63, 0x72, 0x65, 0x74, 0x54, 0x79, 0x70, 0x65, 0x52, 0x04, 0x74,
	0x79, 0x70, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x0c, 0x52, 0x07, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x12, 0x12, 0x0a,
	0x04, 0x6e, 0x6f, 0x74, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x6f, 0x74,
	0x65, 0x22, 0x8e, 0x01, 0x0a, 0x13, 0x53, 0x65, 0x63, 0x72, 0x65, 0x74, 0x55, 0x70, 0x64, 0x61,
	0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x12, 0x52, 0x02, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x12, 0x0a,
	0x04, 0x6e, 0x6f, 0x74, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x6f, 0x74,
	0x65, 0x12, 0x25, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0e, 0x32,
	0x11, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x53, 0x65, 0x63, 0x72, 0x65, 0x74, 0x54, 0x79,
	0x70, 0x65, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x70, 0x61, 0x79, 0x6c,
	0x6f, 0x61, 0x64, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x07, 0x70, 0x61, 0x79, 0x6c, 0x6f,
	0x61, 0x64, 0x22, 0x78, 0x0a, 0x17, 0x53, 0x65, 0x63, 0x72, 0x65, 0x74, 0x46, 0x69, 0x6c, 0x65,
	0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a,
	0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d,
	0x65, 0x12, 0x1b, 0x0a, 0x09, 0x66, 0x69, 0x6c, 0x65, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x66, 0x69, 0x6c, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x18,
	0x0a, 0x07, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52,
	0x07, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x6f, 0x74, 0x65,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x6f, 0x74, 0x65, 0x22, 0x6f, 0x0a, 0x18,
	0x53, 0x65, 0x63, 0x72, 0x65, 0x74, 0x46, 0x69, 0x6c, 0x65, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x12, 0x52, 0x02, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x1b, 0x0a, 0x09,
	0x66, 0x69, 0x6c, 0x65, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x08, 0x66, 0x69, 0x6c, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x6f, 0x74,
	0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x6f, 0x74, 0x65, 0x2a, 0x38, 0x0a,
	0x0a, 0x53, 0x65, 0x63, 0x72, 0x65, 0x74, 0x54, 0x79, 0x70, 0x65, 0x12, 0x0c, 0x0a, 0x08, 0x50,
	0x41, 0x53, 0x53, 0x57, 0x4f, 0x52, 0x44, 0x10, 0x00, 0x12, 0x08, 0x0a, 0x04, 0x43, 0x41, 0x52,
	0x44, 0x10, 0x01, 0x12, 0x08, 0x0a, 0x04, 0x4e, 0x4f, 0x54, 0x45, 0x10, 0x02, 0x12, 0x08, 0x0a,
	0x04, 0x46, 0x49, 0x4c, 0x45, 0x10, 0x03, 0x32, 0x91, 0x03, 0x0a, 0x06, 0x4b, 0x65, 0x65, 0x70,
	0x65, 0x72, 0x12, 0x2a, 0x0a, 0x05, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x12, 0x13, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2e, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x0c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x2c,
	0x0a, 0x06, 0x53, 0x69, 0x67, 0x6e, 0x55, 0x70, 0x12, 0x14, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2e, 0x53, 0x69, 0x67, 0x6e, 0x55, 0x70, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x0c,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x36, 0x0a, 0x0b,
	0x4c, 0x69, 0x73, 0x74, 0x53, 0x65, 0x63, 0x72, 0x65, 0x74, 0x73, 0x12, 0x12, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x13, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x30, 0x0a, 0x09, 0x47, 0x65, 0x74, 0x53, 0x65, 0x63, 0x72, 0x65,
	0x74, 0x12, 0x14, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x53, 0x65, 0x63, 0x72, 0x65, 0x74,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x0d, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e,
	0x53, 0x65, 0x63, 0x72, 0x65, 0x74, 0x12, 0x39, 0x0a, 0x0c, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65,
	0x53, 0x65, 0x63, 0x72, 0x65, 0x74, 0x12, 0x1a, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x53,
	0x65, 0x63, 0x72, 0x65, 0x74, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x0d, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x53, 0x65, 0x63, 0x72, 0x65,
	0x74, 0x12, 0x39, 0x0a, 0x0c, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x53, 0x65, 0x63, 0x72, 0x65,
	0x74, 0x12, 0x1a, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x53, 0x65, 0x63, 0x72, 0x65, 0x74,
	0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x0d, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x53, 0x65, 0x63, 0x72, 0x65, 0x74, 0x12, 0x4d, 0x0a, 0x0a,
	0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x46, 0x69, 0x6c, 0x65, 0x12, 0x1e, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2e, 0x53, 0x65, 0x63, 0x72, 0x65, 0x74, 0x46, 0x69, 0x6c, 0x65, 0x43, 0x72, 0x65,
	0x61, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1f, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2e, 0x53, 0x65, 0x63, 0x72, 0x65, 0x74, 0x46, 0x69, 0x6c, 0x65, 0x43, 0x72, 0x65,
	0x61, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x0f, 0x5a, 0x0d, 0x2e,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_internal_proto_keeper_proto_rawDescOnce sync.Once
	file_internal_proto_keeper_proto_rawDescData = file_internal_proto_keeper_proto_rawDesc
)

func file_internal_proto_keeper_proto_rawDescGZIP() []byte {
	file_internal_proto_keeper_proto_rawDescOnce.Do(func() {
		file_internal_proto_keeper_proto_rawDescData = protoimpl.X.CompressGZIP(file_internal_proto_keeper_proto_rawDescData)
	})
	return file_internal_proto_keeper_proto_rawDescData
}

var file_internal_proto_keeper_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_internal_proto_keeper_proto_msgTypes = make([]protoimpl.MessageInfo, 12)
var file_internal_proto_keeper_proto_goTypes = []any{
	(SecretType)(0),                  // 0: proto.SecretType
	(*LoginRequest)(nil),             // 1: proto.LoginRequest
	(*SignUpRequest)(nil),            // 2: proto.SignUpRequest
	(*Token)(nil),                    // 3: proto.Token
	(*ListRequest)(nil),              // 4: proto.ListRequest
	(*SecretMeta)(nil),               // 5: proto.SecretMeta
	(*ListResponse)(nil),             // 6: proto.ListResponse
	(*Secret)(nil),                   // 7: proto.Secret
	(*SecretRequest)(nil),            // 8: proto.SecretRequest
	(*SecretCreateRequest)(nil),      // 9: proto.SecretCreateRequest
	(*SecretUpdateRequest)(nil),      // 10: proto.SecretUpdateRequest
	(*SecretFileCreateRequest)(nil),  // 11: proto.SecretFileCreateRequest
	(*SecretFileCreateResponse)(nil), // 12: proto.SecretFileCreateResponse
}
var file_internal_proto_keeper_proto_depIdxs = []int32{
	0,  // 0: proto.SecretMeta.type:type_name -> proto.SecretType
	5,  // 1: proto.ListResponse.secrets:type_name -> proto.SecretMeta
	0,  // 2: proto.Secret.type:type_name -> proto.SecretType
	0,  // 3: proto.SecretRequest.type:type_name -> proto.SecretType
	0,  // 4: proto.SecretCreateRequest.type:type_name -> proto.SecretType
	0,  // 5: proto.SecretUpdateRequest.type:type_name -> proto.SecretType
	1,  // 6: proto.Keeper.Login:input_type -> proto.LoginRequest
	2,  // 7: proto.Keeper.SignUp:input_type -> proto.SignUpRequest
	4,  // 8: proto.Keeper.ListSecrets:input_type -> proto.ListRequest
	8,  // 9: proto.Keeper.GetSecret:input_type -> proto.SecretRequest
	9,  // 10: proto.Keeper.CreateSecret:input_type -> proto.SecretCreateRequest
	10, // 11: proto.Keeper.UpdateSecret:input_type -> proto.SecretUpdateRequest
	11, // 12: proto.Keeper.CreateFile:input_type -> proto.SecretFileCreateRequest
	3,  // 13: proto.Keeper.Login:output_type -> proto.Token
	3,  // 14: proto.Keeper.SignUp:output_type -> proto.Token
	6,  // 15: proto.Keeper.ListSecrets:output_type -> proto.ListResponse
	7,  // 16: proto.Keeper.GetSecret:output_type -> proto.Secret
	7,  // 17: proto.Keeper.CreateSecret:output_type -> proto.Secret
	7,  // 18: proto.Keeper.UpdateSecret:output_type -> proto.Secret
	12, // 19: proto.Keeper.CreateFile:output_type -> proto.SecretFileCreateResponse
	13, // [13:20] is the sub-list for method output_type
	6,  // [6:13] is the sub-list for method input_type
	6,  // [6:6] is the sub-list for extension type_name
	6,  // [6:6] is the sub-list for extension extendee
	0,  // [0:6] is the sub-list for field type_name
}

func init() { file_internal_proto_keeper_proto_init() }
func file_internal_proto_keeper_proto_init() {
	if File_internal_proto_keeper_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_internal_proto_keeper_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*LoginRequest); i {
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
		file_internal_proto_keeper_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*SignUpRequest); i {
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
		file_internal_proto_keeper_proto_msgTypes[2].Exporter = func(v any, i int) any {
			switch v := v.(*Token); i {
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
		file_internal_proto_keeper_proto_msgTypes[3].Exporter = func(v any, i int) any {
			switch v := v.(*ListRequest); i {
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
		file_internal_proto_keeper_proto_msgTypes[4].Exporter = func(v any, i int) any {
			switch v := v.(*SecretMeta); i {
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
		file_internal_proto_keeper_proto_msgTypes[5].Exporter = func(v any, i int) any {
			switch v := v.(*ListResponse); i {
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
		file_internal_proto_keeper_proto_msgTypes[6].Exporter = func(v any, i int) any {
			switch v := v.(*Secret); i {
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
		file_internal_proto_keeper_proto_msgTypes[7].Exporter = func(v any, i int) any {
			switch v := v.(*SecretRequest); i {
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
		file_internal_proto_keeper_proto_msgTypes[8].Exporter = func(v any, i int) any {
			switch v := v.(*SecretCreateRequest); i {
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
		file_internal_proto_keeper_proto_msgTypes[9].Exporter = func(v any, i int) any {
			switch v := v.(*SecretUpdateRequest); i {
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
		file_internal_proto_keeper_proto_msgTypes[10].Exporter = func(v any, i int) any {
			switch v := v.(*SecretFileCreateRequest); i {
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
		file_internal_proto_keeper_proto_msgTypes[11].Exporter = func(v any, i int) any {
			switch v := v.(*SecretFileCreateResponse); i {
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
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_internal_proto_keeper_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   12,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_internal_proto_keeper_proto_goTypes,
		DependencyIndexes: file_internal_proto_keeper_proto_depIdxs,
		EnumInfos:         file_internal_proto_keeper_proto_enumTypes,
		MessageInfos:      file_internal_proto_keeper_proto_msgTypes,
	}.Build()
	File_internal_proto_keeper_proto = out.File
	file_internal_proto_keeper_proto_rawDesc = nil
	file_internal_proto_keeper_proto_goTypes = nil
	file_internal_proto_keeper_proto_depIdxs = nil
}
