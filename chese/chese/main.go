package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/twitchyliquid64/chese/chese/keyhole"
	"github.com/twitchyliquid64/chese/config"
)

var (
	debug bool
)

var (
	configPath string
	debugFlag  = flag.Bool("debug", false, "enable debug logging")
)

func readFlags() {
	flag.Parse()
	if flag.NArg() < 1 {
		fmt.Println("No flags or arguments --> Autopilot mode.")

		// use the config file in the home directory if one exists, otherwise in etc.
		usr, _ := user.Current()
		if fileExists(path.Join(usr.HomeDir, ".config/chese/client.json")) {
			configPath = path.Join(usr.HomeDir, ".config/chese/client.json")
		} else {
			configPath = "/etc/chese/client.json"
		}
		fmt.Printf("Reading config from %q\n", configPath)
	} else {
		configPath = os.Args[1]
	}
	if *debugFlag {
		debug = true
	}
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// initialization step following the load of config.
func setupGlobals(conf *config.Client) {
	if conf.Debug {
		debug = true // it can be set elsewhere
	}
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

	conf, err := config.ReadClientConfig(configPath)
	if err != nil {
		return
	}
	if debug {
		fmt.Println("Client configuration")
		spew.Dump(conf)
	}

	serv, err := keyhole.New(conf.Keyhole.Listener, conf.ResourcePath, debug, conf)
	if err != nil {
		return
	}
	defer serv.Close()
	if debug {
		fmt.Println("Keyhole server")
		spew.Dump(serv)
	}
	serv.Start()

	err = startChrome(conf)
	if err != nil {
		return
	}
}

func startChrome(conf *config.Client) error {
	args := []string{"--no-first-run", "--disable-default-apps", "--no-default-browser-check", "--enforce-webrtc-ip-permission-check", "--prerender-from-omnibox=disabled", "--user-data-dir=" + conf.Chrome.DataDirectory}
	if !conf.Chrome.AllowReferrers {
		args = append(args, "--no-referrers")
	}
	args = append(args, "--load-extension="+path.Join(conf.ResourcePath, "theme"))

	proxyAddr := conf.Keyhole.Listener
	if strings.HasPrefix(proxyAddr, ":") {
		proxyAddr = "localhost" + proxyAddr
	}
	args = append(args, "--proxy-server="+proxyAddr)

	args = append(args, "http://keyhole.internal/keyhole/landing")

	cmd := exec.Command(conf.Chrome.BinPath, args...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	return cmd.Run()
}
