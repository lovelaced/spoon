package main

import (
    "fmt"
    "os"
    . "github.com/rthornton128/goncurses"
)

func main() {
	stdscr, err := Init()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	defer End()
	rows, cols := stdscr.MaxYX()
	title := "muh panels"
	stdscr.Print("sup this is text")
	stdscr.MovePrint(3, 0, "q to quit")
//	stdscr.Refresh()
	window, _ := NewWindow(3, cols, rows-4, 0)
	window.Box(0, 0)
	mx, my := window.MaxYX()
	window.MovePrint(mx/2, my/2, title)
	StartColor()
	InitPair(1, C_BLACK, C_CYAN)
//	window.Color(int16(1))
	NewPanel(window)

	main:
		for {
			UpdatePanels()
			Update()
			ch := stdscr.GetChar()
			switch Key(ch) {
			case 'q':
				break main
			case KEY_TAB:
				// rotate focus between feed, expanded feed (if present), and bbar
			case KEY_RIGHT:
				// if focus is on bbar, scrolls through feeds
			case KEY_LEFT:
				// same shit
			case KEY_UP:
				// scroll through feed or expanded feed
			case KEY_DOWN:
				// same shit
			case KEY_ENTER:
				// expands feed if feed is focused, selects feed if bbar is focused

			}
		}
}
