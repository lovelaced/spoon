package main

import (
	//"bufio"
	"fmt"
	"os"
        "io"
        "log"
	//"strings"
	"time"
        "encoding/json"
        "launchpad.net/go-xdg"
	"github.com/ChimeraCoder/anaconda"
	. "github.com/rthornton128/goncurses"
)

func main() {
	// check to see if the config file exists
	// if it doesn't, create it and ask for API keys

	var tweets []anaconda.Tweet
        conf := loadConfig()
        api, tweets = createAPI(conf.Twitter.Api_key, conf.Twitter.Api_secret, 
                                conf.Twitter.Access_token, conf.Twitter.Access_secret)

	if _, err := xdg.Config.Find("/spoon/config.json"); err != nil {
		fmt.Println(`Hello there! It looks like this is your first time running spoon
			or you've misplaced your configuration file. We're gonna need those Twitter keys
			again - actually, do you even use Twitter?`)
	}

	//reader := bufio.NewReader(os.Stdin)
	//fmt.Print("y/n: ")
	//text, _ := reader.ReadString('\n')
	//if strings.TrimSpace(text) == "y" {
	//	fmt.Println(`Alright, now hand over your keys.
	//		Consumer key first, then the secret consumer key.`)
	//	ck, _ := reader.ReadString('\n')
	//	csk, _ := reader.ReadString('\n')
	//	fmt.Println(`Okay, great. Now we'll also need your Access Token and
	//		Access Token Secret to fire up the reader.`)
	//	at, _ := reader.ReadString('\n')
	//	ats, _ := reader.ReadString('\n')
	//	ck = strings.TrimSpace(ck)
	//	csk = strings.TrimSpace(csk)
	//	at = strings.TrimSpace(at)
	//	ats = strings.TrimSpace(ats)
	//	api, tweets = createAPI(ck, csk, at, ats)
	//}

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
	InitPair(1, C_BLACK, C_YELLOW)
	InitPair(2, C_WHITE, C_BLACK)
	InitPair(3, C_BLUE, C_BLACK)
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
		updateTimeline(api)

		for i:=0; i<len(tweets); i++ {
			t, _ := time.Parse(time.RubyDate, tweets[i].CreatedAt)
			stdscr.ColorOn(2)
			stdscr.Print(t.Format("15:04:05") + " ")
			stdscr.ColorOff(2)
			stdscr.ColorOn(3)
			stdscr.AttrOn(A_BOLD)
			stdscr.Print(tweets[i].User.ScreenName)
			stdscr.AttrOff(A_BOLD)
			stdscr.ColorOff(3)
			stdscr.Print("  ")
			UseDefaultColors()
			stdscr.Println(tweets[i].Text)
		}

		nrows, ncols := stdscr.MaxYX()
		if nrows != mx || ncols != my {
//			goto redraw
		}
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

func loadConfig() Config{
    conf := Config{}
    var file io.Reader
    path, err := xdg.Config.Ensure("spoon/config.json")
    if err == nil {
        file, _ = os.Open(path)
        fmt.Println(path)
    } else {
        fmt.Println(err.Error)
    }
    decoder := json.NewDecoder(file)
    err = decoder.Decode(&conf)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(conf)
    return conf
}
