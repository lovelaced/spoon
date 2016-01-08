package main

import (
	"bufio"
	"fmt"
	. "github.com/rthornton128/goncurses"
	"os"
	"time"
	"strings"
)

func main() {


	// check to see if the config file exists
	// if it doesn't, create it and ask for API keys


	if _, err := os.Stat("$HOME/.config/spoon/config"); err != nil {
		fmt.Println(`Hello there! It looks like this is your first time running spoon
			or you've misplaced your configuration file. We're gonna need those Twitter keys
			again - actually, do you even use Twitter?`)
	}
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("y/n: ")
	text, _ := reader.ReadString('\n')
	if strings.TrimSpace(text) == "y" {
		fmt.Println(`Alright, now hand over your keys.
			Consumer key first, then the secret consumer key.`)
		ck, _ := reader.ReadString('\n')
		csk, _ := reader.ReadString('\n')
		fmt.Println(`Okay, great. Now we'll also need your Access Token and
			Access Token Secret to fire up the reader.`)
		at, _ := reader.ReadString('\n')
		ats, _ := reader.ReadString('\n')
		createAPI(ck, csk, at, ats)
	}

	stdscr, err := Init()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	defer End()
	rows, cols := stdscr.MaxYX()
	stdscr.Print("feeds go here")
	stdscr.MovePrint(3, 0, "q to quit")
	stdscr.Refresh()
	window, _ := NewWindow(1, cols, rows-1, 0)
	mx, my := window.MaxYX()
	StartColor()
	InitPair(1, C_BLACK, C_YELLOW)
	window.ColorOn(int16(1))
	// TODO: change time format so it's better
	barinfo := time.Now().Format(time.RFC822)

	window.MovePrint(mx/2, my-len(barinfo)-2, barinfo)
	window.ColorOff(int16(1))
	bgc := ColorPair(int16(1))
	window.SetBackground(bgc)
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
