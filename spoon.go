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
	api, tweets = createAPI(conf.Twitter.Api_key, conf.Twitter.Api_secret,
		conf.Twitter.Access_token, conf.Twitter.Access_secret)

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
	window, _ := NewWindow(1, rows, cols-1, 0)
	my, mx := window.MaxYX()

	// Initialize colors
	StartColor()
	UseDefaultColors()
	InitPair(1, C_BLACK, C_YELLOW)
	InitPair(2, C_BLACK, C_WHITE)
	InitPair(3, C_BLUE, -1)
	InitPair(4, C_YELLOW, -1)
	window.ColorOn(int16(1))
	// TODO: change time format so it's better maybe

	// Set up lower right bar information
	barinfo := time.Now().Format(time.RFC822)

	window.MovePrint(my/2, mx-len(barinfo)-2, barinfo)
	window.ColorOff(int16(1))
	bgc := ColorPair(int16(1))
	window.SetBackground(bgc)
	Echo(false)
	NewPanel(window)

	// set up a different panel for each feed
	//TODO: Get total number of feeds from config file
	//TODO: Get names of feeds from config file
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
			window.AttrOn(A_BOLD)
		}
		window.MovePrint(my/2, totalLength+1, "["+names[i]+"] ")
		if i == 0 {
			window.AttrOff(A_BOLD)
		}
		totalLength += len(names[i]) + 3
	}

	// this sets up a hashmap of names:windows but I'm not sure we need it
	var win *Window
	m := make(map[string]*Window)
	for i := totalFeeds - 1; i >= 0; i-- {
		win, _ = NewWindow(cols-1, rows, 0, 0)
		feedPanels[i] = NewPanel(win)
		m[names[i]] = win
	}

	//create a list which holds feed items, we need one for each feedBuffer
	feedList := make([]FeedItem, 1)
	var feedBuffer FeedBuffer
	feedBuffer.currSelected = 0
	feedBuffer.currPrinted = 0
	feedBuffer.items = feedList
	go updateWindow(win, tweets, feedBuffer)
	//go updateTwitterWindow(win, tweets, feedList)
	active := 0

main:
	for {
		UpdatePanels()
		Update()
		ch := win.GetChar()

		// check for user input
		switch Key(ch) {
		case 'q':
			break main
		case KEY_TAB:
			active += 1
			if active > totalFeeds-1 {
				active = 0
			}
			feedPanels[active].Top()
			totalLength = 1
			for i := 0; i < len(names); i++ {
				if names[active] == names[i] {
					window.AttrOn(A_BOLD)
				}
				window.MovePrint(my/2, totalLength+1, "["+names[i]+"] ")
				if names[active] == names[i] {
					window.AttrOff(A_BOLD)
				}
				totalLength += len(names[i]) + 3
			}
		case ':':
			//TODO: fix this
			_, mx = window.MaxYX()
			Echo(true)
			window.MovePrint(my, 1, len(feedBuffer.items))
			window.Move(1, mx)
			window.ColorOn(1)
			//	window.Move(1, my)
			command, _ := window.GetString(mx)
			parseCommand(command, window)
			Echo(false)

		case KEY_RIGHT:
		//TODO: if focus is on bbar, scrolls through feeds
		case KEY_LEFT:
		//TODO: same shit
		case KEY_UP:
			//TODO: scroll through feed or expanded feed
			feedBuffer.currSelected -= 1
			win.ScrollOk(true)
			printFeed(win, feedBuffer)
		case KEY_DOWN:
			//TODO: same shit
			win.ScrollOk(true)
			feedBuffer.currSelected += 1
			printFeed(win, feedBuffer)
		case KEY_ENTER:
			//TODO: expands feed if feed is focused, selects feed if bbar is focused

		}
	}
}

func updateWindow(win *Window, tweets []anaconda.Tweet, feedBuffer FeedBuffer) {
	for {
		feedList := feedBuffer.items
		UpdatePanels()
		win.Erase()
		win.NoutRefresh()
		//TODO: if Twitter...
		tweets = updateTimeline(api)
		feedBuffer.items = processTweets(feedList, tweets)
		//feedBuffer = updateTwitterWindow(win, tweets, feedBuffer)
		printFeed(win, feedBuffer)
		win.NoutRefresh()
		Update()
		time.Sleep(10 * time.Second)
	}
}

func printFeed(win *Window, feedBuffer FeedBuffer) {
	my, mx := win.MaxYX()
	feedList := feedBuffer.items
	// check to see if the feedlist is more than the number
	// of lines we have in our window
	iterations := my - 1
	if (feedBuffer.currPrinted)+iterations > len(feedBuffer.items)-1 {
		iterations = len(feedList) - feedBuffer.currPrinted - 1
	}
	for i := feedBuffer.currPrinted; i < iterations; i++ {
		_, my := win.MaxYX()
		lineLength := 1
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
			lineLength += 8
		} else {
			win.Print(feedList[i].Name)
			lineLength += len(feedList[i].Name)
		}
		win.AttrOff(A_BOLD)
		win.ColorOff(3)
		win.ColorOn(4)
		win.Print(" │ ")
		lineLength += 3
		win.ColorOff(4)
		UseDefaultColors()
		if i == feedBuffer.currSelected {
			win.ColorOn(2)
		}
		text := feedList[i].Body
		if len(text) > my-lineLength-7 {

			win.Println(strings.TrimSpace(text[:my-lineLength-7]) + "...")
		} else {
			win.Println(text)
		}
		if i == feedBuffer.currSelected {
			win.ColorOff(2)
		}
		feedBuffer.currPrinted++

	}
}

func parseCommand(command string, window *Window) {
	switch command {
	case "/clear":
		window.Clear()
		window.Refresh()
	case "/exit":
		End()
		os.Exit(0)
	}
}
