package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/davecgh/go-spew/spew"
	"github.com/twitchyliquid64/chese/config"
	"github.com/twitchyliquid64/chese/server/proxy"
)

var (
	debug bool
)

var (
	debugFlag = flag.Bool("debug", false, "enable debug logging")
)

func readFlags() {
	flag.Parse()
	if flag.NArg() < 1 {
		fmt.Fprintf(os.Stderr, "USAGE: %s [<flags>] <config-path>\n", os.Args[0])
		os.Exit(1)
	}
	if *debugFlag {
		debug = true
	}
}

// initialization step following the load of config.
func setupGlobals(conf *config.Server) {
	if conf.Debug {
		debug = true // it can be set elsewhere
	}
}

func waitInterrupt() os.Signal {
	sig := make(chan os.Signal, 2)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	return <-sig
}

func main() {
	readFlags()

	// print any errors after all other
	// defers (safe shutdown routines) have finished
	var err error
	defer func() {
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(2)
		}
	}()

	conf, err := config.ReadServerConfig(flag.Arg(0))
	if err != nil {
		return
	}
	if debug {
		fmt.Println("Server configuration")
		spew.Dump(conf)
	}

	serv, err := proxy.New(conf.Listener, conf.ResourcePath, conf.TLS, debug)
	if err != nil {
		return
	}
	if debug {
		fmt.Printf("Server listening on %q\n", conf.Listener)
	}
	go serv.Start()
	defer serv.Close()

	for {
		sig := waitInterrupt()
		fmt.Printf(" Got %v\n", sig)
		if sig != syscall.SIGHUP {
			break
		}
	}
}
