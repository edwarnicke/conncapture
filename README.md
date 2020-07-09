conncapture is a hack to allow a grpc Client or Server to retrieve the net.Conn they are currently operating over.

Normally, this is a terrible idea.  There is at least one use case in which it makes sense:
* Sending file descriptors over unix file sockets

In order to accomplish net.Conn capture, a [credentials.TransportCredentials](https://godoc.org/google.golang.org/grpc/credentials#TransportCredentials)
wrapper is provided by [conncapture.TransportCredentials(cred credentials.TransportCredentials)](https://github.com/edwarnicke/conncapture/blob/47c400e/transport_credentials.go#L37)
to wrap your real TLS or other TransportCredentials (or nil if you have none).

The result is that the [peer.Addr](https://godoc.org/google.golang.org/grpc/peer#Peer) will be a wrapped version of the actual [net.Conn](https://golang.org/pkg/net/#Conn)
that meets the [net.Addr](https://golang.org/pkg/net/#Addr) interface.  This can be retrieved in servers with [peer.FromContext(ctx)](https://godoc.org/google.golang.org/grpc/peer#FromContext) or in clients with
[grpc.Peer(p *peer.Peer)](https://godoc.org/google.golang.org/grpc#Peer).
