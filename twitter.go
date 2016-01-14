package main

import (
	an "github.com/ChimeraCoder/anaconda"
	. "github.com/gbin/goncurses"
	"log"
	"net/url"
	"strconv"
	"strings"
	"time"
)

var api *an.TwitterApi

func createAPI(ck string, csk string, at string, ats string) (*an.TwitterApi, []an.Tweet) {
	an.SetConsumerKey(ck)
	an.SetConsumerSecret(csk)
	api := an.NewTwitterApi(at, ats)
	v := url.Values{}
	v.Set("count", "200")
	home, err := api.GetHomeTimeline(v)
	if err != nil {
		log.Fatal(err)
	}
	return api, home
}

func updateTimeline(api *an.TwitterApi, id string) []an.Tweet {
	v := url.Values{}
	v.Set("since_id", id)
	home, err := api.GetHomeTimeline(v)
	if err != nil {
		log.Fatal(err)
	}
	time.Sleep(5 * time.Second)
	return home
}

func processTweets(feedBuffer *FeedBuffer, tweets []an.Tweet) {
	var extraSlice []FeedItem
	feedList := feedBuffer.Items
	for i := 0; i < len(tweets); i++ {
		if i == 0 {
			feedBuffer.LastUpdate = strconv.FormatInt(tweets[0].Id, 10)
		}
		t, _ := time.Parse(time.RubyDate, tweets[i].CreatedAt)
		text := strings.Replace(tweets[i].Text, "\n", " ", -1)
		var newFeed FeedItem
		newFeed.Body = text
		newFeed.Name = tweets[i].User.ScreenName
		newFeed.Read = false
		newFeed.Time = t

		extraSlice = append(extraSlice, newFeed)
	}
	feedList = append(extraSlice, feedList...)
	feedBuffer.Items = feedList
}

//TODO: Actually move this here maybe?
func updateTwitterWindow(win *Window, tweets []an.Tweet, feedBuffer FeedBuffer) {
}
