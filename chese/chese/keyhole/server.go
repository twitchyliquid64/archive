package keyhole

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"time"

	"github.com/twitchyliquid64/chese/config"
)

const (
	keyholeHost   = "keyhole.internal"
	keyholePrefix = "/keyhole/"
)

// Server represents a configurable and runnable keyhole instance.
type Server struct {
	resourceDir string
	serv        *http.Server
	handler     *http.ServeMux
	debug       bool
	conf        *config.Client

	// indicates to server go-routines to shut down.
	close chan bool
}

// New creates a new keyhole instance.
func New(listener, resourceDir string, debug bool, conf *config.Client) (*Server, error) {
	if err := testPortOpen(listener); err != nil {
		return nil, err
	}

	s := &Server{
		resourceDir: path.Join(resourceDir, "keyhole"),
		handler:     http.NewServeMux(),
		serv:        &http.Server{Addr: listener},
		close:       make(chan bool),
		debug:       debug,
		conf:        conf,
	}
	s.serv.Handler = s
	s.handler.Handle("/keyhole/static/", http.StripPrefix("/keyhole/static/", http.FileServer(http.Dir(path.Join(resourceDir, "keyhole", "static")))))
	s.handler.HandleFunc("/keyhole/landing", s.handleLanding)
	return s, nil
}

// Close shuts down the server and releases its internal resources.
func (s *Server) Close() error {
	err := s.serv.Close()
	if err != nil {
		return err
	}
	close(s.close)
	return nil
}

// Start is called to start listening and serve requests.
func (s *Server) Start() {
	go func() {
		err := s.serv.ListenAndServe()
		if err != http.ErrServerClosed {
			fmt.Fprintf(os.Stderr, "[KEYHOLE] ListenAndServe() error: %v\n", err)
		}
	}()
}

func (s *Server) handleLanding(rw http.ResponseWriter, req *http.Request) {
	renderPage(path.Join(s.resourceDir, "landing.html"), nil, rw)
}

// ServeHTTP is called for all requests keyhole recieves.
func (s *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if s.debug {
		fmt.Printf("[KEYHOLE-REQUEST] %s %s %s\n", req.Method, req.RequestURI, req.Proto)
	}

	if req.Host == keyholeHost && req.Method == "GET" && strings.HasPrefix(req.URL.Path, keyholePrefix) {
		s.handler.ServeHTTP(w, req)
		return
	}

	// proxy the request
	if req.Method == "CONNECT" {
		s.proxyConnect(w, req)
		return
	}
	s.proxyHTTP(w, req)

}

type writeCloserConn interface {
	CloseWrite() error
}

// proxyConnect handles a CONNECT tunnel, making an identical request to our proxy
// and transparently tunnelling the data through it.
func (s *Server) proxyConnect(rw http.ResponseWriter, req *http.Request) {
	conn, err := dial(s.conf)
	if err != nil {
		http.Error(rw, "Server Error", http.StatusInternalServerError)
		fmt.Printf("[KEYHOLE-CONNECT] dial() Error: %v\n", err)
		return
	}
	defer conn.Close()

	proxyReq := &http.Request{
		Method: "CONNECT",
		URL:    &url.URL{Opaque: req.RequestURI},
		Host:   req.Host,
	}
	proxyReq.Write(conn)
	br := bufio.NewReader(conn)
	resp, err := http.ReadResponse(br, proxyReq)
	if err != nil {
		http.Error(rw, "Server Error", http.StatusBadGateway)
		fmt.Printf("[KEYHOLE-CONNECT] readResponse() Error: %v\n", err)
		return
	}
	if resp.StatusCode != 200 {
		http.Error(rw, "Server Error", http.StatusBadGateway)
		fmt.Printf("[KEYHOLE-CONNECT] remote Error: %v\n", strings.SplitN(resp.Status, " ", 2))
		return
	}
	rw.WriteHeader(http.StatusOK)

	browserConn, bufrw, err := rw.(http.Hijacker).Hijack()
	if err != nil {
		fmt.Printf("[KEYHOLE-CONNECT] hijack Error: %v\n", err)
		return
	}
	defer browserConn.Close()

	browserWriteDone := make(chan bool)
	go func() {
		_, err2 := io.Copy(conn, bufrw)
		fmt.Printf("[KEYHOLE-COPY] browser->conn Error: %v\n", err2)
		conn.(writeCloserConn).CloseWrite()
		close(browserWriteDone)
	}()

	_, err = io.Copy(bufrw, conn)
	fmt.Printf("[KEYHOLE-COPY] conn->browser Error: %v\n", err)
	browserConn.(writeCloserConn).CloseWrite()
	select {
	case <-browserWriteDone:
	case <-time.After(time.Second * 15):
	}
}

// proxyHTTP handles non-CONNECT HTTP requests, proxying the request to
// a forwarder.
func (s *Server) proxyHTTP(rw http.ResponseWriter, req *http.Request) {
	req.RequestURI = ""
	delHopHeaders(req.Header)

	proxy, _ := url.Parse("http://127.0.0.1:8080")
	client := &http.Client{
		Transport: &http.Transport{
			Proxy:                 http.ProxyURL(proxy),
			Dial:                  dialer(s.conf),
			MaxIdleConns:          0,
			ExpectContinueTimeout: 1 * time.Second,
		},
	}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(rw, "Server Error", http.StatusInternalServerError)
		fmt.Printf("[KEYHOLE-REQUEST] client.Do(%q) Error: %v\n", req.URL.Path, err)
		return
	}
	defer resp.Body.Close()

	delHopHeaders(resp.Header)
	copyHeaders(rw.Header(), resp.Header)
	rw.WriteHeader(resp.StatusCode)
	io.Copy(rw, resp.Body)
}

// returns an error if a listener address cannot be opened.
func testPortOpen(listener string) error {
	ln, err := net.Listen("tcp", listener)
	if err != nil {
		return err
	}
	return ln.Close()
}

// The following is credited to https://gist.github.com/yowu/f7dc34bd4736a65ff28d
// Hop-by-hop headers. These are removed when sent to the backend.
// http://www.w3.org/Protocols/rfc2616/rfc2616-sec13.html
var hopHeaders = []string{
	"Connection",
	"Keep-Alive",
	"Proxy-Authenticate",
	"Proxy-Authorization",
	"Te", // canonicalized version of "TE"
	"Trailers",
	"Transfer-Encoding",
	"Upgrade",
}

func delHopHeaders(header http.Header) {
	for _, h := range hopHeaders {
		header.Del(h)
	}
}

func copyHeaders(responseHeaders, sourceHeaders http.Header) {
	for k := range responseHeaders {
		responseHeaders.Del(k)
	}
	for k, vs := range sourceHeaders {
		for _, v := range vs {
			responseHeaders.Add(k, v)
		}
	}
}
