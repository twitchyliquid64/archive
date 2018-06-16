package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"

	termbox "github.com/nsf/termbox-go"
	"github.com/twitchyliquid64/chese/config"
)

func editServerConfig() error {
	var serverConf config.Server
	if !fileNotExists(*serverConfigPath) {
		d, err := ioutil.ReadFile(*serverConfigPath)
		if err != nil {
			return err
		}
		json.Unmarshal(d, &serverConf)
	}

	form := FormUI{
		Fg:       termbox.ColorBlack,
		Bg:       termbox.ColorYellow,
		Title:    "Chese Setup --> Configure Server",
		Subtitle: *serverConfigPath,
		Items: []formItem{
			formItem{
				Name:  "Listener",
				Value: ":8427",
			},
			formItem{
				Name:  "Server certificate file",
				Value: *servCertPath,
			},
			formItem{
				Name:  "Server key file",
				Value: *servKeyPath,
			},
			formItem{
				Name:  "CA certificate file",
				Value: *caCertPath,
			},
			formItem{
				Name:  "Resources path",
				Value: "/usr/share/chese-server/resources",
			},
			formItem{
				Name:  "Debug (true/false)",
				Value: "false",
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

		_, _, listenerErr := net.SplitHostPort(items[0].Value)
		if listenerErr != nil {
			isValid = false
			items[0].FgOverride = termbox.ColorRed
			form.Statustext = fmt.Sprintf("Invalid listener address (error %q)", listenerErr.Error())
		} else {
			items[0].FgOverride = 0
		}
		for i := 1; i < (len(items) - 1); i++ {
			if fileNotExists(items[i].Value) {
				items[i].FgOverride = termbox.ColorRed
				form.Statustext = fmt.Sprintf("%q does not exist!", items[i].Value)
			} else {
				items[i].FgOverride = 0
			}
		}

		if items[5].Value != "true" && items[5].Value != "false" {
			items[5].FgOverride = termbox.ColorRed
		} else {
			items[5].FgOverride = 0
		}

		return isValid
	}
	err := form.Run()
	if err != nil {
		return err
	}
	termbox.Close()

	serverConf.Listener = form.Items[0].Value
	serverConf.TLS.MainCertPath = form.Items[1].Value
	serverConf.TLS.MainKeyPath = form.Items[2].Value
	serverConf.TLS.CACertPath = form.Items[3].Value
	serverConf.ResourcePath = form.Items[4].Value
	serverConf.Debug = form.Items[5].Value == "true"

	b, err := json.MarshalIndent(serverConf, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(*serverConfigPath, b, 0755)
}
