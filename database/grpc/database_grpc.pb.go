// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.12
// source: database.proto

package grpc

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

// DatabaseHelperClient is the client API for DatabaseHelper service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type DatabaseHelperClient interface {
	AddRow(ctx context.Context, in *AddRowRequest, opts ...grpc.CallOption) (*Response, error)
	UpdateValue(ctx context.Context, in *UpdateValueRequest, opts ...grpc.CallOption) (*Response, error)
	DeleteRows(ctx context.Context, in *DeleteRowsRequest, opts ...grpc.CallOption) (*Response, error)
	Query(ctx context.Context, in *QueryRequest, opts ...grpc.CallOption) (*Response, error)
	ResetDB(ctx context.Context, in *EmptyRequest, opts ...grpc.CallOption) (*Response, error)
}

type databaseHelperClient struct {
	cc grpc.ClientConnInterface
}

func NewDatabaseHelperClient(cc grpc.ClientConnInterface) DatabaseHelperClient {
	return &databaseHelperClient{cc}
}

func (c *databaseHelperClient) AddRow(ctx context.Context, in *AddRowRequest, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/database.DatabaseHelper/AddRow", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *databaseHelperClient) UpdateValue(ctx context.Context, in *UpdateValueRequest, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/database.DatabaseHelper/UpdateValue", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *databaseHelperClient) DeleteRows(ctx context.Context, in *DeleteRowsRequest, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/database.DatabaseHelper/DeleteRows", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *databaseHelperClient) Query(ctx context.Context, in *QueryRequest, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/database.DatabaseHelper/Query", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *databaseHelperClient) ResetDB(ctx context.Context, in *EmptyRequest, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/database.DatabaseHelper/ResetDB", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DatabaseHelperServer is the server API for DatabaseHelper service.
// All implementations must embed UnimplementedDatabaseHelperServer
// for forward compatibility
type DatabaseHelperServer interface {
	AddRow(context.Context, *AddRowRequest) (*Response, error)
	UpdateValue(context.Context, *UpdateValueRequest) (*Response, error)
	DeleteRows(context.Context, *DeleteRowsRequest) (*Response, error)
	Query(context.Context, *QueryRequest) (*Response, error)
	ResetDB(context.Context, *EmptyRequest) (*Response, error)
	mustEmbedUnimplementedDatabaseHelperServer()
}

// UnimplementedDatabaseHelperServer must be embedded to have forward compatible implementations.
type UnimplementedDatabaseHelperServer struct {
}

func (UnimplementedDatabaseHelperServer) AddRow(context.Context, *AddRowRequest) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddRow not implemented")
}
func (UnimplementedDatabaseHelperServer) UpdateValue(context.Context, *UpdateValueRequest) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateValue not implemented")
}
func (UnimplementedDatabaseHelperServer) DeleteRows(context.Context, *DeleteRowsRequest) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteRows not implemented")
}
func (UnimplementedDatabaseHelperServer) Query(context.Context, *QueryRequest) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Query not implemented")
}
func (UnimplementedDatabaseHelperServer) ResetDB(context.Context, *EmptyRequest) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ResetDB not implemented")
}
func (UnimplementedDatabaseHelperServer) mustEmbedUnimplementedDatabaseHelperServer() {}

// UnsafeDatabaseHelperServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to DatabaseHelperServer will
// result in compilation errors.
type UnsafeDatabaseHelperServer interface {
	mustEmbedUnimplementedDatabaseHelperServer()
}

func RegisterDatabaseHelperServer(s grpc.ServiceRegistrar, srv DatabaseHelperServer) {
	s.RegisterService(&DatabaseHelper_ServiceDesc, srv)
}

func _DatabaseHelper_AddRow_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddRowRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DatabaseHelperServer).AddRow(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/database.DatabaseHelper/AddRow",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DatabaseHelperServer).AddRow(ctx, req.(*AddRowRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DatabaseHelper_UpdateValue_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateValueRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DatabaseHelperServer).UpdateValue(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/database.DatabaseHelper/UpdateValue",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DatabaseHelperServer).UpdateValue(ctx, req.(*UpdateValueRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DatabaseHelper_DeleteRows_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteRowsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DatabaseHelperServer).DeleteRows(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/database.DatabaseHelper/DeleteRows",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DatabaseHelperServer).DeleteRows(ctx, req.(*DeleteRowsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DatabaseHelper_Query_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DatabaseHelperServer).Query(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/database.DatabaseHelper/Query",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DatabaseHelperServer).Query(ctx, req.(*QueryRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DatabaseHelper_ResetDB_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EmptyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DatabaseHelperServer).ResetDB(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/database.DatabaseHelper/ResetDB",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DatabaseHelperServer).ResetDB(ctx, req.(*EmptyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// DatabaseHelper_ServiceDesc is the grpc.ServiceDesc for DatabaseHelper service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var DatabaseHelper_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "database.DatabaseHelper",
	HandlerType: (*DatabaseHelperServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AddRow",
			Handler:    _DatabaseHelper_AddRow_Handler,
		},
		{
			MethodName: "UpdateValue",
			Handler:    _DatabaseHelper_UpdateValue_Handler,
		},
		{
			MethodName: "DeleteRows",
			Handler:    _DatabaseHelper_DeleteRows_Handler,
		},
		{
			MethodName: "Query",
			Handler:    _DatabaseHelper_Query_Handler,
		},
		{
			MethodName: "ResetDB",
			Handler:    _DatabaseHelper_ResetDB_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "database.proto",
}
