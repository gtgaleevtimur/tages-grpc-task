// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.12.4
// source: proto/proto.proto

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

// TagesClient is the client API for Tages service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type TagesClient interface {
	Upload(ctx context.Context, opts ...grpc.CallOption) (Tages_UploadClient, error)
	Download(ctx context.Context, in *StringForm, opts ...grpc.CallOption) (Tages_DownloadClient, error)
	List(ctx context.Context, in *ListRequest, opts ...grpc.CallOption) (Tages_ListClient, error)
}

type tagesClient struct {
	cc grpc.ClientConnInterface
}

func NewTagesClient(cc grpc.ClientConnInterface) TagesClient {
	return &tagesClient{cc}
}

func (c *tagesClient) Upload(ctx context.Context, opts ...grpc.CallOption) (Tages_UploadClient, error) {
	stream, err := c.cc.NewStream(ctx, &Tages_ServiceDesc.Streams[0], "/tages.Tages/Upload", opts...)
	if err != nil {
		return nil, err
	}
	x := &tagesUploadClient{stream}
	return x, nil
}

type Tages_UploadClient interface {
	Send(*UploadRequest) error
	CloseAndRecv() (*StringForm, error)
	grpc.ClientStream
}

type tagesUploadClient struct {
	grpc.ClientStream
}

func (x *tagesUploadClient) Send(m *UploadRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *tagesUploadClient) CloseAndRecv() (*StringForm, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(StringForm)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *tagesClient) Download(ctx context.Context, in *StringForm, opts ...grpc.CallOption) (Tages_DownloadClient, error) {
	stream, err := c.cc.NewStream(ctx, &Tages_ServiceDesc.Streams[1], "/tages.Tages/Download", opts...)
	if err != nil {
		return nil, err
	}
	x := &tagesDownloadClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Tages_DownloadClient interface {
	Recv() (*DownloadResponse, error)
	grpc.ClientStream
}

type tagesDownloadClient struct {
	grpc.ClientStream
}

func (x *tagesDownloadClient) Recv() (*DownloadResponse, error) {
	m := new(DownloadResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *tagesClient) List(ctx context.Context, in *ListRequest, opts ...grpc.CallOption) (Tages_ListClient, error) {
	stream, err := c.cc.NewStream(ctx, &Tages_ServiceDesc.Streams[2], "/tages.Tages/List", opts...)
	if err != nil {
		return nil, err
	}
	x := &tagesListClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Tages_ListClient interface {
	Recv() (*StringForm, error)
	grpc.ClientStream
}

type tagesListClient struct {
	grpc.ClientStream
}

func (x *tagesListClient) Recv() (*StringForm, error) {
	m := new(StringForm)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// TagesServer is the server API for Tages service.
// All implementations must embed UnimplementedTagesServer
// for forward compatibility
type TagesServer interface {
	Upload(Tages_UploadServer) error
	Download(*StringForm, Tages_DownloadServer) error
	List(*ListRequest, Tages_ListServer) error
	mustEmbedUnimplementedTagesServer()
}

// UnimplementedTagesServer must be embedded to have forward compatible implementations.
type UnimplementedTagesServer struct {
}

func (UnimplementedTagesServer) Upload(Tages_UploadServer) error {
	return status.Errorf(codes.Unimplemented, "method Upload not implemented")
}
func (UnimplementedTagesServer) Download(*StringForm, Tages_DownloadServer) error {
	return status.Errorf(codes.Unimplemented, "method Download not implemented")
}
func (UnimplementedTagesServer) List(*ListRequest, Tages_ListServer) error {
	return status.Errorf(codes.Unimplemented, "method List not implemented")
}
func (UnimplementedTagesServer) mustEmbedUnimplementedTagesServer() {}

// UnsafeTagesServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to TagesServer will
// result in compilation errors.
type UnsafeTagesServer interface {
	mustEmbedUnimplementedTagesServer()
}

func RegisterTagesServer(s grpc.ServiceRegistrar, srv TagesServer) {
	s.RegisterService(&Tages_ServiceDesc, srv)
}

func _Tages_Upload_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(TagesServer).Upload(&tagesUploadServer{stream})
}

type Tages_UploadServer interface {
	SendAndClose(*StringForm) error
	Recv() (*UploadRequest, error)
	grpc.ServerStream
}

type tagesUploadServer struct {
	grpc.ServerStream
}

func (x *tagesUploadServer) SendAndClose(m *StringForm) error {
	return x.ServerStream.SendMsg(m)
}

func (x *tagesUploadServer) Recv() (*UploadRequest, error) {
	m := new(UploadRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _Tages_Download_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(StringForm)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(TagesServer).Download(m, &tagesDownloadServer{stream})
}

type Tages_DownloadServer interface {
	Send(*DownloadResponse) error
	grpc.ServerStream
}

type tagesDownloadServer struct {
	grpc.ServerStream
}

func (x *tagesDownloadServer) Send(m *DownloadResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _Tages_List_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(ListRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(TagesServer).List(m, &tagesListServer{stream})
}

type Tages_ListServer interface {
	Send(*StringForm) error
	grpc.ServerStream
}

type tagesListServer struct {
	grpc.ServerStream
}

func (x *tagesListServer) Send(m *StringForm) error {
	return x.ServerStream.SendMsg(m)
}

// Tages_ServiceDesc is the grpc.ServiceDesc for Tages service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Tages_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "tages.Tages",
	HandlerType: (*TagesServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Upload",
			Handler:       _Tages_Upload_Handler,
			ClientStreams: true,
		},
		{
			StreamName:    "Download",
			Handler:       _Tages_Download_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "List",
			Handler:       _Tages_List_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "proto/proto.proto",
}
