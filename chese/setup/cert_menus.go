package main

import (
	"fmt"
	"path"
	"strconv"
	"time"

	termbox "github.com/nsf/termbox-go"
	"github.com/twitchyliquid64/chese/tls"
)

func makeServerCertsMenu() error {
	form := FormUI{
		Fg:         termbox.ColorBlack,
		Bg:         termbox.ColorYellow,
		Title:      "Chese Setup --> Mint new CA",
		Subtitle:   "Form",
		Statustext: "Will generate into specified paths.",
		Items: []formItem{
			formItem{
				Name:  "CA Certificate path",
				Value: *caCertPath,
			},
			formItem{
				Name:  "CA Key path",
				Value: *caKeyPath,
			},
			formItem{
				Name:  "Server Certificate path",
				Value: *servCertPath,
			},
			formItem{
				Name:  "Server Key path",
				Value: *servKeyPath,
			},
			formItem{
				Name:  "RSA key size",
				Value: "2048",
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

		_, numErr := strconv.Atoi(items[4].Value)
		if numErr != nil {
			items[4].FgOverride = termbox.ColorRed
			isValid = false
		} else {
			items[4].FgOverride = 0
		}

		if items[0].Value != "" && items[0].Value != "." && fileNotExists(path.Dir(items[0].Value)) {
			form.Statustext = fmt.Sprintf("Directory %q does not exist!", path.Dir(items[0].Value))
		} else if items[1].Value != "" && items[1].Value != "." && fileNotExists(path.Dir(items[1].Value)) {
			form.Statustext = fmt.Sprintf("Directory %q does not exist!", path.Dir(items[1].Value))
		} else if items[2].Value != "" && items[2].Value != "." && fileNotExists(path.Dir(items[2].Value)) {
			form.Statustext = fmt.Sprintf("Directory %q does not exist!", path.Dir(items[2].Value))
		} else if items[3].Value != "" && items[3].Value != "." && fileNotExists(path.Dir(items[3].Value)) {
			form.Statustext = fmt.Sprintf("Directory %q does not exist!", path.Dir(items[3].Value))
		}

		return isValid
	}
	err := form.Run()
	if err != nil {
		return err
	}
	termbox.Close()
	fmt.Printf("Creating certificates / keys...")
	rsaSize, _ := strconv.Atoi(form.Items[4].Value)
	err = tls.MakeServerCert(form.Items[2].Value, form.Items[3].Value, form.Items[0].Value, form.Items[1].Value, rsaSize)
	if err != nil {
		fmt.Println("ERROR.")
		return err
	}
	fmt.Println("DONE.")

	return nil
}

func issueClientCertsMenu() error {
	form := FormUI{
		Fg:         termbox.ColorBlack,
		Bg:         termbox.ColorYellow,
		Title:      "Chese Setup --> Issue Client Certificate",
		Subtitle:   "Form",
		Statustext: "Will generate into client paths.",
		Items: []formItem{
			formItem{
				Name:  "CA Certificate path",
				Value: *caCertPath,
			},
			formItem{
				Name:  "CA Key path",
				Value: *caKeyPath,
			},
			formItem{
				Name:  "Client Certificate path",
				Value: *clientCertPath,
			},
			formItem{
				Name:  "Client Key path",
				Value: *clientKeyPath,
			},
			formItem{
				Name:  "RSA key size",
				Value: "2048",
			},
			formItem{
				Name:  "Validity age (months)",
				Value: "12",
			},
			formItem{
				Name:  "Validity age (days)",
				Value: "0",
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

		_, numErr := strconv.Atoi(items[4].Value)
		if numErr != nil {
			items[4].FgOverride = termbox.ColorRed
			isValid = false
		} else {
			items[4].FgOverride = 0
		}
		_, numErr = strconv.Atoi(items[5].Value)
		if numErr != nil {
			items[5].FgOverride = termbox.ColorRed
			isValid = false
		} else {
			items[5].FgOverride = 0
		}
		_, numErr = strconv.Atoi(items[6].Value)
		if numErr != nil {
			items[6].FgOverride = termbox.ColorRed
			isValid = false
		} else {
			items[6].FgOverride = 0
		}

		if items[2].Value != "" && items[2].Value != "." && fileNotExists(path.Dir(items[2].Value)) {
			form.Statustext = fmt.Sprintf("Directory %q does not exist!", path.Dir(items[2].Value))
		} else if items[3].Value != "" && items[3].Value != "." && fileNotExists(path.Dir(items[3].Value)) {
			form.Statustext = fmt.Sprintf("Directory %q does not exist!", path.Dir(items[3].Value))
		}

		if items[0].Value != "" && fileNotExists(items[0].Value) {
			items[0].FgOverride = termbox.ColorRed
			isValid = false
			form.Statustext = fmt.Sprintf("CA Cert %q does not exist!", items[0].Value)
		}
		if items[1].Value != "" && fileNotExists(items[1].Value) {
			items[1].FgOverride = termbox.ColorRed
			isValid = false
			form.Statustext = fmt.Sprintf("CA Key %q does not exist!", items[1].Value)
		}

		return isValid
	}
	err := form.Run()
	if err != nil {
		return err
	}
	termbox.Close()
	fmt.Printf("Minting client cert/key...")
	rsaSize, _ := strconv.Atoi(form.Items[4].Value)
	months, _ := strconv.Atoi(form.Items[5].Value)
	days, _ := strconv.Atoi(form.Items[6].Value)
	err = tls.IssueClientCert(form.Items[0].Value, form.Items[1].Value, form.Items[2].Value, form.Items[3].Value, rsaSize, time.Now().AddDate(0, months, days))
	if err != nil {
		fmt.Println("ERROR.")
		return err
	}
	fmt.Println("DONE.")

	return nil
}
