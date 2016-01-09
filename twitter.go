package main

import (
	"fmt"
	an "github.com/ChimeraCoder/anaconda"
	"net/url"
//	spew "github.com/davecgh/go-spew/spew"
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
//	return updateTimeline(api)
//	for i := 0; i<len(home); i++ {
//		fmt.Print(home[i].CreatedAt + "   ")
//		fmt.Print(home[i].User.ScreenName + "  ")
//		fmt.Println(home[i].Text)
//	}
//	spew.Dump(home)
//	tweetlist, _ := api.GetListTweets(1, true, v)
//	fmt.Println(tweetlist)
//	for i := 0; i < len(tweetlist); i++ {
//		fmt.Println("hello?")
//		fmt.Println(tweetlist[i].Text)
//	}
}

func updateTimeline(api *an.TwitterApi) []an.Tweet {
	v := url.Values{}
	home, err := api.GetHomeTimeline(v)
	if err != nil {
		fmt.Println("error: ", err)
	}
	return home
}