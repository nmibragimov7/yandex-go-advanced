// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.29.3
// source: shortener.proto

package shortener

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	ShortenerService_MainPage_FullMethodName              = "/shortener.ShortenerService/MainPage"
	ShortenerService_GetItem_FullMethodName               = "/shortener.ShortenerService/GetItem"
	ShortenerService_ShortenHandler_FullMethodName        = "/shortener.ShortenerService/ShortenHandler"
	ShortenerService_PingHandler_FullMethodName           = "/shortener.ShortenerService/PingHandler"
	ShortenerService_ShortenBatchHandler_FullMethodName   = "/shortener.ShortenerService/ShortenBatchHandler"
	ShortenerService_UserUrlsHandler_FullMethodName       = "/shortener.ShortenerService/UserUrlsHandler"
	ShortenerService_UserUrlsDeleteHandler_FullMethodName = "/shortener.ShortenerService/UserUrlsDeleteHandler"
	ShortenerService_TrustedSubnetHandler_FullMethodName  = "/shortener.ShortenerService/TrustedSubnetHandler"
)

// ShortenerServiceClient is the client API for ShortenerService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ShortenerServiceClient interface {
	MainPage(ctx context.Context, in *ShortenRequest, opts ...grpc.CallOption) (*ShortenResponse, error)
	GetItem(ctx context.Context, in *GetItemRequest, opts ...grpc.CallOption) (*GetItemResponse, error)
	ShortenHandler(ctx context.Context, in *ShortenRequest, opts ...grpc.CallOption) (*ShortenResponse, error)
	PingHandler(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*PingResponse, error)
	ShortenBatchHandler(ctx context.Context, in *ShortenBatchRequest, opts ...grpc.CallOption) (*ShortenBatchResponse, error)
	UserUrlsHandler(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*UserUrlsResponse, error)
	UserUrlsDeleteHandler(ctx context.Context, in *UserUrlsDeleteRequest, opts ...grpc.CallOption) (*UserUrlsDeleteResponse, error)
	TrustedSubnetHandler(ctx context.Context, in *TrustedSubnetRequest, opts ...grpc.CallOption) (*TrustedSubnetResponse, error)
}

type shortenerServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewShortenerServiceClient(cc grpc.ClientConnInterface) ShortenerServiceClient {
	return &shortenerServiceClient{cc}
}

