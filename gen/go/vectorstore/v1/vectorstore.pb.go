// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.1
// 	protoc        (unknown)
// source: vectorstore/v1/vectorstore.proto

package vectorstorev1

import (
	_ "google.golang.org/genproto/googleapis/api/annotations"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	structpb "google.golang.org/protobuf/types/known/structpb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type VectorStoreServiceInsertTextsRequestText struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Text     string           `protobuf:"bytes,1,opt,name=text,proto3" json:"text,omitempty"`
	Metadata *structpb.Struct `protobuf:"bytes,2,opt,name=metadata,proto3" json:"metadata,omitempty"`
}

func (x *VectorStoreServiceInsertTextsRequestText) Reset() {
	*x = VectorStoreServiceInsertTextsRequestText{}
	mi := &file_vectorstore_v1_vectorstore_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *VectorStoreServiceInsertTextsRequestText) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VectorStoreServiceInsertTextsRequestText) ProtoMessage() {}

func (x *VectorStoreServiceInsertTextsRequestText) ProtoReflect() protoreflect.Message {
	mi := &file_vectorstore_v1_vectorstore_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use VectorStoreServiceInsertTextsRequestText.ProtoReflect.Descriptor instead.
func (*VectorStoreServiceInsertTextsRequestText) Descriptor() ([]byte, []int) {
	return file_vectorstore_v1_vectorstore_proto_rawDescGZIP(), []int{0}
}

func (x *VectorStoreServiceInsertTextsRequestText) GetText() string {
	if x != nil {
		return x.Text
	}
	return ""
}

func (x *VectorStoreServiceInsertTextsRequestText) GetMetadata() *structpb.Struct {
	if x != nil {
		return x.Metadata
	}
	return nil
}

type VectorStoreServiceInsertTextsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Texts []*VectorStoreServiceInsertTextsRequestText `protobuf:"bytes,1,rep,name=texts,proto3" json:"texts,omitempty"`
}

func (x *VectorStoreServiceInsertTextsRequest) Reset() {
	*x = VectorStoreServiceInsertTextsRequest{}
	mi := &file_vectorstore_v1_vectorstore_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *VectorStoreServiceInsertTextsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VectorStoreServiceInsertTextsRequest) ProtoMessage() {}

func (x *VectorStoreServiceInsertTextsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_vectorstore_v1_vectorstore_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use VectorStoreServiceInsertTextsRequest.ProtoReflect.Descriptor instead.
func (*VectorStoreServiceInsertTextsRequest) Descriptor() ([]byte, []int) {
	return file_vectorstore_v1_vectorstore_proto_rawDescGZIP(), []int{1}
}

func (x *VectorStoreServiceInsertTextsRequest) GetTexts() []*VectorStoreServiceInsertTextsRequestText {
	if x != nil {
		return x.Texts
	}
	return nil
}

type VectorStoreServiceInsertTextsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *VectorStoreServiceInsertTextsResponse) Reset() {
	*x = VectorStoreServiceInsertTextsResponse{}
	mi := &file_vectorstore_v1_vectorstore_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *VectorStoreServiceInsertTextsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VectorStoreServiceInsertTextsResponse) ProtoMessage() {}

func (x *VectorStoreServiceInsertTextsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_vectorstore_v1_vectorstore_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use VectorStoreServiceInsertTextsResponse.ProtoReflect.Descriptor instead.
func (*VectorStoreServiceInsertTextsResponse) Descriptor() ([]byte, []int) {
	return file_vectorstore_v1_vectorstore_proto_rawDescGZIP(), []int{2}
}

type VectorStoreServiceSearchTextRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Text   string           `protobuf:"bytes,1,opt,name=text,proto3" json:"text,omitempty"`
	TopK   int64            `protobuf:"varint,2,opt,name=top_k,json=topK,proto3" json:"top_k,omitempty"`
	Filter *structpb.Struct `protobuf:"bytes,3,opt,name=filter,proto3" json:"filter,omitempty"`
}

