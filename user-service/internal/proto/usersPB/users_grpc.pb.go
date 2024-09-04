// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v3.12.4
// source: users.proto

package usersPB

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	UsersService_RegisterUser_FullMethodName   = "/usersPB.UsersService/RegisterUser"
	UsersService_LoginUser_FullMethodName      = "/usersPB.UsersService/LoginUser"
	UsersService_GetById_FullMethodName        = "/usersPB.UsersService/GetById"
	UsersService_UpdateUser_FullMethodName     = "/usersPB.UsersService/UpdateUser"
	UsersService_DeleteUser_FullMethodName     = "/usersPB.UsersService/DeleteUser"
	UsersService_ListUsers_FullMethodName      = "/usersPB.UsersService/ListUsers"
	UsersService_GetProfileById_FullMethodName = "/usersPB.UsersService/GetProfileById"
	UsersService_UpdateProfile_FullMethodName  = "/usersPB.UsersService/UpdateProfile"
	UsersService_UploadAvatar_FullMethodName   = "/usersPB.UsersService/UploadAvatar"
	UsersService_SearchUser_FullMethodName     = "/usersPB.UsersService/SearchUser"
)

// UsersServiceClient is the client API for UsersService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type UsersServiceClient interface {
	RegisterUser(ctx context.Context, in *User, opts ...grpc.CallOption) (*RegisterResponse, error)
	LoginUser(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*LoginResponse, error)
	GetById(ctx context.Context, in *GetByIdRequest, opts ...grpc.CallOption) (*User, error)
	UpdateUser(ctx context.Context, in *UpdateUserRequest, opts ...grpc.CallOption) (*MsgResponse, error)
	DeleteUser(ctx context.Context, in *GetByIdRequest, opts ...grpc.CallOption) (*MsgResponse, error)
	ListUsers(ctx context.Context, in *ListRequest, opts ...grpc.CallOption) (*ListResponse, error)
	GetProfileById(ctx context.Context, in *GetByIdRequest, opts ...grpc.CallOption) (*Profile, error)
	UpdateProfile(ctx context.Context, in *UpdateProfileRequest, opts ...grpc.CallOption) (*MsgResponse, error)
	UploadAvatar(ctx context.Context, in *UploadAvatarRequest, opts ...grpc.CallOption) (*MsgResponse, error)
	SearchUser(ctx context.Context, in *SearchRequest, opts ...grpc.CallOption) (*ListResponse, error)
}

type usersServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewUsersServiceClient(cc grpc.ClientConnInterface) UsersServiceClient {
	return &usersServiceClient{cc}
}

