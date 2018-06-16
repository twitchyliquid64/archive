package main

import (
	"math"

	termbox "github.com/nsf/termbox-go"
)

type formItem struct {
	Name       string
	Type       string
	Value      string
	FgOverride termbox.Attribute
}

// FormUI represents a form the user can input
type FormUI struct {
	Title, Subtitle string
	Statustext      string
	Fg, Bg          termbox.Attribute
	Items           []formItem

	currentSelection int
	Verify           func([]formItem) bool
}

func (u *FormUI) draw() {
	if u.Verify != nil {
		u.Verify(u.Items)
	}

	const coldef = termbox.ColorDefault
	termbox.Clear(coldef, coldef)
	w, h := termbox.Size()

	startH := 3
	labelColumnSize := maxTextWidth(u.Items)
	for i, item := range u.Items {
		color := termbox.ColorCyan
		if item.FgOverride != 0 {
			color = item.FgOverride
		}
		termPrint(2, startH+i, color, termbox.ColorDefault, item.Name+":")
		if u.currentSelection == i {
			newPrint := padRight(item.Value, " ", w-5-labelColumnSize)
			termPrint(4+labelColumnSize, startH+i, termbox.ColorDefault|termbox.AttrUnderline, termbox.ColorDefault, newPrint)
			termbox.SetCursor(4+labelColumnSize+textWidth(item.Value), startH+i)
		} else {
			termPrint(4+labelColumnSize, startH+i, termbox.ColorDefault, termbox.ColorDefault, item.Value)
		}
	}

	termPrintRightAlign(w-1, h-5, termbox.ColorRed|termbox.AttrBold, termbox.ColorDefault, "┌"+padRight("", " ", 8)+"┐")
	termPrintRightAlign(w-1, h-4, termbox.ColorRed|termbox.AttrBold, termbox.ColorDefault, "│"+padRight("", " ", 8)+"│")
	termPrintRightAlign(w-4, h-4, termbox.ColorDefault|termbox.AttrBold, termbox.ColorDefault, "Done")
	termPrintRightAlign(w-1, h-3, termbox.ColorRed|termbox.AttrBold, termbox.ColorDefault, "└"+padRight("", " ", 8)+"┘")
	if len(u.Items) == u.currentSelection {
		termbox.SetCell(w-3, h-4, '⯇', termbox.ColorDefault|termbox.AttrBold, termbox.ColorDefault)
		termbox.SetCell(w-10, h-4, '⯈', termbox.ColorDefault|termbox.AttrBold, termbox.ColorDefault)
		termbox.HideCursor()
	}

	fill(1, 0, w, 1, termbox.Cell{Ch: ' ', Bg: u.Bg})
	termPrint(1, 0, termbox.ColorBlack, termbox.ColorYellow, u.Title)
	termPrintRightAlign(w, 0, u.Fg, u.Bg, u.Subtitle)
	fill(1, h-1, w, 1, termbox.Cell{Ch: ' ', Bg: u.Bg})
	statusText := "Please fill in the form."
	if u.Statustext != "" {
		statusText = u.Statustext
	}
	termPrint(1, h-1, termbox.ColorBlack, termbox.ColorYellow, statusText)
	termPrintRightAlign(w, h-1, u.Fg, u.Bg, "Exit with <ESC>")
	termbox.Flush()
}

func padRight(input, pad string, length int) string {
	for {
		input += pad
		if len(input) > length {
			return input[0:length]
		}
	}
}

func maxTextWidth(items []formItem) int {
	w := 0
	for _, item := range items {
		w = int(math.Max(float64(textWidth(item.Name)), float64(w)))
	}
	return w
}

// Run presents the form to the user and returns on error, abort or completion.
func (u *FormUI) Run() error {
	termbox.SetInputMode(termbox.InputEsc | termbox.InputMouse)
	defer termbox.HideCursor()
	w, h := termbox.Size()

	for {
		u.draw()
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventResize:
			w, h = termbox.Size()
		case termbox.EventMouse:
			if ev.Key == termbox.MouseLeft {
				if ev.MouseX >= (w-11) && ev.MouseX <= (w-1) && ev.MouseY >= (h-5) && ev.MouseY <= (h-3) { //process DONE being pressed
					if u.Verify != nil {
						if !u.Verify(u.Items) {
							continue
						}
					}
					return nil
				}
				if ev.MouseY >= 3 && ev.MouseY <= (len(u.Items)+2) { //otherwise, check if a field was pressed
					u.currentSelection = ev.MouseY - 3
				}
			}

		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEsc, termbox.KeyCtrlC:
				return errAbortMenu
			case termbox.KeyArrowDown:
				u.currentSelection++
				if u.currentSelection > len(u.Items) {
					u.currentSelection = 0
				}
			case termbox.KeyArrowUp:
				u.currentSelection--
				if u.currentSelection < 0 {
					u.currentSelection = len(u.Items)
				}
			case termbox.KeyBackspace, termbox.KeyBackspace2:
				v := u.Items[u.currentSelection].Value
				if u.currentSelection < len(u.Items) && len(v) > 0 {
					u.Items[u.currentSelection].Value = v[:len(v)-1]
				}
			case termbox.KeyEnter:
				if u.Verify != nil {
					if !u.Verify(u.Items) {
						continue
					}
				}
				if u.currentSelection == len(u.Items) {
					return nil
				}
			default:
				if ev.Ch != 0 && u.currentSelection < len(u.Items) { //typing event on a form field
					u.Items[u.currentSelection].Value += string(ev.Ch)
				}
			}
		case termbox.EventError:
			return ev.Err
		}
	}
}
