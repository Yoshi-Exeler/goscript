// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.21.5
// source: goscript.proto

package encoding

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	anypb "google.golang.org/protobuf/types/known/anypb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Expression struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Left     *Expression       `protobuf:"bytes,1,opt,name=Left,proto3" json:"Left,omitempty"`
	Right    *Expression       `protobuf:"bytes,2,opt,name=Right,proto3" json:"Right,omitempty"`
	Value    *BinaryTypedValue `protobuf:"bytes,3,opt,name=Value,proto3" json:"Value,omitempty"`
	Ref      uint64            `protobuf:"varint,4,opt,name=Ref,proto3" json:"Ref,omitempty"`
	Operator uint32            `protobuf:"varint,5,opt,name=Operator,proto3" json:"Operator,omitempty"`
	Args     []*anypb.Any      `protobuf:"bytes,6,rep,name=Args,proto3" json:"Args,omitempty"`
}

func (x *Expression) Reset() {
	*x = Expression{}
	if protoimpl.UnsafeEnabled {
		mi := &file_goscript_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Expression) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Expression) ProtoMessage() {}

func (x *Expression) ProtoReflect() protoreflect.Message {
	mi := &file_goscript_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Expression.ProtoReflect.Descriptor instead.
func (*Expression) Descriptor() ([]byte, []int) {
	return file_goscript_proto_rawDescGZIP(), []int{0}
}

func (x *Expression) GetLeft() *Expression {
	if x != nil {
		return x.Left
	}
	return nil
}

func (x *Expression) GetRight() *Expression {
	if x != nil {
		return x.Right
	}
	return nil
}

func (x *Expression) GetValue() *BinaryTypedValue {
	if x != nil {
		return x.Value
	}
	return nil
}

func (x *Expression) GetRef() uint64 {
	if x != nil {
		return x.Ref
	}
	return 0
}

func (x *Expression) GetOperator() uint32 {
	if x != nil {
		return x.Operator
	}
	return 0
}

func (x *Expression) GetArgs() []*anypb.Any {
	if x != nil {
		return x.Args
	}
	return nil
}

type BinaryTypedValue struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type  uint32     `protobuf:"varint,1,opt,name=Type,proto3" json:"Type,omitempty"`
	Value *anypb.Any `protobuf:"bytes,2,opt,name=Value,proto3" json:"Value,omitempty"`
}

func (x *BinaryTypedValue) Reset() {
	*x = BinaryTypedValue{}
	if protoimpl.UnsafeEnabled {
		mi := &file_goscript_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BinaryTypedValue) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BinaryTypedValue) ProtoMessage() {}

func (x *BinaryTypedValue) ProtoReflect() protoreflect.Message {
	mi := &file_goscript_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BinaryTypedValue.ProtoReflect.Descriptor instead.
func (*BinaryTypedValue) Descriptor() ([]byte, []int) {
	return file_goscript_proto_rawDescGZIP(), []int{1}
}

func (x *BinaryTypedValue) GetType() uint32 {
	if x != nil {
		return x.Type
	}
	return 0
}

func (x *BinaryTypedValue) GetValue() *anypb.Any {
	if x != nil {
		return x.Value
	}
	return nil
}

type BinaryOperation struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type uint32       `protobuf:"varint,1,opt,name=Type,proto3" json:"Type,omitempty"`
	Args []*anypb.Any `protobuf:"bytes,2,rep,name=Args,proto3" json:"Args,omitempty"`
}

func (x *BinaryOperation) Reset() {
	*x = BinaryOperation{}
	if protoimpl.UnsafeEnabled {
		mi := &file_goscript_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BinaryOperation) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BinaryOperation) ProtoMessage() {}

func (x *BinaryOperation) ProtoReflect() protoreflect.Message {
	mi := &file_goscript_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BinaryOperation.ProtoReflect.Descriptor instead.
func (*BinaryOperation) Descriptor() ([]byte, []int) {
	return file_goscript_proto_rawDescGZIP(), []int{2}
}

func (x *BinaryOperation) GetType() uint32 {
	if x != nil {
		return x.Type
	}
	return 0
}

func (x *BinaryOperation) GetArgs() []*anypb.Any {
	if x != nil {
		return x.Args
	}
	return nil
}

type Program struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SymbolTableSize uint64             `protobuf:"varint,1,opt,name=SymbolTableSize,proto3" json:"SymbolTableSize,omitempty"`
	Operations      []*BinaryOperation `protobuf:"bytes,2,rep,name=Operations,proto3" json:"Operations,omitempty"`
}

