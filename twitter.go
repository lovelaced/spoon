package main

import (
	"fmt"
	an "github.com/ChimeraCoder/anaconda"
	. "github.com/gbin/goncurses"
	"net/url"
	"time"
	//	spew "github.com/davecgh/go-spew/spew"
	"strings"
)

//
//struct tweet {
//
//}

var api *an.TwitterApi

func createAPI(ck string, csk string, at string, ats string) (*an.TwitterApi, []an.Tweet) {
	an.SetConsumerKey(ck)
	an.SetConsumerSecret(csk)
	api := an.NewTwitterApi(at, ats)
	v := url.Values{}
	home, err := api.GetHomeTimeline(v)
	if err != nil {
		fmt.Println("error: ", err)
	}
	return api, home
}

func updateTimeline(api *an.TwitterApi) []an.Tweet {
	v := url.Values{}
	home, err := api.GetHomeTimeline(v)
	if err != nil {
		fmt.Println("error: ", err)
	}
	return home
}

func processTweets(feedList []FeedItem, tweets []an.Tweet) []FeedItem {

	for i := 0; i < len(tweets); i++ {
		t, _ := time.Parse(time.RubyDate, tweets[i].CreatedAt)
		if len(feedList) > 1 && feedList[len(feedList)-1].Time.Before(t) {
			continue
		}
		text := tweets[i].Text
		var newFeed FeedItem
		newFeed.Body = text
		newFeed.Name = tweets[i].User.ScreenName
		newFeed.Read = false
		newFeed.Time = t
		feedList = append(feedList, newFeed)
	}
	return feedList
}

//TODO: Actually move this here
func updateTwitterWindow(win *Window, tweets []an.Tweet, feedBuffer FeedBuffer) FeedBuffer {
	feedList := feedBuffer.items
	for {
		UpdatePanels()
		tweets = updateTimeline(api)
		win.Erase()
		win.NoutRefresh()
		for i := 0; i < len(tweets); i++ {
			t, _ := time.Parse(time.RubyDate, tweets[i].CreatedAt)
			//			if len(feedList) > 1 && feedList[len(feedList)-1].Time.Before(t) {
			//				continue
			//			}
			_, my := win.MaxYX()
			lineLength := 1
			if my > 40 {
				win.Print(t.Format(" 15:04") + " ")
				lineLength = 7
			}
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
			text := strings.TrimSpace(tweets[i].Text)
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
		feedBuffer.items = feedList
		feedBuffer.currPrinted += len(tweets)

		win.NoutRefresh()
		Update()
		time.Sleep(10 * time.Second)
		return feedBuffer
	}
}
