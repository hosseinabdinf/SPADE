// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v5.26.1
// source: spade.proto

package __

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

// CuratorClient is the client API for Curator service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CuratorClient interface {
	GetPublicParams(ctx context.Context, in *PublicParamsReq, opts ...grpc.CallOption) (*PublicParamsRes, error)
}

type curatorClient struct {
	cc grpc.ClientConnInterface
}

func NewCuratorClient(cc grpc.ClientConnInterface) CuratorClient {
	return &curatorClient{cc}
}

func (c *curatorClient) GetPublicParams(ctx context.Context, in *PublicParamsReq, opts ...grpc.CallOption) (*PublicParamsRes, error) {
	out := new(PublicParamsRes)
	err := c.cc.Invoke(ctx, "/spadeproto.Curator/GetPublicParams", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CuratorServer is the server API for Curator service.
// All implementations must embed UnimplementedCuratorServer
// for forward compatibility
type CuratorServer interface {
	GetPublicParams(context.Context, *PublicParamsReq) (*PublicParamsRes, error)
	mustEmbedUnimplementedCuratorServer()
}

// UnimplementedCuratorServer must be embedded to have forward compatible implementations.
type UnimplementedCuratorServer struct {
}

func (UnimplementedCuratorServer) GetPublicParams(context.Context, *PublicParamsReq) (*PublicParamsRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPublicParams not implemented")
}
func (UnimplementedCuratorServer) mustEmbedUnimplementedCuratorServer() {}

// UnsafeCuratorServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CuratorServer will
// result in compilation errors.
type UnsafeCuratorServer interface {
	mustEmbedUnimplementedCuratorServer()
}

func RegisterCuratorServer(s grpc.ServiceRegistrar, srv CuratorServer) {
	s.RegisterService(&Curator_ServiceDesc, srv)
}

func _Curator_GetPublicParams_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PublicParamsReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CuratorServer).GetPublicParams(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/spadeproto.Curator/GetPublicParams",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CuratorServer).GetPublicParams(ctx, req.(*PublicParamsReq))
	}
	return interceptor(ctx, in, info, handler)
}

// Curator_ServiceDesc is the grpc.ServiceDesc for Curator service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Curator_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "spadeproto.Curator",
	HandlerType: (*CuratorServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetPublicParams",
			Handler:    _Curator_GetPublicParams_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "spade.proto",
}
