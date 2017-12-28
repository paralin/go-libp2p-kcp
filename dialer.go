package kcp

import (
	"context"

	"github.com/Sirupsen/logrus"
	tpt "github.com/libp2p/go-libp2p-transport"
	ma "github.com/multiformats/go-multiaddr"
	manet "github.com/multiformats/go-multiaddr-net"
	kcpgo "github.com/xtaci/kcp-go"
	"github.com/xtaci/smux"
)

// dialer dials via kcp.
type dialer struct {
	transport *KcpTransport
}

// Dial connects via kcp to a remote address.
func (d *dialer) Dial(raddr ma.Multiaddr) (tpt.Conn, error) {
	return d.DialContext(context.Background(), raddr)
}

// DialContext connects via kcp to a remote address.
// TODO: respect the context
func (d *dialer) DialContext(ctx context.Context, raddr ma.Multiaddr) (tpt.Conn, error) {
	na, err := manet.ToNetAddr(raddr)
	if err != nil {
		return nil, err
	}

	logrus.WithField("addr", na.String()).Info("Dialing via kcp")
	kcpConn, err := kcpgo.Dial(na.String())
	if err != nil {
		return nil, err
	}

	smuxSess, err := smux.Client(kcpConn, smuxConf())
	if err != nil {
		kcpConn.Close()
		return nil, err
	}

	// TODO: In the initial version use this to emulate tcp.
	baseStream, err := smuxSess.OpenStream()
	if err != nil {
		smuxSess.Close()
		kcpConn.Close()
		return nil, err
	}

	logrus.Info("Wrapping smux stream")
	mconn, err := manet.WrapNetConn(baseStream)
	if err != nil {
		kcpConn.Close()
		return nil, err
	}

	defer logrus.Info("Done DialContext")
	return &conn{Conn: mconn, transport: d.transport}, nil
}

// Matches checks if the dialer matches the maddr.
func (d *dialer) Matches(addr ma.Multiaddr) bool {
	return KcpFmt.Matches(addr)
}
