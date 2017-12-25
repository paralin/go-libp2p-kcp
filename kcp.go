// Package kcp implements a kcp-go based transport for go-libp2p.
package kcp

import (
	autil "github.com/libp2p/go-addr-util"
	tpt "github.com/libp2p/go-libp2p-transport"
	ma "github.com/multiformats/go-multiaddr"
	manet "github.com/multiformats/go-multiaddr-net"
	"github.com/pkg/errors"
	mafmt "github.com/whyrusleeping/mafmt"
	kcpgo "github.com/xtaci/kcp-go"
)

// KcpProtocol is the multiaddr protocol definition for this transport.
var KcpProtocol = ma.Protocol{
	// Code is the protocol code.
	// TODO: determine how to reserve one
	Code:  482,
	Name:  "kcp",
	VCode: ma.CodeToVarint(482),
}

// KcpFmt is the multiaddr formatter for KcpProtocol
var KcpFmt = mafmt.And(mafmt.UDP, mafmt.Base(KcpProtocol.Code))

var KcpCodec = &manet.NetCodec{
	NetAddrNetworks:  []string{"kcp"},
	ProtocolName:     "kcp",
	ConvertMultiaddr: KcpMultiaddrToNetAddr,
	ParseNetAddr:     NetAddrToKcpMultiaddr,
}

// baseMultiaddr is a multiaddr with just the KCP address in it for decapsulation.
var baseMultiaddr ma.Multiaddr

func init() {
	err := ma.AddProtocol(KcpProtocol)
	if err != nil {
		panic(errors.Errorf("error registering kcp protocol: %s", err.Error()))
	}

	manet.RegisterNetCodec(KcpCodec)
	baseMultiaddr, err = ma.NewMultiaddr("/kcp")
	if err != nil {
		panic(err)
	}
	autil.AddTransport("/ip4/udp/kcp")
}

// KcpTransport is the KCP go-libp2p transport.
type KcpTransport struct{}

// type assertion
var _ tpt.Transport = (*KcpTransport)(nil)

// Matches checks if a multiaddr matches the kcp protocol.
func (t *KcpTransport) Matches(a ma.Multiaddr) bool {
	return KcpFmt.Matches(a)
}

// Dialer is the dialer for the kcp transport
func (t *KcpTransport) Dialer(_ ma.Multiaddr, opts ...tpt.DialOpt) (tpt.Dialer, error) {
	return &dialer{transport: t}, nil
}

// Listen is the listener for the kcp transport.
func (t *KcpTransport) Listen(a ma.Multiaddr) (tpt.Listener, error) {
	na, err := manet.ToNetAddr(a)
	if err != nil {
		return nil, err
	}

	kcpList, err := kcpgo.Listen(na.String())
	if err != nil {
		return nil, err
	}

	tlist, err := manet.WrapNetListener(kcpList)
	if err != nil {
		kcpList.Close()
		return nil, err
	}

	return &listener{Listener: tlist, transport: t}, nil
}
