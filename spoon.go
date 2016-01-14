package main

import (
	//"bufio"
	"fmt"
	"log"
	"os"
	//"strconv"
	//"strings"
	"github.com/ChimeraCoder/anaconda"
	. "github.com/gbin/goncurses"
	xdg "launchpad.net/go-xdg"
	"strings"
	"time"
)

func main() {
	// check to see if the config file exists
	// if it doesn't, create it and ask for API keys

	var tweets []anaconda.Tweet
	conf := loadConfig()
	err := initData()
	if err != nil {
		log.Fatal(err)
	}

	if _, err := xdg.Config.Find("/spoon/config.json"); err != nil {
		fmt.Println(`Hello there! It looks like this is your first time running spoon
			or you've misplaced your configuration file. We're gonna need those Twitter keys
			again - actually, do you even use Twitter?`)
	}
	// Init ncurses, standard window
	stdscr, err := Init()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	defer End()
	stdscr.Refresh()
	cols, rows := stdscr.MaxYX()
	bbar, _ := NewWindow(1, rows, cols-1, 0)
	my, mx := bbar.MaxYX()

	// Initialize colors
	StartColor()
	UseDefaultColors()
	InitPair(1, C_BLACK, C_YELLOW)
	InitPair(2, C_BLACK, C_WHITE)
	InitPair(3, C_BLUE, -1)
	InitPair(4, C_YELLOW, -1)
	bbar.ColorOn(int16(1))
	// TODO: change time format so it's better maybe

	// Set up lower right bar information
	barinfo := time.Now().Format(time.RFC822)

	bbar.MovePrint(my/2, mx-len(barinfo)-2, barinfo)
	bbar.ColorOff(int16(1))
	bgc := ColorPair(int16(1))
	bbar.SetBackground(bgc)
	Echo(false)
	NewPanel(bbar)
	// set up a different panel for each feed
	//TODO: Get total number of feeds from config file
	//TODO: Get names of feeds from config file
	//TODO: Set them as feedBuffer.Name
	totalFeeds := 3
	feedPanels := make([]*Panel, totalFeeds)
	var names [3]string
	names[0] = "all"
	names[1] = "twitter"
	names[2] = "rss"
	totalLength := 1

	// print the names of the feeds into the lower left of the bar
	for i := 0; i < len(names); i++ {
		if i == 0 {
			bbar.AttrOn(A_BOLD)
		}
		bbar.MovePrint(my/2, totalLength+1, "["+names[i]+"] ")
		if i == 0 {
			bbar.AttrOff(A_BOLD)
		}
		totalLength += len(names[i]) + 3
	}

	// this sets up a hashmap of names:windows but I'm not sure we need it
	var win *Window
	for i := totalFeeds - 1; i >= 0; i-- {
		win, _ = NewWindow(cols-1, rows, 0, 0)
		feedPanels[i] = NewPanel(win)
	}
	win.Keypad(true)

	//create a list which holds feed items, we need one for each feedBuffer
	feedList := make([]FeedItem, 1)

	//TODO: Create one feedbuffer for each feed and populate allBuffers with them

	//allBuffers := make([]FeedBuffer, totalFeeds)

	var feedBuffer FeedBuffer
	feedBuffer.CurrSelected = 0
	feedBuffer.CurrPrinted = 0
	feedBuffer.Items = feedList

	api, tweets = createAPI(conf.Twitter.Api_key, conf.Twitter.Api_secret,
		conf.Twitter.Access_token, conf.Twitter.Access_secret)
	//create channel to make sure goroutines shut down on exit
	endSignal := make(chan bool)
	//buffChan := make(chan FeedBuffer)
	go updateWindow(win, tweets, &feedBuffer, endSignal)
	//go updateTwitterWindow(win, tweets, feedList)
	active := 0

	for {
		UpdatePanels()
		Update()

		ch := win.GetChar()

		// check for user input
		switch Key(ch) {
		case 'q':
			endSignal <- true
		case KEY_TAB:
			active += 1
			if active > totalFeeds-1 {
				active = 0
			}
			feedPanels[active].Top()
			totalLength = 1
			for i := 0; i < len(names); i++ {
				if names[active] == names[i] {
					bbar.AttrOn(A_BOLD)
				}
				bbar.MovePrint(my/2, totalLength+1, "["+names[i]+"] ")
				if names[active] == names[i] {
					bbar.AttrOff(A_BOLD)
				}
				totalLength += len(names[i]) + 3
			}
		case ':':
			//TODO: fix this
			_, mx = bbar.MaxYX()
			Echo(true)
			bbar.MovePrint(my, 1, len(feedBuffer.Items))
			bbar.Move(1, mx)
			bbar.ColorOn(1)
			//	window.Move(1, my)
			command, _ := bbar.GetString(mx)
			parseCommand(command, bbar)
			Echo(false)

		case KEY_RIGHT:
			//TODO: if focus is on bbar, scrolls through feeds
		case KEY_LEFT:
			//TODO: same shit
		case KEY_UP:
			if feedBuffer.CurrSelected == 0 {
				feedBuffer.CurrSelected = len(feedBuffer.Items) - 1
			} else {
				feedBuffer.CurrSelected -= 1
			}
			//	buffChan <- feedBuffer
		case KEY_DOWN:
			if feedBuffer.CurrSelected == len(feedBuffer.Items) {
				feedBuffer.CurrSelected = 0
			} else {
				feedBuffer.CurrSelected += 1
			}
			//		buffChan <- feedBuffer
		case KEY_ENTER:
			//TODO: expands feed if feed is focused, selects feed if bbar is focused

		}
	}
}

