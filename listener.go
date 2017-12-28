package kcp

import (
	tpt "github.com/libp2p/go-libp2p-transport"
	manet "github.com/multiformats/go-multiaddr-net"
	"github.com/xtaci/smux"
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

	// maConn is wrapping an accepted KCP session.
	smuxSess, err := smux.Server(maConn, smuxConf())
	if err != nil {
		maConn.Close()
		return nil, err
	}

	stream, err := smuxSess.AcceptStream()
	if err != nil {
		if stream != nil {
			stream.Close()
		}
		smuxSess.Close()
		maConn.Close()
		return nil, err
	}

	streamWrap, err := manet.WrapNetConn(stream)
	if err != nil {
		if streamWrap != nil {
			streamWrap.Close()
		}
		stream.Close()
		smuxSess.Close()
		maConn.Close()
		return nil, err
	}

	return &conn{Conn: streamWrap, transport: l.transport}, nil
}
