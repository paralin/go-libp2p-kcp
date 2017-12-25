package kcp

import (
	tpt "github.com/libp2p/go-libp2p-transport"
	manet "github.com/multiformats/go-multiaddr-net"
)

// conn wraps a manet.Conn
type conn struct {
	manet.Conn
	transport *KcpTransport
}

// Transport returns the underlying transport
func (c *conn) Transport() tpt.Transport {
	return c.transport
}