func updateWindow(win *Window, tweets []anaconda.Tweet, feedBuffer *FeedBuffer,
	endSignal chan bool) {

	go func(chan bool) {
		for {
			shutdown := <-endSignal
			if shutdown {
				End()
				os.Exit(0)
			}
		}
	}(endSignal)

	for {
		UpdatePanels()
		win.Erase()
		win.NoutRefresh()
		//TODO: if not Twitter...
		//TODO: make this nonblocking
		if feedBuffer.LastUpdate != "" {
			tweets = updateTimeline(api, feedBuffer.LastUpdate)
		}
		processTweets(feedBuffer, tweets)
		//feedBuffer = updateTwitterWindow(win, tweets, feedBuffer)
		printFeed(win, feedBuffer)
		win.NoutRefresh()
		Update()
	}
}

func printFeed(win *Window, feedBuffer *FeedBuffer) {

	my, mx := win.MaxYX()
	feedList := feedBuffer.Items
	// check to see if the feedlist is more than the number
	// of lines we have in our window
	iterations := my - 1
	if (feedBuffer.CurrPrinted)+iterations > len(feedBuffer.Items) {
		iterations = len(feedList) - feedBuffer.CurrPrinted
	}
	for i := feedBuffer.CurrPrinted; i < iterations; i++ {
		_, my := win.MaxYX()

		lineLength := 1
		if feedList[i].Time.IsZero() {
			continue
		}
		if i == feedBuffer.CurrSelected {
			win.AttrOn(A_REVERSE)
		}
		// only print time if the window is wider than 40 cols
		if mx > 40 {
			win.Print(feedList[i].Time.Format(" 15:04") + " ")
			lineLength = 7
		}
		win.ColorOn(3)
		win.AttrOn(A_BOLD)
		padding := 8 - len(feedList[i].Name)
		for i := 0; i < padding; i++ {
			win.Print(" ")
		}
		if len(feedList[i].Name) > 8 {
			win.Print(feedList[i].Name[:8])
		} else {
			win.Print(feedList[i].Name)
		}
		win.AttrOff(A_BOLD)
		win.ColorOff(3)
		win.ColorOn(4)
		win.Print(" │ ")
		lineLength += 8
		win.ColorOff(4)
		text := feedList[i].Body
		if len(text) > my-lineLength-7 {
			win.Println(strings.TrimSpace(text[:my-lineLength-7]) + "...")
		} else {
			win.Println(text)
		}
		if i == feedBuffer.CurrSelected {
			win.AttrOff(A_REVERSE)
		}
		win.ColorOn(1)
		win.MovePrint(my-3, mx/2, len(feedBuffer.Items))
		win.ColorOff(1)
		feedBuffer.CurrPrinted++

	}
}

func parseCommand(command string, window *Window) {
	switch command {
	case "/clear":
		window.Clear()
		window.NoutRefresh()
	case "/exit":
		End()
		os.Exit(0)
	}
}
