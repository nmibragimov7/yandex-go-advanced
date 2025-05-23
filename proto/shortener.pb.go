// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v5.29.3
// source: shortener.proto

package shortener

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type ShortenRecord struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	ShortUrl      string                 `protobuf:"bytes,1,opt,name=short_url,json=shortUrl,proto3" json:"short_url,omitempty"`
	OriginalUrl   string                 `protobuf:"bytes,2,opt,name=original_url,json=originalUrl,proto3" json:"original_url,omitempty"`
	UserId        int64                  `protobuf:"varint,3,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	IsDeleted     bool                   `protobuf:"varint,4,opt,name=is_deleted,json=isDeleted,proto3" json:"is_deleted,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ShortenRecord) Reset() {
	*x = ShortenRecord{}
	mi := &file_shortener_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ShortenRecord) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ShortenRecord) ProtoMessage() {}

func (x *ShortenRecord) ProtoReflect() protoreflect.Message {
	mi := &file_shortener_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ShortenRecord.ProtoReflect.Descriptor instead.
func (*ShortenRecord) Descriptor() ([]byte, []int) {
	return file_shortener_proto_rawDescGZIP(), []int{0}
}

func (x *ShortenRecord) GetShortUrl() string {
	if x != nil {
		return x.ShortUrl
	}
	return ""
}

func (x *ShortenRecord) GetOriginalUrl() string {
	if x != nil {
		return x.OriginalUrl
	}
	return ""
}

func (x *ShortenRecord) GetUserId() int64 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *ShortenRecord) GetIsDeleted() bool {
	if x != nil {
		return x.IsDeleted
	}
	return false
}

type ShortenRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Url           string                 `protobuf:"bytes,1,opt,name=url,proto3" json:"url,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ShortenRequest) Reset() {
	*x = ShortenRequest{}
	mi := &file_shortener_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ShortenRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ShortenRequest) ProtoMessage() {}

func (x *ShortenRequest) ProtoReflect() protoreflect.Message {
	mi := &file_shortener_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ShortenRequest.ProtoReflect.Descriptor instead.
func (*ShortenRequest) Descriptor() ([]byte, []int) {
	return file_shortener_proto_rawDescGZIP(), []int{1}
}

func (x *ShortenRequest) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

type ShortenResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Url           string                 `protobuf:"bytes,1,opt,name=url,proto3" json:"url,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ShortenResponse) Reset() {
	*x = ShortenResponse{}
	mi := &file_shortener_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ShortenResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ShortenResponse) ProtoMessage() {}

func (x *ShortenResponse) ProtoReflect() protoreflect.Message {
	mi := &file_shortener_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ShortenResponse.ProtoReflect.Descriptor instead.
func (*ShortenResponse) Descriptor() ([]byte, []int) {
	return file_shortener_proto_rawDescGZIP(), []int{2}
}

func (x *ShortenResponse) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

type GetItemRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetItemRequest) Reset() {
	*x = GetItemRequest{}
	mi := &file_shortener_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetItemRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetItemRequest) ProtoMessage() {}

func (x *GetItemRequest) ProtoReflect() protoreflect.Message {
	mi := &file_shortener_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetItemRequest.ProtoReflect.Descriptor instead.
func (*GetItemRequest) Descriptor() ([]byte, []int) {
	return file_shortener_proto_rawDescGZIP(), []int{3}
}

func (x *GetItemRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type GetItemResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Url           string                 `protobuf:"bytes,1,opt,name=url,proto3" json:"url,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetItemResponse) Reset() {
	*x = GetItemResponse{}
	mi := &file_shortener_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetItemResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetItemResponse) ProtoMessage() {}

func (x *GetItemResponse) ProtoReflect() protoreflect.Message {
	mi := &file_shortener_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetItemResponse.ProtoReflect.Descriptor instead.
func (*GetItemResponse) Descriptor() ([]byte, []int) {
	return file_shortener_proto_rawDescGZIP(), []int{4}
}

func (x *GetItemResponse) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

type PingResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Message       string                 `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *PingResponse) Reset() {
	*x = PingResponse{}
	mi := &file_shortener_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PingResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PingResponse) ProtoMessage() {}

func (x *PingResponse) ProtoReflect() protoreflect.Message {
	mi := &file_shortener_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PingResponse.ProtoReflect.Descriptor instead.
func (*PingResponse) Descriptor() ([]byte, []int) {
	return file_shortener_proto_rawDescGZIP(), []int{5}
}

func (x *PingResponse) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

type ShortenBatchRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Items         []*ShortenBatchItem    `protobuf:"bytes,1,rep,name=items,proto3" json:"items,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ShortenBatchRequest) Reset() {
	*x = ShortenBatchRequest{}
	mi := &file_shortener_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ShortenBatchRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ShortenBatchRequest) ProtoMessage() {}

func (x *ShortenBatchRequest) ProtoReflect() protoreflect.Message {
	mi := &file_shortener_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ShortenBatchRequest.ProtoReflect.Descriptor instead.
func (*ShortenBatchRequest) Descriptor() ([]byte, []int) {
	return file_shortener_proto_rawDescGZIP(), []int{6}
}

func (x *ShortenBatchRequest) GetItems() []*ShortenBatchItem {
	if x != nil {
		return x.Items
	}
	return nil
}

type ShortenBatchItem struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	CorrelationId string                 `protobuf:"bytes,1,opt,name=correlation_id,json=correlationId,proto3" json:"correlation_id,omitempty"`
	OriginalUrl   string                 `protobuf:"bytes,2,opt,name=original_url,json=originalUrl,proto3" json:"original_url,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ShortenBatchItem) Reset() {
	*x = ShortenBatchItem{}
	mi := &file_shortener_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ShortenBatchItem) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ShortenBatchItem) ProtoMessage() {}

