// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             (unknown)
// source: cowboy.proto

package pb

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
	Cowboy_GetShot_FullMethodName = "/Cowboy/GetShot"
)

// CowboyClient is the client API for Cowboy service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CowboyClient interface {
	GetShot(ctx context.Context, in *GetShotRequest, opts ...grpc.CallOption) (*GetShotResponse, error)
}

type cowboyClient struct {
	cc grpc.ClientConnInterface
}

func NewCowboyClient(cc grpc.ClientConnInterface) CowboyClient {
	return &cowboyClient{cc}
}

func (c *cowboyClient) GetShot(ctx context.Context, in *GetShotRequest, opts ...grpc.CallOption) (*GetShotResponse, error) {
	out := new(GetShotResponse)
	err := c.cc.Invoke(ctx, Cowboy_GetShot_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CowboyServer is the server API for Cowboy service.
// All implementations must embed UnimplementedCowboyServer
// for forward compatibility
type CowboyServer interface {
	GetShot(context.Context, *GetShotRequest) (*GetShotResponse, error)
	mustEmbedUnimplementedCowboyServer()
}

// UnimplementedCowboyServer must be embedded to have forward compatible implementations.
type UnimplementedCowboyServer struct {
}

func (UnimplementedCowboyServer) GetShot(context.Context, *GetShotRequest) (*GetShotResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetShot not implemented")
}
func (UnimplementedCowboyServer) mustEmbedUnimplementedCowboyServer() {}

// UnsafeCowboyServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CowboyServer will
// result in compilation errors.
type UnsafeCowboyServer interface {
	mustEmbedUnimplementedCowboyServer()
}

func RegisterCowboyServer(s grpc.ServiceRegistrar, srv CowboyServer) {
	s.RegisterService(&Cowboy_ServiceDesc, srv)
}

func _Cowboy_GetShot_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetShotRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CowboyServer).GetShot(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Cowboy_GetShot_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CowboyServer).GetShot(ctx, req.(*GetShotRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Cowboy_ServiceDesc is the grpc.ServiceDesc for Cowboy service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Cowboy_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "Cowboy",
	HandlerType: (*CowboyServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetShot",
			Handler:    _Cowboy_GetShot_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "cowboy.proto",
}
