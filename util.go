package kcp

import (
	"net"

	ma "github.com/multiformats/go-multiaddr"
	manet "github.com/multiformats/go-multiaddr-net"
	"github.com/pkg/errors"
)

// KcpMultiaddrToNetAddr converts a kcp multiaddress to a net.addr.
func KcpMultiaddrToNetAddr(maddr ma.Multiaddr) (net.Addr, error) {
	protos := maddr.Protocols()
	last := protos[len(protos)-1]
	if last.Name != "kcp" {
		return nil, errors.Errorf("not a kcp multiaddr: %s", maddr.String())
	}

	maddrBase := maddr.Decapsulate(baseMultiaddr)
	return manet.ToNetAddr(maddrBase)
}

// NetAddrToKcpMultiaddr converts a net address to a kcp multiaddress.addr.
func NetAddrToKcpMultiaddr(addr net.Addr) (ma.Multiaddr, error) {
	maddrBase, err := manet.FromNetAddr(addr)
	if err != nil {
		return nil, err
	}

	return maddrBase.Encapsulate(baseMultiaddr), nil
}
