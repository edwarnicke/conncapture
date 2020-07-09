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

// Package conncapture provides a credentials.TransportCredential wrapper that allows capturing the net.Conn
// as the peer.Addr which can be accessed by the server using peer.FromContext
// or the client using grpc.PeerCallOption.  Be very careful in its use, as *most* things you might do with it
// *will* screwup your grpc.  There is at least one case (sending of file descriptors over unix file sockets out of band)
// in which it is useful.
package conncapture

import (
	"context"
	"net"

	"google.golang.org/grpc/credentials"
)

type transportCredentials struct {
	credentials.TransportCredentials
}

// TransportCredentials wraps cred (which may be nil) such that peer.Addr *is* the net.Conn and can be access by the server using peer.FromContext
// or the client using grpc.PeerCallOption
func TransportCredentials(cred credentials.TransportCredentials) credentials.TransportCredentials {
	return &transportCredentials{cred}
}

func (t *transportCredentials) ClientHandshake(ctx context.Context, authority string, conn net.Conn) (net.Conn, credentials.AuthInfo, error) {
	var authInfo credentials.AuthInfo
	var err error
	if t.TransportCredentials != nil {
		conn, authInfo, err = t.TransportCredentials.ClientHandshake(ctx, authority, conn)
		if err != nil {
			return nil, nil, err
		}
	}
	return ConnToAddrInfo(conn), authInfo, err
}

func (t *transportCredentials) ServerHandshake(conn net.Conn) (net.Conn, credentials.AuthInfo, error) {
	var authInfo credentials.AuthInfo
	var err error
	if t.TransportCredentials != nil {
		conn, authInfo, err = t.TransportCredentials.ServerHandshake(conn)
		if err != nil {
			return nil, nil, err
		}
	}
	return ConnToAddrInfo(conn), authInfo, err
}
func (t *transportCredentials) Clone() credentials.TransportCredentials {
	if t.TransportCredentials != nil {
		return &transportCredentials{
			TransportCredentials: t.TransportCredentials.Clone(),
		}
	}
	return &transportCredentials{}
}

func (t *transportCredentials) Info() credentials.ProtocolInfo {
	if t.TransportCredentials == nil {
		return credentials.ProtocolInfo{}
	}
	return t.TransportCredentials.Info()
}

func (t *transportCredentials) OverrideServerName(s string) error {
	if t.TransportCredentials == nil {
		return nil
	}
	return t.TransportCredentials.OverrideServerName(s)
}
