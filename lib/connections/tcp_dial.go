// Copyright (C) 2016 The Syncthing Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this file,
// You can obtain one at http://mozilla.org/MPL/2.0/.

package connections

import (
	"crypto/tls"
	"net/url"
	"time"
	//
	//	"github.com/syncthing/syncthing/lib/config"
	//	"github.com/syncthing/syncthing/lib/dialer"
	//	"github.com/syncthing/syncthing/lib/protocol"

	"github.com/bitgitsync/syncthing/lib/config"
	"github.com/bitgitsync/syncthing/lib/dialer"
	"github.com/bitgitsync/syncthing/lib/protocol"
)

const tcpPriority = 10

func init() {
	factory := &tcpDialerFactory{}
	for _, scheme := range []string{"tcp", "tcp4", "tcp6"} {
		dialers[scheme] = factory
	}
}

type tcpDialer struct {
	cfg    *config.Wrapper
	tlsCfg *tls.Config
}

func (d *tcpDialer) Dial(id protocol.DeviceID, uri *url.URL) (IntermediateConnection, error) {
	uri = fixupPort(uri)

	conn, err := dialer.DialTimeout(uri.Scheme, uri.Host, 10*time.Second)
	if err != nil {
		l.Debugln(err)
		return IntermediateConnection{}, err
	}

	tc := tls.Client(conn, d.tlsCfg)
	err = tlsTimedHandshake(tc)
	if err != nil {
		tc.Close()
		return IntermediateConnection{}, err
	}

	return IntermediateConnection{tc, "TCP (Client)", tcpPriority}, nil
}

func (d *tcpDialer) RedialFrequency() time.Duration {
	return time.Duration(d.cfg.Options().ReconnectIntervalS) * time.Second
}

type tcpDialerFactory struct{}

func (tcpDialerFactory) New(cfg *config.Wrapper, tlsCfg *tls.Config) genericDialer {
	return &tcpDialer{
		cfg:    cfg,
		tlsCfg: tlsCfg,
	}
}

func (tcpDialerFactory) Priority() int {
	return tcpPriority
}

func (tcpDialerFactory) Enabled(cfg config.Configuration) bool {
	return true
}

func (tcpDialerFactory) String() string {
	return "TCP Dialer"
}
