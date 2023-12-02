// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.21.12
// source: booking.proto

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
	BookingService_LetsBook_FullMethodName   = "/pb.BookingService/LetsBook"
	BookingService_DontBook_FullMethodName   = "/pb.BookingService/DontBook"
	BookingService_BookingsOf_FullMethodName = "/pb.BookingService/BookingsOf"
)

// BookingServiceClient is the client API for BookingService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type BookingServiceClient interface {
	LetsBook(ctx context.Context, in *BookingReq, opts ...grpc.CallOption) (*BookingsRes, error)
	DontBook(ctx context.Context, in *DontReq, opts ...grpc.CallOption) (*DontRes, error)
	BookingsOf(ctx context.Context, in *BookingsOfReq, opts ...grpc.CallOption) (*BookingsRes, error)
}

type bookingServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewBookingServiceClient(cc grpc.ClientConnInterface) BookingServiceClient {
	return &bookingServiceClient{cc}
}

func (c *bookingServiceClient) LetsBook(ctx context.Context, in *BookingReq, opts ...grpc.CallOption) (*BookingsRes, error) {
	out := new(BookingsRes)
	err := c.cc.Invoke(ctx, BookingService_LetsBook_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *bookingServiceClient) DontBook(ctx context.Context, in *DontReq, opts ...grpc.CallOption) (*DontRes, error) {
	out := new(DontRes)
	err := c.cc.Invoke(ctx, BookingService_DontBook_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *bookingServiceClient) BookingsOf(ctx context.Context, in *BookingsOfReq, opts ...grpc.CallOption) (*BookingsRes, error) {
	out := new(BookingsRes)
	err := c.cc.Invoke(ctx, BookingService_BookingsOf_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// BookingServiceServer is the server API for BookingService service.
// All implementations must embed UnimplementedBookingServiceServer
// for forward compatibility
type BookingServiceServer interface {
	LetsBook(context.Context, *BookingReq) (*BookingsRes, error)
	DontBook(context.Context, *DontReq) (*DontRes, error)
	BookingsOf(context.Context, *BookingsOfReq) (*BookingsRes, error)
	mustEmbedUnimplementedBookingServiceServer()
}

// UnimplementedBookingServiceServer must be embedded to have forward compatible implementations.
type UnimplementedBookingServiceServer struct {
}

func (UnimplementedBookingServiceServer) LetsBook(context.Context, *BookingReq) (*BookingsRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LetsBook not implemented")
}
func (UnimplementedBookingServiceServer) DontBook(context.Context, *DontReq) (*DontRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DontBook not implemented")
}
func (UnimplementedBookingServiceServer) BookingsOf(context.Context, *BookingsOfReq) (*BookingsRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method BookingsOf not implemented")
}
func (UnimplementedBookingServiceServer) mustEmbedUnimplementedBookingServiceServer() {}

// UnsafeBookingServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to BookingServiceServer will
// result in compilation errors.
type UnsafeBookingServiceServer interface {
	mustEmbedUnimplementedBookingServiceServer()
}

func RegisterBookingServiceServer(s grpc.ServiceRegistrar, srv BookingServiceServer) {
	s.RegisterService(&BookingService_ServiceDesc, srv)
}

func _BookingService_LetsBook_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BookingReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BookingServiceServer).LetsBook(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: BookingService_LetsBook_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BookingServiceServer).LetsBook(ctx, req.(*BookingReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _BookingService_DontBook_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DontReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BookingServiceServer).DontBook(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: BookingService_DontBook_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BookingServiceServer).DontBook(ctx, req.(*DontReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _BookingService_BookingsOf_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BookingsOfReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BookingServiceServer).BookingsOf(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: BookingService_BookingsOf_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BookingServiceServer).BookingsOf(ctx, req.(*BookingsOfReq))
	}
	return interceptor(ctx, in, info, handler)
}

// BookingService_ServiceDesc is the grpc.ServiceDesc for BookingService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var BookingService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "pb.BookingService",
	HandlerType: (*BookingServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "LetsBook",
			Handler:    _BookingService_LetsBook_Handler,
		},
		{
			MethodName: "DontBook",
			Handler:    _BookingService_DontBook_Handler,
		},
		{
			MethodName: "BookingsOf",
			Handler:    _BookingService_BookingsOf_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "booking.proto",
}