func (x *ShortenBatchItem) ProtoReflect() protoreflect.Message {
	mi := &file_shortener_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ShortenBatchItem.ProtoReflect.Descriptor instead.
func (*ShortenBatchItem) Descriptor() ([]byte, []int) {
	return file_shortener_proto_rawDescGZIP(), []int{7}
}

func (x *ShortenBatchItem) GetCorrelationId() string {
	if x != nil {
		return x.CorrelationId
	}
	return ""
}

func (x *ShortenBatchItem) GetOriginalUrl() string {
	if x != nil {
		return x.OriginalUrl
	}
	return ""
}

type ShortenBatchResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Results       []*ShortenBatchResult  `protobuf:"bytes,1,rep,name=results,proto3" json:"results,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ShortenBatchResponse) Reset() {
	*x = ShortenBatchResponse{}
	mi := &file_shortener_proto_msgTypes[8]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ShortenBatchResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ShortenBatchResponse) ProtoMessage() {}

func (x *ShortenBatchResponse) ProtoReflect() protoreflect.Message {
	mi := &file_shortener_proto_msgTypes[8]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ShortenBatchResponse.ProtoReflect.Descriptor instead.
func (*ShortenBatchResponse) Descriptor() ([]byte, []int) {
	return file_shortener_proto_rawDescGZIP(), []int{8}
}

func (x *ShortenBatchResponse) GetResults() []*ShortenBatchResult {
	if x != nil {
		return x.Results
	}
	return nil
}

type ShortenBatchResult struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	CorrelationId string                 `protobuf:"bytes,1,opt,name=correlation_id,json=correlationId,proto3" json:"correlation_id,omitempty"`
	ShortUrl      string                 `protobuf:"bytes,2,opt,name=short_url,json=shortUrl,proto3" json:"short_url,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ShortenBatchResult) Reset() {
	*x = ShortenBatchResult{}
	mi := &file_shortener_proto_msgTypes[9]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ShortenBatchResult) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ShortenBatchResult) ProtoMessage() {}

func (x *ShortenBatchResult) ProtoReflect() protoreflect.Message {
	mi := &file_shortener_proto_msgTypes[9]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ShortenBatchResult.ProtoReflect.Descriptor instead.
func (*ShortenBatchResult) Descriptor() ([]byte, []int) {
	return file_shortener_proto_rawDescGZIP(), []int{9}
}

func (x *ShortenBatchResult) GetCorrelationId() string {
	if x != nil {
		return x.CorrelationId
	}
	return ""
}

func (x *ShortenBatchResult) GetShortUrl() string {
	if x != nil {
		return x.ShortUrl
	}
	return ""
}

type UserUrlsResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Results       []*UserUrlsResult      `protobuf:"bytes,1,rep,name=results,proto3" json:"results,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UserUrlsResponse) Reset() {
	*x = UserUrlsResponse{}
	mi := &file_shortener_proto_msgTypes[10]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UserUrlsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserUrlsResponse) ProtoMessage() {}

func (x *UserUrlsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_shortener_proto_msgTypes[10]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UserUrlsResponse.ProtoReflect.Descriptor instead.
func (*UserUrlsResponse) Descriptor() ([]byte, []int) {
	return file_shortener_proto_rawDescGZIP(), []int{10}
}

func (x *UserUrlsResponse) GetResults() []*UserUrlsResult {
	if x != nil {
		return x.Results
	}
	return nil
}

type UserUrlsResult struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	ShortUrl      string                 `protobuf:"bytes,1,opt,name=short_url,json=shortUrl,proto3" json:"short_url,omitempty"`
	OriginalUrl   string                 `protobuf:"bytes,2,opt,name=original_url,json=originalUrl,proto3" json:"original_url,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UserUrlsResult) Reset() {
	*x = UserUrlsResult{}
	mi := &file_shortener_proto_msgTypes[11]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UserUrlsResult) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserUrlsResult) ProtoMessage() {}

