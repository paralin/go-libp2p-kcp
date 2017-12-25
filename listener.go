package kcp

import (
	tpt "github.com/libp2p/go-libp2p-transport"
	manet "github.com/multiformats/go-multiaddr-net"
)

// listener wraps a net.Listener.
type listener struct {
	manet.Listener
	transport *KcpTransport
}

// Accept accepts an incoming conn.
func (l *listener) Accept() (tpt.Conn, error) {
	maConn, err := l.Listener.Accept()
	if err != nil {
		return nil, err
	}

	return &conn{Conn: maConn, transport: l.transport}, nil
}
