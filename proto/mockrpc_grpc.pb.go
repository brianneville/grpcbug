// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.19.4
// source: mockrpc.proto

package proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// MockRPCClient is the client API for MockRPC service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MockRPCClient interface {
	Get(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*GetResponse, error)
}

type mockRPCClient struct {
	cc grpc.ClientConnInterface
}

func NewMockRPCClient(cc grpc.ClientConnInterface) MockRPCClient {
	return &mockRPCClient{cc}
}

func (c *mockRPCClient) Get(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*GetResponse, error) {
	out := new(GetResponse)
	err := c.cc.Invoke(ctx, "/grpcbug.MockRPC/Get", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MockRPCServer is the server API for MockRPC service.
// All implementations must embed UnimplementedMockRPCServer
// for forward compatibility
type MockRPCServer interface {
	Get(context.Context, *GetRequest) (*GetResponse, error)
	mustEmbedUnimplementedMockRPCServer()
}

// UnimplementedMockRPCServer must be embedded to have forward compatible implementations.
type UnimplementedMockRPCServer struct {
}

func (UnimplementedMockRPCServer) Get(context.Context, *GetRequest) (*GetResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}
func (UnimplementedMockRPCServer) mustEmbedUnimplementedMockRPCServer() {}

// UnsafeMockRPCServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MockRPCServer will
// result in compilation errors.
type UnsafeMockRPCServer interface {
	mustEmbedUnimplementedMockRPCServer()
}

func RegisterMockRPCServer(s grpc.ServiceRegistrar, srv MockRPCServer) {
	s.RegisterService(&MockRPC_ServiceDesc, srv)
}

func _MockRPC_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MockRPCServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpcbug.MockRPC/Get",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MockRPCServer).Get(ctx, req.(*GetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// MockRPC_ServiceDesc is the grpc.ServiceDesc for MockRPC service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var MockRPC_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "grpcbug.MockRPC",
	HandlerType: (*MockRPCServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Get",
			Handler:    _MockRPC_Get_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "mockrpc.proto",
}