func (x *UserUrlsResult) ProtoReflect() protoreflect.Message {
	mi := &file_shortener_proto_msgTypes[11]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UserUrlsResult.ProtoReflect.Descriptor instead.
func (*UserUrlsResult) Descriptor() ([]byte, []int) {
	return file_shortener_proto_rawDescGZIP(), []int{11}
}

func (x *UserUrlsResult) GetShortUrl() string {
	if x != nil {
		return x.ShortUrl
	}
	return ""
}

func (x *UserUrlsResult) GetOriginalUrl() string {
	if x != nil {
		return x.OriginalUrl
	}
	return ""
}

type UserUrlsDeleteRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Items         []*UserUrlsDeleteItem  `protobuf:"bytes,1,rep,name=items,proto3" json:"items,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UserUrlsDeleteRequest) Reset() {
	*x = UserUrlsDeleteRequest{}
	mi := &file_shortener_proto_msgTypes[12]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UserUrlsDeleteRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserUrlsDeleteRequest) ProtoMessage() {}

func (x *UserUrlsDeleteRequest) ProtoReflect() protoreflect.Message {
	mi := &file_shortener_proto_msgTypes[12]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UserUrlsDeleteRequest.ProtoReflect.Descriptor instead.
func (*UserUrlsDeleteRequest) Descriptor() ([]byte, []int) {
	return file_shortener_proto_rawDescGZIP(), []int{12}
}

func (x *UserUrlsDeleteRequest) GetItems() []*UserUrlsDeleteItem {
	if x != nil {
		return x.Items
	}
	return nil
}

type UserUrlsDeleteItem struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	ShortUrl      string                 `protobuf:"bytes,1,opt,name=short_url,json=shortUrl,proto3" json:"short_url,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UserUrlsDeleteItem) Reset() {
	*x = UserUrlsDeleteItem{}
	mi := &file_shortener_proto_msgTypes[13]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UserUrlsDeleteItem) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserUrlsDeleteItem) ProtoMessage() {}

func (x *UserUrlsDeleteItem) ProtoReflect() protoreflect.Message {
	mi := &file_shortener_proto_msgTypes[13]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UserUrlsDeleteItem.ProtoReflect.Descriptor instead.
func (*UserUrlsDeleteItem) Descriptor() ([]byte, []int) {
	return file_shortener_proto_rawDescGZIP(), []int{13}
}

func (x *UserUrlsDeleteItem) GetShortUrl() string {
	if x != nil {
		return x.ShortUrl
	}
	return ""
}

type ShortenBatchUpdateRecord struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	ShortUrl      string                 `protobuf:"bytes,1,opt,name=short_url,json=shortUrl,proto3" json:"short_url,omitempty"`
	UserId        int64                  `protobuf:"varint,2,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ShortenBatchUpdateRecord) Reset() {
	*x = ShortenBatchUpdateRecord{}
	mi := &file_shortener_proto_msgTypes[14]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ShortenBatchUpdateRecord) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ShortenBatchUpdateRecord) ProtoMessage() {}

func (x *ShortenBatchUpdateRecord) ProtoReflect() protoreflect.Message {
	mi := &file_shortener_proto_msgTypes[14]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ShortenBatchUpdateRecord.ProtoReflect.Descriptor instead.
func (*ShortenBatchUpdateRecord) Descriptor() ([]byte, []int) {
	return file_shortener_proto_rawDescGZIP(), []int{14}
}

func (x *ShortenBatchUpdateRecord) GetShortUrl() string {
	if x != nil {
		return x.ShortUrl
	}
	return ""
}

func (x *ShortenBatchUpdateRecord) GetUserId() int64 {
	if x != nil {
		return x.UserId
	}
	return 0
}

type UserUrlsDeleteResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Message       string                 `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UserUrlsDeleteResponse) Reset() {
	*x = UserUrlsDeleteResponse{}
	mi := &file_shortener_proto_msgTypes[15]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UserUrlsDeleteResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserUrlsDeleteResponse) ProtoMessage() {}

