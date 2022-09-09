// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package languagev1

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// PolicyServiceClient is the client API for PolicyService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PolicyServiceClient interface {
	AllPolicies(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*AllPoliciesResponse, error)
}

type policyServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewPolicyServiceClient(cc grpc.ClientConnInterface) PolicyServiceClient {
	return &policyServiceClient{cc}
}

func (c *policyServiceClient) AllPolicies(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*AllPoliciesResponse, error) {
	out := new(AllPoliciesResponse)
	err := c.cc.Invoke(ctx, "/aperture.policy.language.v1.PolicyService/AllPolicies", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PolicyServiceServer is the server API for PolicyService service.
// All implementations should embed UnimplementedPolicyServiceServer
// for forward compatibility
type PolicyServiceServer interface {
	AllPolicies(context.Context, *emptypb.Empty) (*AllPoliciesResponse, error)
}

// UnimplementedPolicyServiceServer should be embedded to have forward compatible implementations.
type UnimplementedPolicyServiceServer struct {
}

func (UnimplementedPolicyServiceServer) AllPolicies(context.Context, *emptypb.Empty) (*AllPoliciesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AllPolicies not implemented")
}

// UnsafePolicyServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PolicyServiceServer will
// result in compilation errors.
type UnsafePolicyServiceServer interface {
	mustEmbedUnimplementedPolicyServiceServer()
}

func RegisterPolicyServiceServer(s grpc.ServiceRegistrar, srv PolicyServiceServer) {
	s.RegisterService(&PolicyService_ServiceDesc, srv)
}

func _PolicyService_AllPolicies_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PolicyServiceServer).AllPolicies(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/aperture.policy.language.v1.PolicyService/AllPolicies",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PolicyServiceServer).AllPolicies(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

// PolicyService_ServiceDesc is the grpc.ServiceDesc for PolicyService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var PolicyService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "aperture.policy.language.v1.PolicyService",
	HandlerType: (*PolicyServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AllPolicies",
			Handler:    _PolicyService_AllPolicies_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "aperture/policy/language/v1/policy.proto",
}