func (c *usersServiceClient) RegisterUser(ctx context.Context, in *User, opts ...grpc.CallOption) (*RegisterResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(RegisterResponse)
	err := c.cc.Invoke(ctx, UsersService_RegisterUser_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *usersServiceClient) LoginUser(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*LoginResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(LoginResponse)
	err := c.cc.Invoke(ctx, UsersService_LoginUser_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *usersServiceClient) GetById(ctx context.Context, in *GetByIdRequest, opts ...grpc.CallOption) (*User, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(User)
	err := c.cc.Invoke(ctx, UsersService_GetById_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *usersServiceClient) UpdateUser(ctx context.Context, in *UpdateUserRequest, opts ...grpc.CallOption) (*MsgResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(MsgResponse)
	err := c.cc.Invoke(ctx, UsersService_UpdateUser_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *usersServiceClient) DeleteUser(ctx context.Context, in *GetByIdRequest, opts ...grpc.CallOption) (*MsgResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(MsgResponse)
	err := c.cc.Invoke(ctx, UsersService_DeleteUser_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *usersServiceClient) ListUsers(ctx context.Context, in *ListRequest, opts ...grpc.CallOption) (*ListResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ListResponse)
	err := c.cc.Invoke(ctx, UsersService_ListUsers_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *usersServiceClient) GetProfileById(ctx context.Context, in *GetByIdRequest, opts ...grpc.CallOption) (*Profile, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Profile)
	err := c.cc.Invoke(ctx, UsersService_GetProfileById_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *usersServiceClient) UpdateProfile(ctx context.Context, in *UpdateProfileRequest, opts ...grpc.CallOption) (*MsgResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(MsgResponse)
	err := c.cc.Invoke(ctx, UsersService_UpdateProfile_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *usersServiceClient) UploadAvatar(ctx context.Context, in *UploadAvatarRequest, opts ...grpc.CallOption) (*MsgResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(MsgResponse)
	err := c.cc.Invoke(ctx, UsersService_UploadAvatar_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *usersServiceClient) SearchUser(ctx context.Context, in *SearchRequest, opts ...grpc.CallOption) (*ListResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ListResponse)
	err := c.cc.Invoke(ctx, UsersService_SearchUser_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// UsersServiceServer is the server API for UsersService service.
// All implementations must embed UnimplementedUsersServiceServer
// for forward compatibility.
type UsersServiceServer interface {
	RegisterUser(context.Context, *User) (*RegisterResponse, error)
	LoginUser(context.Context, *LoginRequest) (*LoginResponse, error)
	GetById(context.Context, *GetByIdRequest) (*User, error)
	UpdateUser(context.Context, *UpdateUserRequest) (*MsgResponse, error)
	DeleteUser(context.Context, *GetByIdRequest) (*MsgResponse, error)
	ListUsers(context.Context, *ListRequest) (*ListResponse, error)
	GetProfileById(context.Context, *GetByIdRequest) (*Profile, error)
	UpdateProfile(context.Context, *UpdateProfileRequest) (*MsgResponse, error)
	UploadAvatar(context.Context, *UploadAvatarRequest) (*MsgResponse, error)
	SearchUser(context.Context, *SearchRequest) (*ListResponse, error)
	mustEmbedUnimplementedUsersServiceServer()
}

// UnimplementedUsersServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedUsersServiceServer struct{}

func (UnimplementedUsersServiceServer) RegisterUser(context.Context, *User) (*RegisterResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RegisterUser not implemented")
}
func (UnimplementedUsersServiceServer) LoginUser(context.Context, *LoginRequest) (*LoginResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LoginUser not implemented")
}
func (UnimplementedUsersServiceServer) GetById(context.Context, *GetByIdRequest) (*User, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetById not implemented")
}
func (UnimplementedUsersServiceServer) UpdateUser(context.Context, *UpdateUserRequest) (*MsgResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateUser not implemented")
}
func (UnimplementedUsersServiceServer) DeleteUser(context.Context, *GetByIdRequest) (*MsgResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteUser not implemented")
}
func (UnimplementedUsersServiceServer) ListUsers(context.Context, *ListRequest) (*ListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListUsers not implemented")
}
func (UnimplementedUsersServiceServer) GetProfileById(context.Context, *GetByIdRequest) (*Profile, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetProfileById not implemented")
}
func (UnimplementedUsersServiceServer) UpdateProfile(context.Context, *UpdateProfileRequest) (*MsgResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateProfile not implemented")
}
func (UnimplementedUsersServiceServer) UploadAvatar(context.Context, *UploadAvatarRequest) (*MsgResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UploadAvatar not implemented")
}
func (UnimplementedUsersServiceServer) SearchUser(context.Context, *SearchRequest) (*ListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SearchUser not implemented")
}
func (UnimplementedUsersServiceServer) mustEmbedUnimplementedUsersServiceServer() {}
func (UnimplementedUsersServiceServer) testEmbeddedByValue()                      {}

// UnsafeUsersServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to UsersServiceServer will
// result in compilation errors.
type UnsafeUsersServiceServer interface {
	mustEmbedUnimplementedUsersServiceServer()
}

func RegisterUsersServiceServer(s grpc.ServiceRegistrar, srv UsersServiceServer) {
	// If the following call pancis, it indicates UnimplementedUsersServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&UsersService_ServiceDesc, srv)
}

func _UsersService_RegisterUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(User)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UsersServiceServer).RegisterUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UsersService_RegisterUser_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UsersServiceServer).RegisterUser(ctx, req.(*User))
	}
	return interceptor(ctx, in, info, handler)
}

func _UsersService_LoginUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoginRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UsersServiceServer).LoginUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UsersService_LoginUser_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UsersServiceServer).LoginUser(ctx, req.(*LoginRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UsersService_GetById_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetByIdRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UsersServiceServer).GetById(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UsersService_GetById_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UsersServiceServer).GetById(ctx, req.(*GetByIdRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UsersService_UpdateUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UsersServiceServer).UpdateUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UsersService_UpdateUser_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UsersServiceServer).UpdateUser(ctx, req.(*UpdateUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UsersService_DeleteUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetByIdRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UsersServiceServer).DeleteUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UsersService_DeleteUser_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UsersServiceServer).DeleteUser(ctx, req.(*GetByIdRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UsersService_ListUsers_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UsersServiceServer).ListUsers(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UsersService_ListUsers_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UsersServiceServer).ListUsers(ctx, req.(*ListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UsersService_GetProfileById_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetByIdRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UsersServiceServer).GetProfileById(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UsersService_GetProfileById_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UsersServiceServer).GetProfileById(ctx, req.(*GetByIdRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UsersService_UpdateProfile_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateProfileRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UsersServiceServer).UpdateProfile(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UsersService_UpdateProfile_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UsersServiceServer).UpdateProfile(ctx, req.(*UpdateProfileRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UsersService_UploadAvatar_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UploadAvatarRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UsersServiceServer).UploadAvatar(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UsersService_UploadAvatar_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UsersServiceServer).UploadAvatar(ctx, req.(*UploadAvatarRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UsersService_SearchUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SearchRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UsersServiceServer).SearchUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UsersService_SearchUser_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UsersServiceServer).SearchUser(ctx, req.(*SearchRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// UsersService_ServiceDesc is the grpc.ServiceDesc for UsersService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var UsersService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "usersPB.UsersService",
	HandlerType: (*UsersServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "RegisterUser",
			Handler:    _UsersService_RegisterUser_Handler,
		},
		{
			MethodName: "LoginUser",
			Handler:    _UsersService_LoginUser_Handler,
		},
		{
			MethodName: "GetById",
			Handler:    _UsersService_GetById_Handler,
		},
		{
			MethodName: "UpdateUser",
			Handler:    _UsersService_UpdateUser_Handler,
		},
		{
			MethodName: "DeleteUser",
			Handler:    _UsersService_DeleteUser_Handler,
		},
		{
			MethodName: "ListUsers",
			Handler:    _UsersService_ListUsers_Handler,
		},
		{
			MethodName: "GetProfileById",
			Handler:    _UsersService_GetProfileById_Handler,
		},
		{
			MethodName: "UpdateProfile",
			Handler:    _UsersService_UpdateProfile_Handler,
		},
		{
			MethodName: "UploadAvatar",
			Handler:    _UsersService_UploadAvatar_Handler,
		},
		{
			MethodName: "SearchUser",
			Handler:    _UsersService_SearchUser_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "users.proto",
}
