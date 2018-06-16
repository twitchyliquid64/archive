package proxy

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"path"
	"sync"

	"github.com/elazarl/goproxy"
	"github.com/twitchyliquid64/chese/config"
	"github.com/twitchyliquid64/chese/tls"
)

// Server represents a configurable and runnable proxy.
type Server struct {
	resourceDir string
	proxy       *goproxy.ProxyHttpServer
	serv        *http.Server
	debug       bool

	routines sync.WaitGroup
}

// New creates a new proxy server.
func New(listener, resourceDir string, tlsInfo config.ServerCerts, debug bool) (*Server, error) {
	if err := testPortOpen(listener); err != nil {
		return nil, err
	}

	tlsConf, err := tls.ConfigureTLS(tlsInfo.MainCertPath, tlsInfo.MainKeyPath, tlsInfo.CACertPath)
	if err != nil {
		return nil, err
	}

	s := &Server{
		resourceDir: path.Join(resourceDir, "proxy"),
		proxy:       goproxy.NewProxyHttpServer(),
		serv:        &http.Server{Addr: listener, TLSConfig: tlsConf},
		debug:       debug,
	}
	s.proxy.Verbose = debug
	s.serv.Handler = s.proxy //TODO: Fixme
	return s, nil
}

// Close shuts down the server and releases its internal resources.
func (s *Server) Close() error {
	err := s.serv.Close()
	if err != nil {
		return err
	}
	s.routines.Wait()
	return nil
}

// Start is called to start listening and serve requests.
func (s *Server) Start() {
	go func() {
		s.routines.Add(1)
		defer s.routines.Done()
		err := s.serv.ListenAndServeTLS("", "")
		if err != http.ErrServerClosed {
			fmt.Fprintf(os.Stderr, "[PROXY] ListenAndServe() error: %v\n", err)
		}
	}()
}

// ServeHTTP is called for all requests keyhole recieves.
func (s *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if s.debug {
		fmt.Printf("[PROXY-REQUEST] %s (%s) %s\n", req.Method, req.Host, req.URL.Path)
	}
	s.proxy.ServeHTTP(w, req)
}

func testPortOpen(listener string) error {
	ln, err := net.Listen("tcp", listener)
	if err != nil {
		return err
	}
	return ln.Close()
}
