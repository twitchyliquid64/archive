package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"os/user"
	"path"

	termbox "github.com/nsf/termbox-go"
	"github.com/twitchyliquid64/chese/config"
)

func checkCreateUserConfig() error {
	usr, err := user.Current()
	if err != nil {
		return err
	}
	if fileNotExists(path.Join(usr.HomeDir, ".config/chese/client.json")) {
		if fileNotExists(path.Join(usr.HomeDir, ".config/chese")) { //make folder
			if err := os.MkdirAll(path.Join(usr.HomeDir, ".config/chese"), 0750); err != nil {
				return err
			}
		}
		d, err := ioutil.ReadFile("/etc/chese/client.json")
		if err != nil {
			return err
		}
		err = ioutil.WriteFile(path.Join(usr.HomeDir, ".config/chese/client.json"), d, 0750)
		if err != nil {
			return err
		}
	}
	*clientConfigPath = path.Join(usr.HomeDir, ".config/chese/client.json")
	return nil
}

func editClientConfig() error {
	// if they dont already have a per-user config, copy one across, creating if necessary.
	// now, all config will be written to the per-user config at ~/.config/chese/client.json
	if *clientConfigPath == "DEFAULT" {
		if err := checkCreateUserConfig(); err != nil {
			return err
		}
	}

	var clientConf config.Client
	var forwarderConf config.Forwarder
	if !fileNotExists(*clientConfigPath) {
		d, err := ioutil.ReadFile(*clientConfigPath)
		if err != nil {
			return err
		}
		json.Unmarshal(d, &clientConf)
		if len(clientConf.Keyhole.Forwarders) > 0 {
			forwarderConf = clientConf.Keyhole.Forwarders[0]
		}
	}

	form := FormUI{
		Fg:       termbox.ColorBlack,
		Bg:       termbox.ColorYellow,
		Title:    "Chese Setup --> Configure Local Client",
		Subtitle: *clientConfigPath,
		Items: []formItem{
			formItem{
				Name:  "Chrome binary",
				Value: defaultString(clientConf.Chrome.BinPath, "/usr/bin/google-chrome-stable"),
			},
			formItem{
				Name:  "Profile storage directory",
				Value: defaultString(clientConf.Chrome.DataDirectory, "/tmp/chese-test"),
			},
			formItem{
				Name:  "Listener",
				Value: defaultString(clientConf.Keyhole.Listener, ":8427"),
			},
			formItem{
				Name:  "Resources path",
				Value: defaultString(clientConf.ResourcePath, "/usr/share/chese-client/resources"),
			},
			formItem{
				Name:  "Forwarder Address",
				Value: defaultString(forwarderConf.Address, ""),
			},
			formItem{
				Name:  "CA-cert path",
				Value: defaultString(forwarderConf.CaCertPath, ""),
			},
			formItem{
				Name:  "Client-cert path",
				Value: defaultString(forwarderConf.ClientCertPath, ""),
			},
			formItem{
				Name:  "Client-key path",
				Value: defaultString(forwarderConf.ClientKeyPath, ""),
			},
		},
	}
	form.Verify = func(items []formItem) bool {
		isValid := true
		for i, item := range items {
			if item.Value == "" {
				isValid = false
				items[i].FgOverride = termbox.ColorRed
			} else {
				items[i].FgOverride = 0
			}
		}

		if items[0].Value != "" && items[0].Value != "." && fileNotExists(path.Dir(items[0].Value)) {
			form.Statustext = fmt.Sprintf("Directory %q does not exist!", path.Dir(items[0].Value))
		}

		_, _, listenerErr := net.SplitHostPort(items[2].Value)
		if listenerErr != nil {
			isValid = false
			items[2].FgOverride = termbox.ColorRed
			form.Statustext = fmt.Sprintf("Invalid listener address (error %q)", listenerErr.Error())
		} else {
			items[2].FgOverride = 0
		}

		_, _, forwarderAddrErr := net.SplitHostPort(items[4].Value)
		if forwarderAddrErr != nil {
			isValid = false
			items[4].FgOverride = termbox.ColorRed
			form.Statustext = fmt.Sprintf("Invalid forwarder address (error %q)", forwarderAddrErr.Error())
		} else {
			items[4].FgOverride = 0
		}

		for i := 5; i < 8; i++ {
			if items[i].Value == "" || fileNotExists(items[i].Value) {
				isValid = false
				items[i].FgOverride = termbox.ColorRed
			} else {
				items[i].FgOverride = 0
			}
		}

		if isValid {
			form.Statustext = "Please fill in the form."
		}
		return isValid
	}
	err := form.Run()
	if err != nil {
		return err
	}
	termbox.Close()

	clientConf.Chrome.BinPath = form.Items[0].Value
	clientConf.Chrome.DataDirectory = form.Items[1].Value
	clientConf.Keyhole.Listener = form.Items[2].Value
	clientConf.ResourcePath = form.Items[3].Value
	if len(clientConf.Keyhole.Forwarders) == 0 {
		clientConf.Keyhole.Forwarders = []config.Forwarder{config.Forwarder{Type: "tls-pinned"}}
	}
	clientConf.Keyhole.Forwarders[0].Address = form.Items[4].Value
	clientConf.Keyhole.Forwarders[0].CaCertPath = form.Items[5].Value
	clientConf.Keyhole.Forwarders[0].ClientCertPath = form.Items[5].Value
	clientConf.Keyhole.Forwarders[0].ClientKeyPath = form.Items[6].Value

	b, err := json.MarshalIndent(clientConf, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(*clientConfigPath, b, 0755)
}