func (x *VectorStoreServiceSearchTextRequest) Reset() {
	*x = VectorStoreServiceSearchTextRequest{}
	mi := &file_vectorstore_v1_vectorstore_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *VectorStoreServiceSearchTextRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VectorStoreServiceSearchTextRequest) ProtoMessage() {}

func (x *VectorStoreServiceSearchTextRequest) ProtoReflect() protoreflect.Message {
	mi := &file_vectorstore_v1_vectorstore_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use VectorStoreServiceSearchTextRequest.ProtoReflect.Descriptor instead.
func (*VectorStoreServiceSearchTextRequest) Descriptor() ([]byte, []int) {
	return file_vectorstore_v1_vectorstore_proto_rawDescGZIP(), []int{3}
}

func (x *VectorStoreServiceSearchTextRequest) GetText() string {
	if x != nil {
		return x.Text
	}
	return ""
}

func (x *VectorStoreServiceSearchTextRequest) GetTopK() int64 {
	if x != nil {
		return x.TopK
	}
	return 0
}

func (x *VectorStoreServiceSearchTextRequest) GetFilter() *structpb.Struct {
	if x != nil {
		return x.Filter
	}
	return nil
}

type VectorStoreServiceSearchTextResponseSimilarText struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Text     string           `protobuf:"bytes,1,opt,name=text,proto3" json:"text,omitempty"`
	Score    float32          `protobuf:"fixed32,2,opt,name=score,proto3" json:"score,omitempty"`
	Metadata *structpb.Struct `protobuf:"bytes,3,opt,name=metadata,proto3" json:"metadata,omitempty"`
}

func (x *VectorStoreServiceSearchTextResponseSimilarText) Reset() {
	*x = VectorStoreServiceSearchTextResponseSimilarText{}
	mi := &file_vectorstore_v1_vectorstore_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *VectorStoreServiceSearchTextResponseSimilarText) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VectorStoreServiceSearchTextResponseSimilarText) ProtoMessage() {}

func (x *VectorStoreServiceSearchTextResponseSimilarText) ProtoReflect() protoreflect.Message {
	mi := &file_vectorstore_v1_vectorstore_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use VectorStoreServiceSearchTextResponseSimilarText.ProtoReflect.Descriptor instead.
func (*VectorStoreServiceSearchTextResponseSimilarText) Descriptor() ([]byte, []int) {
	return file_vectorstore_v1_vectorstore_proto_rawDescGZIP(), []int{4}
}

func (x *VectorStoreServiceSearchTextResponseSimilarText) GetText() string {
	if x != nil {
		return x.Text
	}
	return ""
}

func (x *VectorStoreServiceSearchTextResponseSimilarText) GetScore() float32 {
	if x != nil {
		return x.Score
	}
	return 0
}

func (x *VectorStoreServiceSearchTextResponseSimilarText) GetMetadata() *structpb.Struct {
	if x != nil {
		return x.Metadata
	}
	return nil
}

type VectorStoreServiceSearchTextResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SimilarTexts []*VectorStoreServiceSearchTextResponseSimilarText `protobuf:"bytes,1,rep,name=similar_texts,json=similarTexts,proto3" json:"similar_texts,omitempty"`
}

func (x *VectorStoreServiceSearchTextResponse) Reset() {
	*x = VectorStoreServiceSearchTextResponse{}
	mi := &file_vectorstore_v1_vectorstore_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *VectorStoreServiceSearchTextResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VectorStoreServiceSearchTextResponse) ProtoMessage() {}

