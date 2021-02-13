// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package grpcapi

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion7

// BARTClient is the client API for BART service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type BARTClient interface {
	// Sends a request to classify.
	Classify(ctx context.Context, in *ClassifyRequest, opts ...grpc.CallOption) (*ClassifyReply, error)
	// Sends a request to classify-nli.
	ClassifyNLI(ctx context.Context, in *ClassifyNLIRequest, opts ...grpc.CallOption) (*ClassifyReply, error)
	// Send a request to generate.
	Generate(ctx context.Context, in *GenerateRequest, opts ...grpc.CallOption) (*GenerateReply, error)
}

type bARTClient struct {
	cc grpc.ClientConnInterface
}

func NewBARTClient(cc grpc.ClientConnInterface) BARTClient {
	return &bARTClient{cc}
}

func (c *bARTClient) Classify(ctx context.Context, in *ClassifyRequest, opts ...grpc.CallOption) (*ClassifyReply, error) {
	out := new(ClassifyReply)
	err := c.cc.Invoke(ctx, "/bart.grpcapi.BART/Classify", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *bARTClient) ClassifyNLI(ctx context.Context, in *ClassifyNLIRequest, opts ...grpc.CallOption) (*ClassifyReply, error) {
	out := new(ClassifyReply)
	err := c.cc.Invoke(ctx, "/bart.grpcapi.BART/ClassifyNLI", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *bARTClient) Generate(ctx context.Context, in *GenerateRequest, opts ...grpc.CallOption) (*GenerateReply, error) {
	out := new(GenerateReply)
	err := c.cc.Invoke(ctx, "/bart.grpcapi.BART/Generate", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// BARTServer is the server API for BART service.
// All implementations must embed UnimplementedBARTServer
// for forward compatibility
type BARTServer interface {
	// Sends a request to classify.
	Classify(context.Context, *ClassifyRequest) (*ClassifyReply, error)
	// Sends a request to classify-nli.
	ClassifyNLI(context.Context, *ClassifyNLIRequest) (*ClassifyReply, error)
	// Send a request to generate.
	Generate(context.Context, *GenerateRequest) (*GenerateReply, error)
	mustEmbedUnimplementedBARTServer()
}

// UnimplementedBARTServer must be embedded to have forward compatible implementations.
type UnimplementedBARTServer struct {
}

func (UnimplementedBARTServer) Classify(context.Context, *ClassifyRequest) (*ClassifyReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Classify not implemented")
}
func (UnimplementedBARTServer) ClassifyNLI(context.Context, *ClassifyNLIRequest) (*ClassifyReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ClassifyNLI not implemented")
}
func (UnimplementedBARTServer) Generate(context.Context, *GenerateRequest) (*GenerateReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Generate not implemented")
}
func (UnimplementedBARTServer) mustEmbedUnimplementedBARTServer() {}

// UnsafeBARTServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to BARTServer will
// result in compilation errors.
type UnsafeBARTServer interface {
	mustEmbedUnimplementedBARTServer()
}

func RegisterBARTServer(s grpc.ServiceRegistrar, srv BARTServer) {
	s.RegisterService(&_BART_serviceDesc, srv)
}

func _BART_Classify_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ClassifyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BARTServer).Classify(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/bart.grpcapi.BART/Classify",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BARTServer).Classify(ctx, req.(*ClassifyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _BART_ClassifyNLI_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ClassifyNLIRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BARTServer).ClassifyNLI(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/bart.grpcapi.BART/ClassifyNLI",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BARTServer).ClassifyNLI(ctx, req.(*ClassifyNLIRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _BART_Generate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GenerateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BARTServer).Generate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/bart.grpcapi.BART/Generate",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BARTServer).Generate(ctx, req.(*GenerateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _BART_serviceDesc = grpc.ServiceDesc{
	ServiceName: "bart.grpcapi.BART",
	HandlerType: (*BARTServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Classify",
			Handler:    _BART_Classify_Handler,
		},
		{
			MethodName: "ClassifyNLI",
			Handler:    _BART_ClassifyNLI_Handler,
		},
		{
			MethodName: "Generate",
			Handler:    _BART_Generate_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "bart.proto",
}
