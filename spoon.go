package main

import (
	//"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	//"strings"
	"github.com/ChimeraCoder/anaconda"
	. "github.com/rthornton128/goncurses"
	xdg "launchpad.net/go-xdg"
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

	stdscr, err := Init()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	defer End()

	stdscr.Refresh()
	rows, cols := stdscr.MaxYX()
	window, _ := NewWindow(1, cols, rows-1, 0)
	mx, my := window.MaxYX()

	StartColor()
	UseDefaultColors()
	InitPair(1, C_BLACK, C_YELLOW)
	InitPair(2, C_WHITE, -1)
	InitPair(3, C_BLUE, -1)
	InitPair(4, C_YELLOW, -1)
	window.ColorOn(int16(1))
	// TODO: change time format so it's better maybe
	barinfo := time.Now().Format(time.RFC822)

	window.MovePrint(mx/2, my-len(barinfo)-2, barinfo)
	window.ColorOff(int16(1))
	bgc := ColorPair(int16(1))
	window.SetBackground(bgc)
	NewPanel(window)
	//TODO: Get total number of feeds from config file
	totalFeeds := 3
	var feeds [3]*Panel
	m := make(map[string]*Window)
	var names [3]string
	names[0] = "all"
	names[1] = "twitter"
	names[2] = "rss"
	totalLength := 1

	for i := 0; i < len(names); i++ {
		if i == 0 {
			window.AttrOn(A_BOLD)
		}
		window.MovePrint(mx/2, totalLength+1, "["+names[i]+"] ")
		if i == 0 {
			window.AttrOff(A_BOLD)
		}
		totalLength += len(names[i]) + 3
	}
	var win *Window
	for i := totalFeeds - 1; i >= 0; i-- {
		win, _ = NewWindow(rows-1, cols, 0, 0)
		feeds[i] = NewPanel(win)
		m[names[i]] = win
	}
	feedList := make([]FeedItem, 1)
	go updateWindow(win, tweets, feedList)
	//go updateTwitterWindow(win, tweets, feedList)
	active := 0

main:
	for {
		UpdatePanels()
		Update()
		ch := win.GetChar()
		switch Key(ch) {
		case 'q':
			break main
		case KEY_TAB:
			active += 1
			if active > totalFeeds-1 {
				active = 0
			}
			feeds[active].Top()
			totalLength = 1
			for i := 0; i < len(names); i++ {
				if names[active] == names[i] {
					window.AttrOn(A_BOLD)
				}
				window.MovePrint(mx/2, totalLength+1, "["+names[i]+"] ")
				if names[active] == names[i] {
					window.AttrOff(A_BOLD)
				}
				totalLength += len(names[i]) + 3
			}
		case KEY_RIGHT:
		//TODO: if focus is on bbar, scrolls through feeds
		case KEY_LEFT:
		//TODO: same shit
		case KEY_UP:
		//TODO: scroll through feed or expanded feed
		//	win.Scroll(1)
		case KEY_DOWN:
		//TODO: same shit
		//	win.Scroll(-1)
		case KEY_ENTER:
			//TODO: expands feed if feed is focused, selects feed if bbar is focused

		}
	}
}

func updateWindow(win *Window, tweets []anaconda.Tweet, feedList []FeedItem) {
	for {
		UpdatePanels()
		tweets = updateTimeline(api)
		win.Erase()
		win.NoutRefresh()
		for i := 0; i < len(tweets); i++ {
			t, _ := time.Parse(time.RubyDate, tweets[i].CreatedAt)
			win.ColorOn(2)
			_, my := win.MaxYX()
			lineLength := 1
			if my > 40 {
				win.Print(t.Format(" 15:04") + " ")
				lineLength = 7
			}
			win.ColorOff(2)
			win.ColorOn(3)
			win.AttrOn(A_BOLD)
			padding := 8 - len(tweets[i].User.ScreenName)
			for i := 0; i < padding; i++ {
				win.Print(" ")
			}
			if len(tweets[i].User.ScreenName) > 8 {
				win.Print(tweets[i].User.ScreenName[:8])
				lineLength += 8
			} else {
				win.Print(tweets[i].User.ScreenName)
				lineLength += len(tweets[i].User.ScreenName)
			}
			win.AttrOff(A_BOLD)
			win.ColorOff(3)
			win.ColorOn(4)
			win.Print(" | ")
			lineLength += 3
			win.ColorOff(4)
			UseDefaultColors()
			text := strconv.QuoteToASCII(tweets[i].Text)
			var newFeed FeedItem
			newFeed.Body = text
			newFeed.Name = tweets[i].User.ScreenName
			newFeed.Read = false
			newFeed.Time = t
			feedList = append(feedList, newFeed)
			if len(text) > my-lineLength-7 {
				win.Println(text[:my-lineLength-7] + "...")
			} else {
				win.Println(text)
			}
		}
		win.NoutRefresh()
		Update()
		time.Sleep(10 * time.Second)
	}
}