func (x *VectorStoreServiceSearchTextResponse) ProtoReflect() protoreflect.Message {
	mi := &file_vectorstore_v1_vectorstore_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use VectorStoreServiceSearchTextResponse.ProtoReflect.Descriptor instead.
func (*VectorStoreServiceSearchTextResponse) Descriptor() ([]byte, []int) {
	return file_vectorstore_v1_vectorstore_proto_rawDescGZIP(), []int{5}
}

func (x *VectorStoreServiceSearchTextResponse) GetSimilarTexts() []*VectorStoreServiceSearchTextResponseSimilarText {
	if x != nil {
		return x.SimilarTexts
	}
	return nil
}

var File_vectorstore_v1_vectorstore_proto protoreflect.FileDescriptor

var file_vectorstore_v1_vectorstore_proto_rawDesc = []byte{
	0x0a, 0x20, 0x76, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2f, 0x76, 0x31,
	0x2f, 0x76, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x0e, 0x76, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e,
	0x76, 0x31, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61,
	0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2f, 0x73, 0x74, 0x72, 0x75, 0x63, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x73,
	0x0a, 0x28, 0x56, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x53, 0x74, 0x6f, 0x72, 0x65, 0x53, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x49, 0x6e, 0x73, 0x65, 0x72, 0x74, 0x54, 0x65, 0x78, 0x74, 0x73, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x54, 0x65, 0x78, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x65,
	0x78, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x74, 0x65, 0x78, 0x74, 0x12, 0x33,
	0x0a, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x17, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x75, 0x63, 0x74, 0x52, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x64,
	0x61, 0x74, 0x61, 0x22, 0x76, 0x0a, 0x24, 0x56, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x53, 0x74, 0x6f,
	0x72, 0x65, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x49, 0x6e, 0x73, 0x65, 0x72, 0x74, 0x54,
	0x65, 0x78, 0x74, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x4e, 0x0a, 0x05, 0x74,
	0x65, 0x78, 0x74, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x38, 0x2e, 0x76, 0x65, 0x63,
	0x74, 0x6f, 0x72, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x56, 0x65, 0x63, 0x74,
	0x6f, 0x72, 0x53, 0x74, 0x6f, 0x72, 0x65, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x49, 0x6e,
	0x73, 0x65, 0x72, 0x74, 0x54, 0x65, 0x78, 0x74, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x54, 0x65, 0x78, 0x74, 0x52, 0x05, 0x74, 0x65, 0x78, 0x74, 0x73, 0x22, 0x27, 0x0a, 0x25, 0x56,
	0x65, 0x63, 0x74, 0x6f, 0x72, 0x53, 0x74, 0x6f, 0x72, 0x65, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x49, 0x6e, 0x73, 0x65, 0x72, 0x74, 0x54, 0x65, 0x78, 0x74, 0x73, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x22, 0x7f, 0x0a, 0x23, 0x56, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x53, 0x74,
	0x6f, 0x72, 0x65, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x53, 0x65, 0x61, 0x72, 0x63, 0x68,
	0x54, 0x65, 0x78, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x74,
	0x65, 0x78, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x74, 0x65, 0x78, 0x74, 0x12,
	0x13, 0x0a, 0x05, 0x74, 0x6f, 0x70, 0x5f, 0x6b, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x04,
	0x74, 0x6f, 0x70, 0x4b, 0x12, 0x2f, 0x0a, 0x06, 0x66, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x75, 0x63, 0x74, 0x52, 0x06, 0x66,
	0x69, 0x6c, 0x74, 0x65, 0x72, 0x22, 0x90, 0x01, 0x0a, 0x2f, 0x56, 0x65, 0x63, 0x74, 0x6f, 0x72,
	0x53, 0x74, 0x6f, 0x72, 0x65, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x53, 0x65, 0x61, 0x72,
	0x63, 0x68, 0x54, 0x65, 0x78, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x53, 0x69,
	0x6d, 0x69, 0x6c, 0x61, 0x72, 0x54, 0x65, 0x78, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x65, 0x78,
	0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x74, 0x65, 0x78, 0x74, 0x12, 0x14, 0x0a,
	0x05, 0x73, 0x63, 0x6f, 0x72, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x02, 0x52, 0x05, 0x73, 0x63,
	0x6f, 0x72, 0x65, 0x12, 0x33, 0x0a, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x75, 0x63, 0x74, 0x52, 0x08,
	0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x22, 0x8c, 0x01, 0x0a, 0x24, 0x56, 0x65, 0x63,
	0x74, 0x6f, 0x72, 0x53, 0x74, 0x6f, 0x72, 0x65, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x53,
	0x65, 0x61, 0x72, 0x63, 0x68, 0x54, 0x65, 0x78, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x64, 0x0a, 0x0d, 0x73, 0x69, 0x6d, 0x69, 0x6c, 0x61, 0x72, 0x5f, 0x74, 0x65, 0x78,
	0x74, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x3f, 0x2e, 0x76, 0x65, 0x63, 0x74, 0x6f,
	0x72, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x56, 0x65, 0x63, 0x74, 0x6f, 0x72,
	0x53, 0x74, 0x6f, 0x72, 0x65, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x53, 0x65, 0x61, 0x72,
	0x63, 0x68, 0x54, 0x65, 0x78, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x53, 0x69,
	0x6d, 0x69, 0x6c, 0x61, 0x72, 0x54, 0x65, 0x78, 0x74, 0x52, 0x0c, 0x73, 0x69, 0x6d, 0x69, 0x6c,
	0x61, 0x72, 0x54, 0x65, 0x78, 0x74, 0x73, 0x32, 0xcc, 0x02, 0x0a, 0x12, 0x56, 0x65, 0x63, 0x74,
	0x6f, 0x72, 0x53, 0x74, 0x6f, 0x72, 0x65, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x9b,
	0x01, 0x0a, 0x0b, 0x49, 0x6e, 0x73, 0x65, 0x72, 0x74, 0x54, 0x65, 0x78, 0x74, 0x73, 0x12, 0x34,
	0x2e, 0x76, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x76, 0x31, 0x2e,
	0x56, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x53, 0x74, 0x6f, 0x72, 0x65, 0x53, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x49, 0x6e, 0x73, 0x65, 0x72, 0x74, 0x54, 0x65, 0x78, 0x74, 0x73, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x35, 0x2e, 0x76, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x73, 0x74, 0x6f,
	0x72, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x56, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x53, 0x74, 0x6f, 0x72,
	0x65, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x49, 0x6e, 0x73, 0x65, 0x72, 0x74, 0x54, 0x65,
	0x78, 0x74, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x1f, 0x82, 0xd3, 0xe4,
	0x93, 0x02, 0x19, 0x3a, 0x01, 0x2a, 0x22, 0x14, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x2f,
	0x69, 0x6e, 0x73, 0x65, 0x72, 0x74, 0x5f, 0x74, 0x65, 0x78, 0x74, 0x73, 0x12, 0x97, 0x01, 0x0a,
	0x0a, 0x53, 0x65, 0x61, 0x72, 0x63, 0x68, 0x54, 0x65, 0x78, 0x74, 0x12, 0x33, 0x2e, 0x76, 0x65,
	0x63, 0x74, 0x6f, 0x72, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x56, 0x65, 0x63,
	0x74, 0x6f, 0x72, 0x53, 0x74, 0x6f, 0x72, 0x65, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x53,
	0x65, 0x61, 0x72, 0x63, 0x68, 0x54, 0x65, 0x78, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x34, 0x2e, 0x76, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x76,
	0x31, 0x2e, 0x56, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x53, 0x74, 0x6f, 0x72, 0x65, 0x53, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x53, 0x65, 0x61, 0x72, 0x63, 0x68, 0x54, 0x65, 0x78, 0x74, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x1e, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x18, 0x3a, 0x01,
	0x2a, 0x22, 0x13, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x2f, 0x73, 0x65, 0x61, 0x72, 0x63,
	0x68, 0x5f, 0x74, 0x65, 0x78, 0x74, 0x42, 0x44, 0x5a, 0x42, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62,
	0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x61, 0x72, 0x69, 0x61, 0x33, 0x70, 0x70, 0x70, 0x2f, 0x72, 0x61,
	0x67, 0x2d, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x67, 0x6f, 0x2f,
	0x76, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2f, 0x76, 0x31, 0x3b, 0x76,
	0x65, 0x63, 0x74, 0x6f, 0x72, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_vectorstore_v1_vectorstore_proto_rawDescOnce sync.Once
	file_vectorstore_v1_vectorstore_proto_rawDescData = file_vectorstore_v1_vectorstore_proto_rawDesc
)

func file_vectorstore_v1_vectorstore_proto_rawDescGZIP() []byte {
	file_vectorstore_v1_vectorstore_proto_rawDescOnce.Do(func() {
		file_vectorstore_v1_vectorstore_proto_rawDescData = protoimpl.X.CompressGZIP(file_vectorstore_v1_vectorstore_proto_rawDescData)
	})
	return file_vectorstore_v1_vectorstore_proto_rawDescData
}

var file_vectorstore_v1_vectorstore_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_vectorstore_v1_vectorstore_proto_goTypes = []any{
	(*VectorStoreServiceInsertTextsRequestText)(nil),        // 0: vectorstore.v1.VectorStoreServiceInsertTextsRequestText
	(*VectorStoreServiceInsertTextsRequest)(nil),            // 1: vectorstore.v1.VectorStoreServiceInsertTextsRequest
	(*VectorStoreServiceInsertTextsResponse)(nil),           // 2: vectorstore.v1.VectorStoreServiceInsertTextsResponse
	(*VectorStoreServiceSearchTextRequest)(nil),             // 3: vectorstore.v1.VectorStoreServiceSearchTextRequest
	(*VectorStoreServiceSearchTextResponseSimilarText)(nil), // 4: vectorstore.v1.VectorStoreServiceSearchTextResponseSimilarText
	(*VectorStoreServiceSearchTextResponse)(nil),            // 5: vectorstore.v1.VectorStoreServiceSearchTextResponse
	(*structpb.Struct)(nil),                                 // 6: google.protobuf.Struct
}
var file_vectorstore_v1_vectorstore_proto_depIdxs = []int32{
	6, // 0: vectorstore.v1.VectorStoreServiceInsertTextsRequestText.metadata:type_name -> google.protobuf.Struct
	0, // 1: vectorstore.v1.VectorStoreServiceInsertTextsRequest.texts:type_name -> vectorstore.v1.VectorStoreServiceInsertTextsRequestText
	6, // 2: vectorstore.v1.VectorStoreServiceSearchTextRequest.filter:type_name -> google.protobuf.Struct
	6, // 3: vectorstore.v1.VectorStoreServiceSearchTextResponseSimilarText.metadata:type_name -> google.protobuf.Struct
	4, // 4: vectorstore.v1.VectorStoreServiceSearchTextResponse.similar_texts:type_name -> vectorstore.v1.VectorStoreServiceSearchTextResponseSimilarText
	1, // 5: vectorstore.v1.VectorStoreService.InsertTexts:input_type -> vectorstore.v1.VectorStoreServiceInsertTextsRequest
	3, // 6: vectorstore.v1.VectorStoreService.SearchText:input_type -> vectorstore.v1.VectorStoreServiceSearchTextRequest
	2, // 7: vectorstore.v1.VectorStoreService.InsertTexts:output_type -> vectorstore.v1.VectorStoreServiceInsertTextsResponse
	5, // 8: vectorstore.v1.VectorStoreService.SearchText:output_type -> vectorstore.v1.VectorStoreServiceSearchTextResponse
	7, // [7:9] is the sub-list for method output_type
	5, // [5:7] is the sub-list for method input_type
	5, // [5:5] is the sub-list for extension type_name
	5, // [5:5] is the sub-list for extension extendee
	0, // [0:5] is the sub-list for field type_name
}

func init() { file_vectorstore_v1_vectorstore_proto_init() }
func file_vectorstore_v1_vectorstore_proto_init() {
	if File_vectorstore_v1_vectorstore_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_vectorstore_v1_vectorstore_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_vectorstore_v1_vectorstore_proto_goTypes,
		DependencyIndexes: file_vectorstore_v1_vectorstore_proto_depIdxs,
		MessageInfos:      file_vectorstore_v1_vectorstore_proto_msgTypes,
	}.Build()
	File_vectorstore_v1_vectorstore_proto = out.File
	file_vectorstore_v1_vectorstore_proto_rawDesc = nil
	file_vectorstore_v1_vectorstore_proto_goTypes = nil
	file_vectorstore_v1_vectorstore_proto_depIdxs = nil
}