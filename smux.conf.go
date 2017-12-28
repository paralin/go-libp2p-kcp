package kcp

import (
	"github.com/xtaci/smux"
	"time"
)

// smuxConf is the config for smux server and client
func smuxConf() (conf *smux.Config) {
	conf = smux.DefaultConfig()
	// TODO: potentially tweak timeouts
	conf.KeepAliveInterval = time.Second * 5
	conf.KeepAliveTimeout = time.Second * 13
	return
}
