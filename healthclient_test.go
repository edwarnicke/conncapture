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
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/peer"
)

type testHealthClient struct {
	*testing.T
	client grpc_health_v1.HealthClient
}

func NewTestHealthClient(t *testing.T, client grpc_health_v1.HealthClient) grpc_health_v1.HealthClient {
	return &testHealthClient{T: t, client: client}
}

func (t *testHealthClient) Check(ctx context.Context, in *grpc_health_v1.HealthCheckRequest, opts ...grpc.CallOption) (*grpc_health_v1.HealthCheckResponse, error) {
	p := &peer.Peer{}
	peerOpt := grpc.PeerCallOption{p}
	rsp, err := t.client.Check(ctx, in, append(opts, peerOpt)...)
	require.NoError(t, err)
	require.NotNil(t, p)
	require.NotNil(t, p.Addr)
	conn, ok := p.Addr.(net.Conn)
	require.True(t, ok)
	require.NotNil(t, conn)
	require.Equal(t, p.Addr.Network(), conn.RemoteAddr().Network())
	require.Equal(t, p.Addr.String(), conn.RemoteAddr().String())
	return rsp, nil
}

func (t *testHealthClient) Watch(ctx context.Context, in *grpc_health_v1.HealthCheckRequest, opts ...grpc.CallOption) (grpc_health_v1.Health_WatchClient, error) {
	p := &peer.Peer{}
	peerOpt := grpc.PeerCallOption{p}
	watchClient, err := t.client.Watch(ctx, in, append(opts, peerOpt)...)
	require.NoError(t, err)
	require.NotNil(t, p)
	require.NotNil(t, p.Addr)
	conn, ok := p.Addr.(net.Conn)
	require.True(t, ok)
	require.NotNil(t, conn)
	require.Equal(t, p.Addr.Network(), conn.RemoteAddr().Network())
	require.Equal(t, p.Addr.String(), conn.RemoteAddr().String())
	return watchClient, nil
}
