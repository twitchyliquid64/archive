package main

import (
	termbox "github.com/nsf/termbox-go"
)

func defaultString(value, def string) string {
	if value != "" {
		return value
	}
	return def
}

func mainMenu() error {
	for {
		ui := &SelectionUI{
			Fg:       termbox.ColorBlack,
			Bg:       termbox.ColorYellow,
			Title:    "Chese Setup",
			Subtitle: "Main Menu",
			Items:    []SelectionItem{},
		}

		if !fileNotExists(*serverConfigPath) || *serverConfigPath != "/etc/chese/server.json" {
			ui.Items = append(ui.Items, SelectionItem{
				ID:   "mint-server",
				Text: "Mint new CA/Server TLS Certificates",
			},
				SelectionItem{
					ID:   "issue-cert",
					Text: "Issue new Client Certificate",
				},
				SelectionItem{
					ID:   "edit-server-config",
					Text: "Edit Server Configuration",
				})
		}

		if !fileNotExists(*clientConfigPath) || *clientConfigPath != "DEFAULT" {
			ui.Items = append(ui.Items, SelectionItem{
				ID:   "edit-client-config",
				Text: "Edit Client Configuration",
			})
		}

		selection, err := ui.Run()
		if err != nil {
			return err
		}

		switch selection {
		case "edit-server-config":
			if err := editServerConfig(); err != errAbortMenu {
				return err
			}

		case "edit-client-config":
			if err := editClientConfig(); err != errAbortMenu {
				return err
			}

		case "mint-server":
			if err := makeServerCertsMenu(); err != errAbortMenu {
				return err
			}
		case "issue-cert":
			if err := issueClientCertsMenu(); err != errAbortMenu {
				return err
			}
		}
	}
}
