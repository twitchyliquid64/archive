package keyhole

import (
	rtls "crypto/tls"
	"errors"
	"math/rand"
	"net"

	"github.com/twitchyliquid64/chese/config"
	"github.com/twitchyliquid64/chese/tls"
)

type dialFunc func(network, addr string) (net.Conn, error)

func dialer(conf *config.Client) dialFunc {
	forwarder := conf.Keyhole.Forwarders[rand.Intn(len(conf.Keyhole.Forwarders))]
	return func(network, addr string) (net.Conn, error) {
		return dialForwarder(forwarder)
	}
}

// dial selects a forwarder based on configuration, before opening a socket to it
func dial(conf *config.Client) (net.Conn, error) {
	if len(conf.Keyhole.Forwarders) == 0 {
		return nil, errors.New("no forwarders configured")
	}

	// initial implementation round-robin's between forwarders.
	return dialForwarder(conf.Keyhole.Forwarders[rand.Intn(len(conf.Keyhole.Forwarders))])
}

func dialForwarder(forwarder config.Forwarder) (net.Conn, error) {
	switch forwarder.Type {
	case "tcp-insecure":
		return net.Dial("tcp", forwarder.Address)
	case "tls-pinned":
		if forwarder.CachedTLSConfig == nil {
			var err error
			forwarder.CachedTLSConfig, err = tls.ConfigureTLS(forwarder.ClientCertPath, forwarder.ClientKeyPath, forwarder.CaCertPath)
			if err != nil {
				return nil, err
			}
		}
		return rtls.Dial("tcp", forwarder.Address, forwarder.CachedTLSConfig)
	}
	return nil, errors.New("invalid forwarder type")
}
