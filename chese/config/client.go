package config

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/user"
	"strings"
)

// Client represents client configuration.
type Client struct {
	// Chrome represents all information relevant to the launching and configuration of Chrome.
	Chrome struct {
		BinPath        string `json:"bin-path"`
		DataDirectory  string `json:"data-directory"`
		AllowReferrers bool   `json:"allow-referrers"`
	} `json:"chrome"`

	// Keyhole represents all information relevant to configuring the local proxy (aka keyhole)
	Keyhole struct {
		Listener string `json:"listener"`

		// Forwarders are default routes where traffic not matching local rules are proxied.
		// This implementation round-robin's between forwarders.
		Forwarders []Forwarder `json:"forwarders"`
	} `json:"keyhole"`

	// Debug tells the client to print verbose information to assist with debugging or testing.
	Debug bool `json:"debug"`

	// ResourcePath tells the client the directory which contains resource files.
	// ResourcePath may contain '~', which will be expanded to the users home directory.
	ResourcePath string `json:"resource-path"`

	// Path represents where the on-disk configuration is stored - used for writing back configuration.
	Path string `json:"-"`
}

// Forwarder represents a remote proxy server which traffic can be forwarded to.
type Forwarder struct {
	Type    string `json:"type"`
	Address string `json:"address"`

	CachedTLSConfig *tls.Config `json:"-"`
	CaCertPath      string      `json:"ca-cert-path"`
	ClientCertPath  string      `json:"client-cert-path"`
	ClientKeyPath   string      `json:"client-key-path"`
}

// ReadClientConfig loads configuration from a file on disk.
func ReadClientConfig(fpath string) (*Client, error) {
	var m = &Client{}

	confF, err := os.Open(fpath)

	if err != nil {
		return nil, err
	}
	defer confF.Close()

	dec := json.NewDecoder(confF)

	if err = dec.Decode(&m); err == io.EOF {
	} else if err != nil {
		return nil, err
	}
	m.Path = fpath

	if strings.Contains(m.ResourcePath, "~") {
		usr, err := user.Current()
		if err != nil {
			return nil, err
		}
		m.ResourcePath = strings.Replace(m.ResourcePath, "~", usr.HomeDir, -1)
	}

	return m, validateClientConfig(m)
}

func validateClientConfig(configuration *Client) error {
	if configuration.Chrome.BinPath == "" {
		return errors.New("chrome.bin-path not specified")
	}
	if configuration.Chrome.DataDirectory == "" {
		return errors.New("chrome.data-directory not specified")
	}

	if configuration.ResourcePath == "" {
		return errors.New("resource-path not specified")
	}

	if configuration.Keyhole.Listener == "" {
		return errors.New("keyhole.listener not specified")
	}
	for i, forwarder := range configuration.Keyhole.Forwarders {
		switch forwarder.Type {
		case "tcp-insecure", "tls-insecure", "tls-pinned":
		default:
			return fmt.Errorf("keyhole.forwarder[%d].type is invalid", i)
		}
		if forwarder.Type == "tls-pinned" {
			if forwarder.CaCertPath == "" {
				return fmt.Errorf("keyhole.forwarder[%d].ca-cert-path not specified", i)
			}
			if forwarder.ClientCertPath == "" {
				return fmt.Errorf("keyhole.forwarder[%d].client-cert-path not specified", i)
			}
			if forwarder.ClientKeyPath == "" {
				return fmt.Errorf("keyhole.forwarder[%d].client-key-path not specified", i)
			}
			if fileNotExists(forwarder.CaCertPath) {
				return fmt.Errorf("bad %s: %q does not exist", "ca-cert-path", forwarder.CaCertPath)
			}
			if fileNotExists(forwarder.ClientCertPath) {
				return fmt.Errorf("bad %s: %q does not exist", "client-cert-path", forwarder.ClientCertPath)
			}
			if fileNotExists(forwarder.ClientKeyPath) {
				return fmt.Errorf("bad %s: %q does not exist", "client-key-path", forwarder.ClientKeyPath)
			}
		}
	}
	return nil
}