func (c *shortenerServiceClient) MainPage(ctx context.Context, in *ShortenRequest, opts ...grpc.CallOption) (*ShortenResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ShortenResponse)
	err := c.cc.Invoke(ctx, ShortenerService_MainPage_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *shortenerServiceClient) GetItem(ctx context.Context, in *GetItemRequest, opts ...grpc.CallOption) (*GetItemResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetItemResponse)
	err := c.cc.Invoke(ctx, ShortenerService_GetItem_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *shortenerServiceClient) ShortenHandler(ctx context.Context, in *ShortenRequest, opts ...grpc.CallOption) (*ShortenResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ShortenResponse)
	err := c.cc.Invoke(ctx, ShortenerService_ShortenHandler_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *shortenerServiceClient) PingHandler(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*PingResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(PingResponse)
	err := c.cc.Invoke(ctx, ShortenerService_PingHandler_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *shortenerServiceClient) ShortenBatchHandler(ctx context.Context, in *ShortenBatchRequest, opts ...grpc.CallOption) (*ShortenBatchResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ShortenBatchResponse)
	err := c.cc.Invoke(ctx, ShortenerService_ShortenBatchHandler_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *shortenerServiceClient) UserUrlsHandler(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*UserUrlsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UserUrlsResponse)
	err := c.cc.Invoke(ctx, ShortenerService_UserUrlsHandler_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *shortenerServiceClient) UserUrlsDeleteHandler(ctx context.Context, in *UserUrlsDeleteRequest, opts ...grpc.CallOption) (*UserUrlsDeleteResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UserUrlsDeleteResponse)
	err := c.cc.Invoke(ctx, ShortenerService_UserUrlsDeleteHandler_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *shortenerServiceClient) TrustedSubnetHandler(ctx context.Context, in *TrustedSubnetRequest, opts ...grpc.CallOption) (*TrustedSubnetResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(TrustedSubnetResponse)
	err := c.cc.Invoke(ctx, ShortenerService_TrustedSubnetHandler_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ShortenerServiceServer is the server API for ShortenerService service.
// All implementations must embed UnimplementedShortenerServiceServer
// for forward compatibility.
type ShortenerServiceServer interface {
	MainPage(context.Context, *ShortenRequest) (*ShortenResponse, error)
	GetItem(context.Context, *GetItemRequest) (*GetItemResponse, error)
	ShortenHandler(context.Context, *ShortenRequest) (*ShortenResponse, error)
	PingHandler(context.Context, *emptypb.Empty) (*PingResponse, error)
	ShortenBatchHandler(context.Context, *ShortenBatchRequest) (*ShortenBatchResponse, error)
	UserUrlsHandler(context.Context, *emptypb.Empty) (*UserUrlsResponse, error)
	UserUrlsDeleteHandler(context.Context, *UserUrlsDeleteRequest) (*UserUrlsDeleteResponse, error)
	TrustedSubnetHandler(context.Context, *TrustedSubnetRequest) (*TrustedSubnetResponse, error)
	mustEmbedUnimplementedShortenerServiceServer()
}

// UnimplementedShortenerServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedShortenerServiceServer struct{}

func (UnimplementedShortenerServiceServer) MainPage(context.Context, *ShortenRequest) (*ShortenResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MainPage not implemented")
}
func (UnimplementedShortenerServiceServer) GetItem(context.Context, *GetItemRequest) (*GetItemResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetItem not implemented")
}
func (UnimplementedShortenerServiceServer) ShortenHandler(context.Context, *ShortenRequest) (*ShortenResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ShortenHandler not implemented")
}
func (UnimplementedShortenerServiceServer) PingHandler(context.Context, *emptypb.Empty) (*PingResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PingHandler not implemented")
}
func (UnimplementedShortenerServiceServer) ShortenBatchHandler(context.Context, *ShortenBatchRequest) (*ShortenBatchResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ShortenBatchHandler not implemented")
}
func (UnimplementedShortenerServiceServer) UserUrlsHandler(context.Context, *emptypb.Empty) (*UserUrlsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UserUrlsHandler not implemented")
}
func (UnimplementedShortenerServiceServer) UserUrlsDeleteHandler(context.Context, *UserUrlsDeleteRequest) (*UserUrlsDeleteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UserUrlsDeleteHandler not implemented")
}
func (UnimplementedShortenerServiceServer) TrustedSubnetHandler(context.Context, *TrustedSubnetRequest) (*TrustedSubnetResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TrustedSubnetHandler not implemented")
}
func (UnimplementedShortenerServiceServer) mustEmbedUnimplementedShortenerServiceServer() {}
func (UnimplementedShortenerServiceServer) testEmbeddedByValue()                          {}

// UnsafeShortenerServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ShortenerServiceServer will
// result in compilation errors.
type UnsafeShortenerServiceServer interface {
	mustEmbedUnimplementedShortenerServiceServer()
}

func RegisterShortenerServiceServer(s grpc.ServiceRegistrar, srv ShortenerServiceServer) {
	// If the following call pancis, it indicates UnimplementedShortenerServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&ShortenerService_ServiceDesc, srv)
}

func _ShortenerService_MainPage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ShortenRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ShortenerServiceServer).MainPage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ShortenerService_MainPage_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ShortenerServiceServer).MainPage(ctx, req.(*ShortenRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ShortenerService_GetItem_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetItemRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ShortenerServiceServer).GetItem(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ShortenerService_GetItem_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ShortenerServiceServer).GetItem(ctx, req.(*GetItemRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ShortenerService_ShortenHandler_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ShortenRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ShortenerServiceServer).ShortenHandler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ShortenerService_ShortenHandler_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ShortenerServiceServer).ShortenHandler(ctx, req.(*ShortenRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ShortenerService_PingHandler_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ShortenerServiceServer).PingHandler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ShortenerService_PingHandler_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ShortenerServiceServer).PingHandler(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _ShortenerService_ShortenBatchHandler_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ShortenBatchRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ShortenerServiceServer).ShortenBatchHandler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ShortenerService_ShortenBatchHandler_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ShortenerServiceServer).ShortenBatchHandler(ctx, req.(*ShortenBatchRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ShortenerService_UserUrlsHandler_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ShortenerServiceServer).UserUrlsHandler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ShortenerService_UserUrlsHandler_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ShortenerServiceServer).UserUrlsHandler(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _ShortenerService_UserUrlsDeleteHandler_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserUrlsDeleteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ShortenerServiceServer).UserUrlsDeleteHandler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ShortenerService_UserUrlsDeleteHandler_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ShortenerServiceServer).UserUrlsDeleteHandler(ctx, req.(*UserUrlsDeleteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ShortenerService_TrustedSubnetHandler_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TrustedSubnetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ShortenerServiceServer).TrustedSubnetHandler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ShortenerService_TrustedSubnetHandler_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ShortenerServiceServer).TrustedSubnetHandler(ctx, req.(*TrustedSubnetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ShortenerService_ServiceDesc is the grpc.ServiceDesc for ShortenerService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ShortenerService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "shortener.ShortenerService",
	HandlerType: (*ShortenerServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "MainPage",
			Handler:    _ShortenerService_MainPage_Handler,
		},
		{
			MethodName: "GetItem",
			Handler:    _ShortenerService_GetItem_Handler,
		},
		{
			MethodName: "ShortenHandler",
			Handler:    _ShortenerService_ShortenHandler_Handler,
		},
		{
			MethodName: "PingHandler",
			Handler:    _ShortenerService_PingHandler_Handler,
		},
		{
			MethodName: "ShortenBatchHandler",
			Handler:    _ShortenerService_ShortenBatchHandler_Handler,
		},
		{
			MethodName: "UserUrlsHandler",
			Handler:    _ShortenerService_UserUrlsHandler_Handler,
		},
		{
			MethodName: "UserUrlsDeleteHandler",
			Handler:    _ShortenerService_UserUrlsDeleteHandler_Handler,
		},
		{
			MethodName: "TrustedSubnetHandler",
			Handler:    _ShortenerService_TrustedSubnetHandler_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "shortener.proto",
}