func (x *UserUrlsDeleteResponse) ProtoReflect() protoreflect.Message {
	mi := &file_shortener_proto_msgTypes[15]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UserUrlsDeleteResponse.ProtoReflect.Descriptor instead.
func (*UserUrlsDeleteResponse) Descriptor() ([]byte, []int) {
	return file_shortener_proto_rawDescGZIP(), []int{15}
}

func (x *UserUrlsDeleteResponse) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

type TrustedSubnetRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	XRealIp       string                 `protobuf:"bytes,1,opt,name=x_real_ip,json=xRealIp,proto3" json:"x_real_ip,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *TrustedSubnetRequest) Reset() {
	*x = TrustedSubnetRequest{}
	mi := &file_shortener_proto_msgTypes[16]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *TrustedSubnetRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TrustedSubnetRequest) ProtoMessage() {}

func (x *TrustedSubnetRequest) ProtoReflect() protoreflect.Message {
	mi := &file_shortener_proto_msgTypes[16]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TrustedSubnetRequest.ProtoReflect.Descriptor instead.
func (*TrustedSubnetRequest) Descriptor() ([]byte, []int) {
	return file_shortener_proto_rawDescGZIP(), []int{16}
}

func (x *TrustedSubnetRequest) GetXRealIp() string {
	if x != nil {
		return x.XRealIp
	}
	return ""
}

type TrustedSubnetResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Users         int64                  `protobuf:"varint,1,opt,name=users,proto3" json:"users,omitempty"`
	Urls          int64                  `protobuf:"varint,2,opt,name=urls,proto3" json:"urls,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *TrustedSubnetResponse) Reset() {
	*x = TrustedSubnetResponse{}
	mi := &file_shortener_proto_msgTypes[17]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *TrustedSubnetResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TrustedSubnetResponse) ProtoMessage() {}

