// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.21.12
// source: random.proto

package api_proto

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

const (
	RandomService_RandomGPRC_FullMethodName = "/random.RandomService/RandomGPRC"
)

// RandomServiceClient is the client API for RandomService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RandomServiceClient interface {
	RandomGPRC(ctx context.Context, in *RandomRequest, opts ...grpc.CallOption) (*RandomResponse, error)
}

type randomServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewRandomServiceClient(cc grpc.ClientConnInterface) RandomServiceClient {
	return &randomServiceClient{cc}
}

func (c *randomServiceClient) RandomGPRC(ctx context.Context, in *RandomRequest, opts ...grpc.CallOption) (*RandomResponse, error) {
	out := new(RandomResponse)
	err := c.cc.Invoke(ctx, RandomService_RandomGPRC_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RandomServiceServer is the server API for RandomService service.
// All implementations must embed UnimplementedRandomServiceServer
// for forward compatibility
type RandomServiceServer interface {
	RandomGPRC(context.Context, *RandomRequest) (*RandomResponse, error)
	mustEmbedUnimplementedRandomServiceServer()
}

// UnimplementedRandomServiceServer must be embedded to have forward compatible implementations.
type UnimplementedRandomServiceServer struct {
}

func (UnimplementedRandomServiceServer) RandomGPRC(context.Context, *RandomRequest) (*RandomResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RandomGPRC not implemented")
}
func (UnimplementedRandomServiceServer) mustEmbedUnimplementedRandomServiceServer() {}

// UnsafeRandomServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RandomServiceServer will
// result in compilation errors.
type UnsafeRandomServiceServer interface {
	mustEmbedUnimplementedRandomServiceServer()
}

func RegisterRandomServiceServer(s grpc.ServiceRegistrar, srv RandomServiceServer) {
	s.RegisterService(&RandomService_ServiceDesc, srv)
}

func _RandomService_RandomGPRC_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RandomRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RandomServiceServer).RandomGPRC(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RandomService_RandomGPRC_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RandomServiceServer).RandomGPRC(ctx, req.(*RandomRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// RandomService_ServiceDesc is the grpc.ServiceDesc for RandomService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var RandomService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "random.RandomService",
	HandlerType: (*RandomServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "RandomGPRC",
			Handler:    _RandomService_RandomGPRC_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "random.proto",
}
