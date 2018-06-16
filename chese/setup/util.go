package main

import (
	"os"

	runewidth "github.com/mattn/go-runewidth"
	termbox "github.com/nsf/termbox-go"
)

func fileNotExists(path string) bool {
	_, err := os.Stat(path)
	return os.IsNotExist(err)
}

func fill(x, y, w, h int, cell termbox.Cell) {
	for ly := 0; ly < h; ly++ {
		for lx := 0; lx < w; lx++ {
			termbox.SetCell(x+lx, y+ly, cell.Ch, cell.Fg, cell.Bg)
		}
	}
}

func termPrint(x, y int, fg, bg termbox.Attribute, msg string) {
	for _, c := range msg {
		termbox.SetCell(x, y, c, fg, bg)
		x += runewidth.RuneWidth(c)
	}
}

func textWidth(text string) int {
	i := 0
	for _, c := range text {
		i += runewidth.RuneWidth(c)
	}
	return i
}

func termPrintRightAlign(x, y int, fg, bg termbox.Attribute, msg string) {
	for _, c := range msg {
		x -= runewidth.RuneWidth(c)
	}
	termPrint(x, y, fg, bg, msg)
}

func termPrintCenterAlign(x, y int, fg, bg termbox.Attribute, msg string) {
	i := 0
	for _, c := range msg {
		i += runewidth.RuneWidth(c)
	}
	termPrint(x-(i/2), y, fg, bg, msg)
}

func horizontalLine(y, width int, ch rune, fg, bg termbox.Attribute) {
	for i := 0; i < (width + 1); i++ {
		termbox.SetCell(i, y, ch, fg, bg)
	}
}
