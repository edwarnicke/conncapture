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
	"crypto/tls"
	"io/ioutil"
	"net"
	"path/filepath"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/health/grpc_health_v1"

	"github.com/edwarnicke/conncapture"
)

func TestConnCaptureTransportCredentials_Check_NoTLS(t *testing.T) {
	cred := conncapture.TransportCredentials(nil)
	dir, err := ioutil.TempDir("", "")
	require.NoError(t, err)
	path := filepath.Join(dir, "socket")
	listener, err := net.Listen("unix", path)
	require.NoError(t, err)
	defer func() { _ = listener.Close() }()
	server := grpc.NewServer(grpc.Creds(cred))
	srv := NewTestHealthServer(t)
	grpc_health_v1.RegisterHealthServer(server, srv)
	go func() {
		_ = server.Serve(listener)
	}()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cc, err := grpc.DialContext(ctx, "unix://"+path, grpc.WithTransportCredentials(cred))
	defer func() { _ = cc.Close() }()
	require.NoError(t, err)
	client := grpc_health_v1.NewHealthClient(cc)
	client = NewTestHealthClient(t, client)
	_, err = client.Check(ctx, &grpc_health_v1.HealthCheckRequest{
		Service: "foo",
	})
	require.NoError(t, err)
}

func TestConnCaptureTransportCredentials_Check_TLS(t *testing.T) {
	cert, err := SelfSignedCert()
	require.NoError(t, err)
	require.NotNil(t, cert)
	cred := conncapture.TransportCredentials(credentials.NewTLS(&tls.Config{Certificates: []tls.Certificate{*cert}, InsecureSkipVerify: true})) // #nosec
	dir, err := ioutil.TempDir("", "")
	require.NoError(t, err)
	path := filepath.Join(dir, "socket")
	listener, err := net.Listen("unix", path)
	require.NoError(t, err)
	defer func() { _ = listener.Close() }()
	server := grpc.NewServer(grpc.Creds(cred))
	srv := NewTestHealthServer(t)
	grpc_health_v1.RegisterHealthServer(server, srv)
	go func() {
		_ = server.Serve(listener)
	}()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cc, err := grpc.DialContext(ctx, "unix://"+path, grpc.WithTransportCredentials(cred))
	defer func() { _ = cc.Close() }()
	require.NoError(t, err)
	client := grpc_health_v1.NewHealthClient(cc)
	client = NewTestHealthClient(t, client)
	_, err = client.Check(ctx, &grpc_health_v1.HealthCheckRequest{
		Service: "foo",
	})
	require.NoError(t, err)
}

func TestConnToAddrInfo(t *testing.T) {
	dir, err := ioutil.TempDir("", "")
	require.NoError(t, err)
	path := filepath.Join(dir, "socket")
	listener, err := net.Listen("unix", path)
	require.NoError(t, err)
	defer func() { _ = listener.Close() }()
	go func() {
		_, _ = listener.Accept()
	}()

	conn, err := net.Dial("unix", path)
	require.NoError(t, err)
	require.NotNil(t, conn)
	defer func() { _ = conn.Close() }()
	conn = conncapture.ConnToAddrInfo(conn)
	c, ok := conn.(interface {
		SyscallConn() (syscall.RawConn, error)
	})
	require.True(t, ok)
	require.NotNil(t, c)
	raw, err := c.SyscallConn()
	require.NoError(t, err)
	require.NotNil(t, raw)
}