func (x *TrustedSubnetResponse) ProtoReflect() protoreflect.Message {
	mi := &file_shortener_proto_msgTypes[17]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TrustedSubnetResponse.ProtoReflect.Descriptor instead.
func (*TrustedSubnetResponse) Descriptor() ([]byte, []int) {
	return file_shortener_proto_rawDescGZIP(), []int{17}
}

func (x *TrustedSubnetResponse) GetUsers() int64 {
	if x != nil {
		return x.Users
	}
	return 0
}

func (x *TrustedSubnetResponse) GetUrls() int64 {
	if x != nil {
		return x.Urls
	}
	return 0
}

var File_shortener_proto protoreflect.FileDescriptor

const file_shortener_proto_rawDesc = "" +
	"\n" +
	"\x0fshortener.proto\x12\tshortener\x1a\x1bgoogle/protobuf/empty.proto\"\x87\x01\n" +
	"\rShortenRecord\x12\x1b\n" +
	"\tshort_url\x18\x01 \x01(\tR\bshortUrl\x12!\n" +
	"\foriginal_url\x18\x02 \x01(\tR\voriginalUrl\x12\x17\n" +
	"\auser_id\x18\x03 \x01(\x03R\x06userId\x12\x1d\n" +
	"\n" +
	"is_deleted\x18\x04 \x01(\bR\tisDeleted\"\"\n" +
	"\x0eShortenRequest\x12\x10\n" +
	"\x03url\x18\x01 \x01(\tR\x03url\"#\n" +
	"\x0fShortenResponse\x12\x10\n" +
	"\x03url\x18\x01 \x01(\tR\x03url\" \n" +
	"\x0eGetItemRequest\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\tR\x02id\"#\n" +
	"\x0fGetItemResponse\x12\x10\n" +
	"\x03url\x18\x01 \x01(\tR\x03url\"(\n" +
	"\fPingResponse\x12\x18\n" +
	"\amessage\x18\x01 \x01(\tR\amessage\"H\n" +
	"\x13ShortenBatchRequest\x121\n" +
	"\x05items\x18\x01 \x03(\v2\x1b.shortener.ShortenBatchItemR\x05items\"\\\n" +
	"\x10ShortenBatchItem\x12%\n" +
	"\x0ecorrelation_id\x18\x01 \x01(\tR\rcorrelationId\x12!\n" +
	"\foriginal_url\x18\x02 \x01(\tR\voriginalUrl\"O\n" +
	"\x14ShortenBatchResponse\x127\n" +
	"\aresults\x18\x01 \x03(\v2\x1d.shortener.ShortenBatchResultR\aresults\"X\n" +
	"\x12ShortenBatchResult\x12%\n" +
	"\x0ecorrelation_id\x18\x01 \x01(\tR\rcorrelationId\x12\x1b\n" +
	"\tshort_url\x18\x02 \x01(\tR\bshortUrl\"G\n" +
	"\x10UserUrlsResponse\x123\n" +
	"\aresults\x18\x01 \x03(\v2\x19.shortener.UserUrlsResultR\aresults\"P\n" +
	"\x0eUserUrlsResult\x12\x1b\n" +
	"\tshort_url\x18\x01 \x01(\tR\bshortUrl\x12!\n" +
	"\foriginal_url\x18\x02 \x01(\tR\voriginalUrl\"L\n" +
	"\x15UserUrlsDeleteRequest\x123\n" +
	"\x05items\x18\x01 \x03(\v2\x1d.shortener.UserUrlsDeleteItemR\x05items\"1\n" +
	"\x12UserUrlsDeleteItem\x12\x1b\n" +
	"\tshort_url\x18\x01 \x01(\tR\bshortUrl\"P\n" +
	"\x18ShortenBatchUpdateRecord\x12\x1b\n" +
	"\tshort_url\x18\x01 \x01(\tR\bshortUrl\x12\x17\n" +
	"\auser_id\x18\x02 \x01(\x03R\x06userId\"2\n" +
	"\x16UserUrlsDeleteResponse\x12\x18\n" +
	"\amessage\x18\x01 \x01(\tR\amessage\"2\n" +
	"\x14TrustedSubnetRequest\x12\x1a\n" +
	"\tx_real_ip\x18\x01 \x01(\tR\axRealIp\"A\n" +
	"\x15TrustedSubnetResponse\x12\x14\n" +
	"\x05users\x18\x01 \x01(\x03R\x05users\x12\x12\n" +
	"\x04urls\x18\x02 \x01(\x03R\x04urls2\xf9\x04\n" +
	"\x10ShortenerService\x12A\n" +
	"\bMainPage\x12\x19.shortener.ShortenRequest\x1a\x1a.shortener.ShortenResponse\x12@\n" +
	"\aGetItem\x12\x19.shortener.GetItemRequest\x1a\x1a.shortener.GetItemResponse\x12G\n" +
	"\x0eShortenHandler\x12\x19.shortener.ShortenRequest\x1a\x1a.shortener.ShortenResponse\x12>\n" +
	"\vPingHandler\x12\x16.google.protobuf.Empty\x1a\x17.shortener.PingResponse\x12V\n" +
	"\x13ShortenBatchHandler\x12\x1e.shortener.ShortenBatchRequest\x1a\x1f.shortener.ShortenBatchResponse\x12F\n" +
	"\x0fUserUrlsHandler\x12\x16.google.protobuf.Empty\x1a\x1b.shortener.UserUrlsResponse\x12\\\n" +
	"\x15UserUrlsDeleteHandler\x12 .shortener.UserUrlsDeleteRequest\x1a!.shortener.UserUrlsDeleteResponse\x12Y\n" +
	"\x14TrustedSubnetHandler\x12\x1f.shortener.TrustedSubnetRequest\x1a .shortener.TrustedSubnetResponseB.Z,yandex-go-advanced/proto/shortener;shortenerb\x06proto3"

var (
	file_shortener_proto_rawDescOnce sync.Once
	file_shortener_proto_rawDescData []byte
)

func file_shortener_proto_rawDescGZIP() []byte {
	file_shortener_proto_rawDescOnce.Do(func() {
		file_shortener_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_shortener_proto_rawDesc), len(file_shortener_proto_rawDesc)))
	})
	return file_shortener_proto_rawDescData
}

