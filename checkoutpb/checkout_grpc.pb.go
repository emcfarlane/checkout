// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package checkoutpb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// CheckoutClient is the client API for Checkout service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CheckoutClient interface {
	// Authorize creates a new authorization.
	Authorize(ctx context.Context, in *AuthorizeRequest, opts ...grpc.CallOption) (*Authorization, error)
	// Captures funds for an authorization.
	Capture(ctx context.Context, in *CaptureRequest, opts ...grpc.CallOption) (*Authorization, error)
	// Refunds funds for an authorization.
	Refund(ctx context.Context, in *RefundRequest, opts ...grpc.CallOption) (*Authorization, error)
	// Voids an authorization.
	Void(ctx context.Context, in *VoidRequest, opts ...grpc.CallOption) (*Authorization, error)
}

type checkoutClient struct {
	cc grpc.ClientConnInterface
}

func NewCheckoutClient(cc grpc.ClientConnInterface) CheckoutClient {
	return &checkoutClient{cc}
}

func (c *checkoutClient) Authorize(ctx context.Context, in *AuthorizeRequest, opts ...grpc.CallOption) (*Authorization, error) {
	out := new(Authorization)
	err := c.cc.Invoke(ctx, "/checkout.api.Checkout/Authorize", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *checkoutClient) Capture(ctx context.Context, in *CaptureRequest, opts ...grpc.CallOption) (*Authorization, error) {
	out := new(Authorization)
	err := c.cc.Invoke(ctx, "/checkout.api.Checkout/Capture", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *checkoutClient) Refund(ctx context.Context, in *RefundRequest, opts ...grpc.CallOption) (*Authorization, error) {
	out := new(Authorization)
	err := c.cc.Invoke(ctx, "/checkout.api.Checkout/Refund", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *checkoutClient) Void(ctx context.Context, in *VoidRequest, opts ...grpc.CallOption) (*Authorization, error) {
	out := new(Authorization)
	err := c.cc.Invoke(ctx, "/checkout.api.Checkout/Void", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CheckoutServer is the server API for Checkout service.
// All implementations must embed UnimplementedCheckoutServer
// for forward compatibility
type CheckoutServer interface {
	// Authorize creates a new authorization.
	Authorize(context.Context, *AuthorizeRequest) (*Authorization, error)
	// Captures funds for an authorization.
	Capture(context.Context, *CaptureRequest) (*Authorization, error)
	// Refunds funds for an authorization.
	Refund(context.Context, *RefundRequest) (*Authorization, error)
	// Voids an authorization.
	Void(context.Context, *VoidRequest) (*Authorization, error)
	mustEmbedUnimplementedCheckoutServer()
}

// UnimplementedCheckoutServer must be embedded to have forward compatible implementations.
type UnimplementedCheckoutServer struct {
}

func (*UnimplementedCheckoutServer) Authorize(context.Context, *AuthorizeRequest) (*Authorization, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Authorize not implemented")
}
func (*UnimplementedCheckoutServer) Capture(context.Context, *CaptureRequest) (*Authorization, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Capture not implemented")
}
func (*UnimplementedCheckoutServer) Refund(context.Context, *RefundRequest) (*Authorization, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Refund not implemented")
}
func (*UnimplementedCheckoutServer) Void(context.Context, *VoidRequest) (*Authorization, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Void not implemented")
}
func (*UnimplementedCheckoutServer) mustEmbedUnimplementedCheckoutServer() {}

func RegisterCheckoutServer(s *grpc.Server, srv CheckoutServer) {
	s.RegisterService(&_Checkout_serviceDesc, srv)
}

func _Checkout_Authorize_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AuthorizeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CheckoutServer).Authorize(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/checkout.api.Checkout/Authorize",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CheckoutServer).Authorize(ctx, req.(*AuthorizeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Checkout_Capture_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CaptureRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CheckoutServer).Capture(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/checkout.api.Checkout/Capture",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CheckoutServer).Capture(ctx, req.(*CaptureRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Checkout_Refund_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RefundRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CheckoutServer).Refund(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/checkout.api.Checkout/Refund",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CheckoutServer).Refund(ctx, req.(*RefundRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Checkout_Void_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(VoidRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CheckoutServer).Void(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/checkout.api.Checkout/Void",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CheckoutServer).Void(ctx, req.(*VoidRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Checkout_serviceDesc = grpc.ServiceDesc{
	ServiceName: "checkout.api.Checkout",
	HandlerType: (*CheckoutServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Authorize",
			Handler:    _Checkout_Authorize_Handler,
		},
		{
			MethodName: "Capture",
			Handler:    _Checkout_Capture_Handler,
		},
		{
			MethodName: "Refund",
			Handler:    _Checkout_Refund_Handler,
		},
		{
			MethodName: "Void",
			Handler:    _Checkout_Void_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "checkout.proto",
}