func (x *Program) Reset() {
	*x = Program{}
	if protoimpl.UnsafeEnabled {
		mi := &file_goscript_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Program) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Program) ProtoMessage() {}

func (x *Program) ProtoReflect() protoreflect.Message {
	mi := &file_goscript_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Program.ProtoReflect.Descriptor instead.
func (*Program) Descriptor() ([]byte, []int) {
	return file_goscript_proto_rawDescGZIP(), []int{3}
}

func (x *Program) GetSymbolTableSize() uint64 {
	if x != nil {
		return x.SymbolTableSize
	}
	return 0
}

func (x *Program) GetOperations() []*BinaryOperation {
	if x != nil {
		return x.Operations
	}
	return nil
}

type FunctionArgument struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Expression *Expression `protobuf:"bytes,1,opt,name=Expression,proto3" json:"Expression,omitempty"`
	SymbolRef  uint64      `protobuf:"varint,2,opt,name=SymbolRef,proto3" json:"SymbolRef,omitempty"`
}

func (x *FunctionArgument) Reset() {
	*x = FunctionArgument{}
	if protoimpl.UnsafeEnabled {
		mi := &file_goscript_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FunctionArgument) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FunctionArgument) ProtoMessage() {}

func (x *FunctionArgument) ProtoReflect() protoreflect.Message {
	mi := &file_goscript_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FunctionArgument.ProtoReflect.Descriptor instead.
func (*FunctionArgument) Descriptor() ([]byte, []int) {
	return file_goscript_proto_rawDescGZIP(), []int{4}
}

func (x *FunctionArgument) GetExpression() *Expression {
	if x != nil {
		return x.Expression
	}
	return nil
}

func (x *FunctionArgument) GetSymbolRef() uint64 {
	if x != nil {
		return x.SymbolRef
	}
	return 0
}

type U64Container struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Value uint64 `protobuf:"varint,1,opt,name=Value,proto3" json:"Value,omitempty"`
}

func (x *U64Container) Reset() {
	*x = U64Container{}
	if protoimpl.UnsafeEnabled {
		mi := &file_goscript_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *U64Container) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*U64Container) ProtoMessage() {}

func (x *U64Container) ProtoReflect() protoreflect.Message {
	mi := &file_goscript_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use U64Container.ProtoReflect.Descriptor instead.
func (*U64Container) Descriptor() ([]byte, []int) {
	return file_goscript_proto_rawDescGZIP(), []int{5}
}

func (x *U64Container) GetValue() uint64 {
	if x != nil {
		return x.Value
	}
	return 0
}

type StringContainer struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Value string `protobuf:"bytes,1,opt,name=Value,proto3" json:"Value,omitempty"`
}

func (x *StringContainer) Reset() {
	*x = StringContainer{}
	if protoimpl.UnsafeEnabled {
		mi := &file_goscript_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StringContainer) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StringContainer) ProtoMessage() {}

func (x *StringContainer) ProtoReflect() protoreflect.Message {
	mi := &file_goscript_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StringContainer.ProtoReflect.Descriptor instead.
func (*StringContainer) Descriptor() ([]byte, []int) {
	return file_goscript_proto_rawDescGZIP(), []int{6}
}

func (x *StringContainer) GetValue() string {
	if x != nil {
		return x.Value
	}
	return ""
}

type F64Container struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Value float64 `protobuf:"fixed64,1,opt,name=Value,proto3" json:"Value,omitempty"`
}

func (x *F64Container) Reset() {
	*x = F64Container{}
	if protoimpl.UnsafeEnabled {
		mi := &file_goscript_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *F64Container) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*F64Container) ProtoMessage() {}