var file_shortener_proto_msgTypes = make([]protoimpl.MessageInfo, 18)
var file_shortener_proto_goTypes = []any{
	(*ShortenRecord)(nil),            // 0: shortener.ShortenRecord
	(*ShortenRequest)(nil),           // 1: shortener.ShortenRequest
	(*ShortenResponse)(nil),          // 2: shortener.ShortenResponse
	(*GetItemRequest)(nil),           // 3: shortener.GetItemRequest
	(*GetItemResponse)(nil),          // 4: shortener.GetItemResponse
	(*PingResponse)(nil),             // 5: shortener.PingResponse
	(*ShortenBatchRequest)(nil),      // 6: shortener.ShortenBatchRequest
	(*ShortenBatchItem)(nil),         // 7: shortener.ShortenBatchItem
	(*ShortenBatchResponse)(nil),     // 8: shortener.ShortenBatchResponse
	(*ShortenBatchResult)(nil),       // 9: shortener.ShortenBatchResult
	(*UserUrlsResponse)(nil),         // 10: shortener.UserUrlsResponse
	(*UserUrlsResult)(nil),           // 11: shortener.UserUrlsResult
	(*UserUrlsDeleteRequest)(nil),    // 12: shortener.UserUrlsDeleteRequest
	(*UserUrlsDeleteItem)(nil),       // 13: shortener.UserUrlsDeleteItem
	(*ShortenBatchUpdateRecord)(nil), // 14: shortener.ShortenBatchUpdateRecord
	(*UserUrlsDeleteResponse)(nil),   // 15: shortener.UserUrlsDeleteResponse
	(*TrustedSubnetRequest)(nil),     // 16: shortener.TrustedSubnetRequest
	(*TrustedSubnetResponse)(nil),    // 17: shortener.TrustedSubnetResponse
	(*emptypb.Empty)(nil),            // 18: google.protobuf.Empty
}
var file_shortener_proto_depIdxs = []int32{
	7,  // 0: shortener.ShortenBatchRequest.items:type_name -> shortener.ShortenBatchItem
	9,  // 1: shortener.ShortenBatchResponse.results:type_name -> shortener.ShortenBatchResult
	11, // 2: shortener.UserUrlsResponse.results:type_name -> shortener.UserUrlsResult
	13, // 3: shortener.UserUrlsDeleteRequest.items:type_name -> shortener.UserUrlsDeleteItem
	1,  // 4: shortener.ShortenerService.MainPage:input_type -> shortener.ShortenRequest
	3,  // 5: shortener.ShortenerService.GetItem:input_type -> shortener.GetItemRequest
	1,  // 6: shortener.ShortenerService.ShortenHandler:input_type -> shortener.ShortenRequest
	18, // 7: shortener.ShortenerService.PingHandler:input_type -> google.protobuf.Empty
	6,  // 8: shortener.ShortenerService.ShortenBatchHandler:input_type -> shortener.ShortenBatchRequest
	18, // 9: shortener.ShortenerService.UserUrlsHandler:input_type -> google.protobuf.Empty
	12, // 10: shortener.ShortenerService.UserUrlsDeleteHandler:input_type -> shortener.UserUrlsDeleteRequest
	16, // 11: shortener.ShortenerService.TrustedSubnetHandler:input_type -> shortener.TrustedSubnetRequest
	2,  // 12: shortener.ShortenerService.MainPage:output_type -> shortener.ShortenResponse
	4,  // 13: shortener.ShortenerService.GetItem:output_type -> shortener.GetItemResponse
	2,  // 14: shortener.ShortenerService.ShortenHandler:output_type -> shortener.ShortenResponse
	5,  // 15: shortener.ShortenerService.PingHandler:output_type -> shortener.PingResponse
	8,  // 16: shortener.ShortenerService.ShortenBatchHandler:output_type -> shortener.ShortenBatchResponse
	10, // 17: shortener.ShortenerService.UserUrlsHandler:output_type -> shortener.UserUrlsResponse
	15, // 18: shortener.ShortenerService.UserUrlsDeleteHandler:output_type -> shortener.UserUrlsDeleteResponse
	17, // 19: shortener.ShortenerService.TrustedSubnetHandler:output_type -> shortener.TrustedSubnetResponse
	12, // [12:20] is the sub-list for method output_type
	4,  // [4:12] is the sub-list for method input_type
	4,  // [4:4] is the sub-list for extension type_name
	4,  // [4:4] is the sub-list for extension extendee
	0,  // [0:4] is the sub-list for field type_name
}

func init() { file_shortener_proto_init() }
func file_shortener_proto_init() {
	if File_shortener_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_shortener_proto_rawDesc), len(file_shortener_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   18,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_shortener_proto_goTypes,
		DependencyIndexes: file_shortener_proto_depIdxs,
		MessageInfos:      file_shortener_proto_msgTypes,
	}.Build()
	File_shortener_proto = out.File
	file_shortener_proto_goTypes = nil
	file_shortener_proto_depIdxs = nil
}
