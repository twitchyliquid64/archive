package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/user"
	"strings"
)

// Server represents server configuration.
type Server struct {
	// Listener is the address which the server listens for requests.
	Listener string `json:"listener"`

	// TLS contains configuration for server security.
	TLS ServerCerts `json:"tls"`

	// Debug tells the client to print verbose information to assist with debugging or testing.
	Debug bool `json:"debug"`

	// ResourcePath tells the client the directory which contains resource files.
	// ResourcePath may contain '~', which will be expanded to the users home directory.
	ResourcePath string `json:"resource-path"`

	// Path represents where the on-disk configuration is stored - used for writing back configuration.
	Path string `json:"-"`
}

// ServerCerts stores the paths to TLS certificates and keys needed to run the server.
// All certificates and keys are PEM encoded.
type ServerCerts struct {
	CACertPath   string `json:"ca-cert-path"`
	MainCertPath string `json:"server-cert-path"`
	MainKeyPath  string `json:"server-key-path"`
}

// ReadServerConfig loads configuration from a file on disk.
func ReadServerConfig(fpath string) (*Server, error) {
	var m = &Server{}

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

	return m, validateServerConfig(m)
}

func validateServerConfig(configuration *Server) error {
	if configuration.Listener == "" {
		return errors.New("listener not specified")
	}
	if configuration.ResourcePath == "" {
		return errors.New("resource-path not specified")
	}

	if configuration.TLS.MainCertPath == "" {
		return errors.New("tls.server-cert-path not specified")
	}
	if configuration.TLS.MainKeyPath == "" {
		return errors.New("tls.server-key-path not specified")
	}
	if configuration.TLS.CACertPath == "" {
		return errors.New("tls.ca-cert-path not specified")
	}
	if fileNotExists(configuration.TLS.MainCertPath) {
		return fmt.Errorf("bad %s: %q does not exist", "tls.server-cert-path", configuration.TLS.MainCertPath)
	}
	if fileNotExists(configuration.TLS.MainKeyPath) {
		return fmt.Errorf("bad %s: %q does not exist", "tls.server-key-path", configuration.TLS.MainKeyPath)
	}
	if fileNotExists(configuration.TLS.CACertPath) {
		return fmt.Errorf("bad %s: %q does not exist", "tls.ca-cert-path", configuration.TLS.CACertPath)
	}
	return nil
}

func fileNotExists(path string) bool {
	_, err := os.Stat(path)
	return os.IsNotExist(err)
}
