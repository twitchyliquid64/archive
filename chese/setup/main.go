// chese-setup has two modes of operation:
//  - interactive command-line interface to configure a client/server
//   - Create a new CA/Server pair (defaults to values specified as flags)
//   - Create a new client certificate (defaults to values specified as flags)
//   - Edit client configuration (defaults to current config in /etc/chese/client.json)
//   - Edit server configuration (defaults to current config in /etc/chese/server.json)
//  - command-line non-interactive:
//   - Create a new CA/Server pair
//   - Creating a new client certificate

package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"time"

	termbox "github.com/nsf/termbox-go"
	"github.com/twitchyliquid64/chese/tls"
)

var (
	errAbortMenu       = errors.New("aborted")
	caCertPath         = flag.String("ca-cert", "/etc/chese/ca-cert.pem", "Path to CA certificate file")
	caKeyPath          = flag.String("ca-key", "/etc/chese/ca-key.pem", "Path to CA key file")
	servCertPath       = flag.String("serv-cert", "/etc/chese/serv-cert.pem", "Path to server certificate file")
	servKeyPath        = flag.String("serv-key", "/etc/chese/serv-key.pem", "Path to server key file")
	clientCertPath     = flag.String("client-cert", "/etc/chese/client-cert.pem", "Path to client certificate file")
	clientKeyPath      = flag.String("client-key", "/etc/chese/client-key.pem", "Path to client key file")
	clientConfigPath   = flag.String("client-conf", "DEFAULT", "Path to client configuration file")
	clientCertValidity = flag.Duration("client-cert-validity", time.Duration(time.Hour*24*28*6), "Age of issued client certs")
	serverConfigPath   = flag.String("server-conf", "/etc/chese/server.json", "Path to server configuration file")
)

func main() {
	flag.Parse()

	var err error
	defer func() {
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			os.Exit(1)
		}
	}()

	defer func() {
		if termbox.IsInit {
			termbox.Close()
		}
	}()

	if flag.NArg() == 1 {
		switch flag.Arg(0) {
		case "mint-server":
			err = tls.MakeServerCert(*servCertPath, *servKeyPath, *caCertPath, *caKeyPath, 2048)
		case "issue-cert":
			err = tls.IssueClientCert(*caCertPath, *caKeyPath, *clientCertPath, *clientKeyPath, 2048, time.Now().Add(*clientCertValidity))
		default:
			err = fmt.Errorf("unknown mode %q", flag.Arg(0))
		}
		return
	}

	err = termbox.Init()
	if err != nil {
		return
	}

	if err = mainMenu(); err != nil {
		return
	}
}
