package main

import (
	termbox "github.com/nsf/termbox-go"
)

// SelectionItem represents an option in a selector menu
type SelectionItem struct {
	ID, Text string
}

// SelectionUI implements a root widget to allow the user to make a selection with their arrow keys.
type SelectionUI struct {
	currentSelection int
	Title, Subtitle  string
	Fg, Bg           termbox.Attribute
	Items            []SelectionItem
}

func (u *SelectionUI) draw() {
	const coldef = termbox.ColorDefault
	termbox.Clear(coldef, coldef)
	w, h := termbox.Size()

	startH := (h / 2) - (len(u.Items) / 2)
	for i, item := range u.Items {
		if u.currentSelection == i {
			fill(10, startH+i, w-20, 1, termbox.Cell{Ch: ' ', Bg: termbox.ColorCyan})
			termPrintCenterAlign(w/2, startH+i, termbox.ColorWhite, termbox.ColorCyan, item.Text)
		} else {
			termPrintCenterAlign(w/2, startH+i, termbox.ColorWhite, termbox.ColorDefault, item.Text)
		}
	}

	fill(1, 0, w, 1, termbox.Cell{Ch: ' ', Bg: u.Bg})
	termPrint(1, 0, termbox.ColorBlack, termbox.ColorYellow, u.Title)
	termPrintRightAlign(w, 0, u.Fg, u.Bg, u.Subtitle)
	fill(1, h-1, w, 1, termbox.Cell{Ch: ' ', Bg: u.Bg})
	termPrint(1, h-1, termbox.ColorBlack, termbox.ColorYellow, "Make your selection using the arrow keys and <ENTER>.")
	termPrintRightAlign(w, h-1, u.Fg, u.Bg, "Exit with <ESC>")
	termbox.Flush()
}

// Run presents the selection to the user, returning their selection.
func (u *SelectionUI) Run() (string, error) {
	termbox.SetInputMode(termbox.InputEsc)

	for {
		u.draw()
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEsc, termbox.KeyCtrlC:
				return "", errAbortMenu
			case termbox.KeyArrowUp:
				u.currentSelection--
				if u.currentSelection < 0 {
					u.currentSelection = len(u.Items) - 1
				}
			case termbox.KeyEnter:
				return u.Items[u.currentSelection].ID, nil
			case termbox.KeyArrowDown:
				u.currentSelection++
				if u.currentSelection >= len(u.Items) {
					u.currentSelection = 0
				}
			}
		case termbox.EventError:
			return "", ev.Err
		}
	}
}