func (x *F64Container) ProtoReflect() protoreflect.Message {
	mi := &file_goscript_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use F64Container.ProtoReflect.Descriptor instead.
func (*F64Container) Descriptor() ([]byte, []int) {
	return file_goscript_proto_rawDescGZIP(), []int{7}
}

func (x *F64Container) GetValue() float64 {
	if x != nil {
		return x.Value
	}
	return 0
}

var File_goscript_proto protoreflect.FileDescriptor

var file_goscript_proto_rawDesc = []byte{
	0x0a, 0x0e, 0x67, 0x6f, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x08, 0x65, 0x6e, 0x63, 0x6f, 0x64, 0x69, 0x6e, 0x67, 0x1a, 0x19, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x61, 0x6e, 0x79, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xec, 0x01, 0x0a, 0x0a, 0x45, 0x78, 0x70, 0x72, 0x65, 0x73,
	0x73, 0x69, 0x6f, 0x6e, 0x12, 0x28, 0x0a, 0x04, 0x4c, 0x65, 0x66, 0x74, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x14, 0x2e, 0x65, 0x6e, 0x63, 0x6f, 0x64, 0x69, 0x6e, 0x67, 0x2e, 0x45, 0x78,
	0x70, 0x72, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x04, 0x4c, 0x65, 0x66, 0x74, 0x12, 0x2a,
	0x0a, 0x05, 0x52, 0x69, 0x67, 0x68, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e,
	0x65, 0x6e, 0x63, 0x6f, 0x64, 0x69, 0x6e, 0x67, 0x2e, 0x45, 0x78, 0x70, 0x72, 0x65, 0x73, 0x73,
	0x69, 0x6f, 0x6e, 0x52, 0x05, 0x52, 0x69, 0x67, 0x68, 0x74, 0x12, 0x30, 0x0a, 0x05, 0x56, 0x61,
	0x6c, 0x75, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x65, 0x6e, 0x63, 0x6f,
	0x64, 0x69, 0x6e, 0x67, 0x2e, 0x42, 0x69, 0x6e, 0x61, 0x72, 0x79, 0x54, 0x79, 0x70, 0x65, 0x64,
	0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x05, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x10, 0x0a, 0x03,
	0x52, 0x65, 0x66, 0x18, 0x04, 0x20, 0x01, 0x28, 0x04, 0x52, 0x03, 0x52, 0x65, 0x66, 0x12, 0x1a,
	0x0a, 0x08, 0x4f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x6f, 0x72, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0d,
	0x52, 0x08, 0x4f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x6f, 0x72, 0x12, 0x28, 0x0a, 0x04, 0x41, 0x72,
	0x67, 0x73, 0x18, 0x06, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x41, 0x6e, 0x79, 0x52, 0x04,
	0x41, 0x72, 0x67, 0x73, 0x22, 0x52, 0x0a, 0x10, 0x42, 0x69, 0x6e, 0x61, 0x72, 0x79, 0x54, 0x79,
	0x70, 0x65, 0x64, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x54, 0x79, 0x70, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x04, 0x54, 0x79, 0x70, 0x65, 0x12, 0x2a, 0x0a, 0x05,
	0x56, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x41, 0x6e,
	0x79, 0x52, 0x05, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x22, 0x4f, 0x0a, 0x0f, 0x42, 0x69, 0x6e, 0x61,
	0x72, 0x79, 0x4f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x12, 0x0a, 0x04, 0x54,
	0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x04, 0x54, 0x79, 0x70, 0x65, 0x12,
	0x28, 0x0a, 0x04, 0x41, 0x72, 0x67, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x14, 0x2e,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e,
	0x41, 0x6e, 0x79, 0x52, 0x04, 0x41, 0x72, 0x67, 0x73, 0x22, 0x6e, 0x0a, 0x07, 0x50, 0x72, 0x6f,
	0x67, 0x72, 0x61, 0x6d, 0x12, 0x28, 0x0a, 0x0f, 0x53, 0x79, 0x6d, 0x62, 0x6f, 0x6c, 0x54, 0x61,
	0x62, 0x6c, 0x65, 0x53, 0x69, 0x7a, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0f, 0x53,
	0x79, 0x6d, 0x62, 0x6f, 0x6c, 0x54, 0x61, 0x62, 0x6c, 0x65, 0x53, 0x69, 0x7a, 0x65, 0x12, 0x39,
	0x0a, 0x0a, 0x4f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x02, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x19, 0x2e, 0x65, 0x6e, 0x63, 0x6f, 0x64, 0x69, 0x6e, 0x67, 0x2e, 0x42, 0x69,
	0x6e, 0x61, 0x72, 0x79, 0x4f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x0a, 0x4f,
	0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x22, 0x66, 0x0a, 0x10, 0x46, 0x75, 0x6e,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x41, 0x72, 0x67, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x34, 0x0a,
	0x0a, 0x45, 0x78, 0x70, 0x72, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x14, 0x2e, 0x65, 0x6e, 0x63, 0x6f, 0x64, 0x69, 0x6e, 0x67, 0x2e, 0x45, 0x78, 0x70,
	0x72, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x0a, 0x45, 0x78, 0x70, 0x72, 0x65, 0x73, 0x73,
	0x69, 0x6f, 0x6e, 0x12, 0x1c, 0x0a, 0x09, 0x53, 0x79, 0x6d, 0x62, 0x6f, 0x6c, 0x52, 0x65, 0x66,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x09, 0x53, 0x79, 0x6d, 0x62, 0x6f, 0x6c, 0x52, 0x65,
	0x66, 0x22, 0x24, 0x0a, 0x0c, 0x55, 0x36, 0x34, 0x43, 0x6f, 0x6e, 0x74, 0x61, 0x69, 0x6e, 0x65,
	0x72, 0x12, 0x14, 0x0a, 0x05, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04,
	0x52, 0x05, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x22, 0x27, 0x0a, 0x0f, 0x53, 0x74, 0x72, 0x69, 0x6e,
	0x67, 0x43, 0x6f, 0x6e, 0x74, 0x61, 0x69, 0x6e, 0x65, 0x72, 0x12, 0x14, 0x0a, 0x05, 0x56, 0x61,
	0x6c, 0x75, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x56, 0x61, 0x6c, 0x75, 0x65,
	0x22, 0x24, 0x0a, 0x0c, 0x46, 0x36, 0x34, 0x43, 0x6f, 0x6e, 0x74, 0x61, 0x69, 0x6e, 0x65, 0x72,
	0x12, 0x14, 0x0a, 0x05, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x01, 0x52,
	0x05, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x42, 0x19, 0x48, 0x01, 0x5a, 0x15, 0x67, 0x6f, 0x73, 0x63,
	0x72, 0x69, 0x70, 0x74, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x65, 0x6e, 0x63, 0x6f, 0x64, 0x69, 0x6e,
	0x67, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_goscript_proto_rawDescOnce sync.Once
	file_goscript_proto_rawDescData = file_goscript_proto_rawDesc
)

func file_goscript_proto_rawDescGZIP() []byte {
	file_goscript_proto_rawDescOnce.Do(func() {
		file_goscript_proto_rawDescData = protoimpl.X.CompressGZIP(file_goscript_proto_rawDescData)
	})
	return file_goscript_proto_rawDescData
}

var file_goscript_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_goscript_proto_goTypes = []interface{}{
	(*Expression)(nil),       // 0: encoding.Expression
	(*BinaryTypedValue)(nil), // 1: encoding.BinaryTypedValue
	(*BinaryOperation)(nil),  // 2: encoding.BinaryOperation
	(*Program)(nil),          // 3: encoding.Program
	(*FunctionArgument)(nil), // 4: encoding.FunctionArgument
	(*U64Container)(nil),     // 5: encoding.U64Container
	(*StringContainer)(nil),  // 6: encoding.StringContainer
	(*F64Container)(nil),     // 7: encoding.F64Container
	(*anypb.Any)(nil),        // 8: google.protobuf.Any
}
var file_goscript_proto_depIdxs = []int32{
	0, // 0: encoding.Expression.Left:type_name -> encoding.Expression
	0, // 1: encoding.Expression.Right:type_name -> encoding.Expression
	1, // 2: encoding.Expression.Value:type_name -> encoding.BinaryTypedValue
	8, // 3: encoding.Expression.Args:type_name -> google.protobuf.Any
	8, // 4: encoding.BinaryTypedValue.Value:type_name -> google.protobuf.Any
	8, // 5: encoding.BinaryOperation.Args:type_name -> google.protobuf.Any
	2, // 6: encoding.Program.Operations:type_name -> encoding.BinaryOperation
	0, // 7: encoding.FunctionArgument.Expression:type_name -> encoding.Expression
	8, // [8:8] is the sub-list for method output_type
	8, // [8:8] is the sub-list for method input_type
	8, // [8:8] is the sub-list for extension type_name
	8, // [8:8] is the sub-list for extension extendee
	0, // [0:8] is the sub-list for field type_name
}

func init() { file_goscript_proto_init() }
func file_goscript_proto_init() {
	if File_goscript_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_goscript_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Expression); i {
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
		file_goscript_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BinaryTypedValue); i {
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
		file_goscript_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BinaryOperation); i {
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
		file_goscript_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Program); i {
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
		file_goscript_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FunctionArgument); i {
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
		file_goscript_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*U64Container); i {
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
		file_goscript_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StringContainer); i {
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
		file_goscript_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*F64Container); i {
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
			RawDescriptor: file_goscript_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_goscript_proto_goTypes,
		DependencyIndexes: file_goscript_proto_depIdxs,
		MessageInfos:      file_goscript_proto_msgTypes,
	}.Build()
	File_goscript_proto = out.File
	file_goscript_proto_rawDesc = nil
	file_goscript_proto_goTypes = nil
	file_goscript_proto_depIdxs = nil
}
