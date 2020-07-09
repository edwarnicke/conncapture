// Copyright (c) 2020 Cisco and/or its affiliates.
//
// SPDX-License-Identifier: Apache-2.0
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package conncapture_test

import (
	"context"
	"net"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/peer"
)

type testHealthServer struct {
	*testing.T
}

func NewTestHealthServer(t *testing.T) grpc_health_v1.HealthServer {
	return &testHealthServer{T: t}
}

func (t *testHealthServer) Check(ctx context.Context, request *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
	p, ok := peer.FromContext(ctx)
	require.True(t, ok)
	require.NotNil(t, p)
	require.NotNil(t, p.Addr)
	require.Equal(t, p.Addr.Network(), "unix")
	require.Equal(t, p.Addr.String(), "")
	conn, ok := p.Addr.(net.Conn)
	require.True(t, ok)
	require.NotNil(t, conn)
	return &grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_SERVING,
	}, nil
}

func (t *testHealthServer) Watch(request *grpc_health_v1.HealthCheckRequest, server grpc_health_v1.Health_WatchServer) error {
	p, ok := peer.FromContext(server.Context())
	require.True(t, ok)
	require.NotNil(t, p)
	require.NotNil(t, p.Addr)
	require.Equal(t, p.Addr.Network(), "unix")
	require.Equal(t, p.Addr.String(), "")
	conn, ok := p.Addr.(net.Conn)
	require.True(t, ok)
	require.NotNil(t, conn)
	server.Send(&grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_SERVING,
	})
	<-server.Context().Done()
	return nil
}
