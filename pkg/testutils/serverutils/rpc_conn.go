// Copyright 2025 The Cockroach Authors.
//
// Use of this software is governed by the CockroachDB Software License
// included in the /LICENSE file.

package serverutils

import (
	"github.com/cockroachdb/cockroach/pkg/kv/kvpb"
	"github.com/cockroachdb/cockroach/pkg/server/serverpb"
	"github.com/cockroachdb/cockroach/pkg/sql/execinfrapb"
	"github.com/cockroachdb/cockroach/pkg/ts/tspb"
	"google.golang.org/grpc"
	"storj.io/drpc"
)

// grpcConn implements server.RPCConn that provides methods to create various
// RPC clients. This allows the CLI to interact with the server using gRPC
// without exposing the underlying connection details.
type grpcConn struct {
	conn *grpc.ClientConn
}

func FromGRPCConn(conn *grpc.ClientConn) RPCConn {
	return &grpcConn{conn: conn}
}

func (c *grpcConn) NewStatusClient() serverpb.RPCStatusClient {
	return serverpb.NewGRPCStatusClientAdapter(c.conn)
}

func (c *grpcConn) NewAdminClient() serverpb.RPCAdminClient {
	return serverpb.NewGRPCAdminClientAdapter(c.conn)
}

func (c *grpcConn) NewInitClient() serverpb.RPCInitClient {
	return serverpb.NewGRPCInitClientAdapter(c.conn)
}

func (c *grpcConn) NewTimeSeriesClient() tspb.RPCTimeSeriesClient {
	return tspb.NewGRPCTimeSeriesClientAdapter(c.conn)
}

func (c *grpcConn) NewInternalClient() kvpb.RPCInternalClient {
	return kvpb.NewGRPCInternalClientAdapter(c.conn)
}

func (c *grpcConn) NewDistSQLClient() execinfrapb.RPCDistSQLClient {
	return execinfrapb.NewGRPCDistSQLClientAdapter(c.conn)
}

// drpcConn implements server.RPCConn that provides methods to create various
// RPC clients. This allows the CLI to interact with the server using DRPC
// without exposing the underlying connection details.
type drpcConn struct {
	conn drpc.Conn
}

func FromDRPCConn(conn drpc.Conn) RPCConn {
	return &drpcConn{conn: conn}
}

func (c *drpcConn) NewStatusClient() serverpb.RPCStatusClient {
	return serverpb.NewDRPCStatusClientAdapter(c.conn)
}

func (c *drpcConn) NewAdminClient() serverpb.RPCAdminClient {
	return serverpb.NewDRPCAdminClientAdapter(c.conn)
}

func (c *drpcConn) NewInitClient() serverpb.RPCInitClient {
	return serverpb.NewDRPCInitClientAdapter(c.conn)
}

func (c *drpcConn) NewTimeSeriesClient() tspb.RPCTimeSeriesClient {
	return tspb.NewDRPCTimeSeriesClientAdapter(c.conn)
}

func (c *drpcConn) NewInternalClient() kvpb.RPCInternalClient {
	return kvpb.NewDRPCInternalClientAdapter(c.conn)
}

func (c *drpcConn) NewDistSQLClient() execinfrapb.RPCDistSQLClient {
	return execinfrapb.NewDRPCDistSQLClientAdapter(c.conn)
}
