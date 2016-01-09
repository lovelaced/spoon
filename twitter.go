package main

import (
	"fmt"
	an "github.com/ChimeraCoder/anaconda"
	"net/url"
)

func createAPI(ck string, csk string, at string, ats string) {
	an.SetConsumerKey(ck)
	an.SetConsumerSecret(csk)
	api := an.NewTwitterApi(at, ats)
	fmt.Print("made it")
	v := url.Values{}
	home, err := api.GetHomeTimeline(v)
	if err != nil {
		fmt.Println("error: ", err)
	}
	fmt.Print(home)
//	tweetlist, _ := api.GetListTweets(1, true, v)
//	fmt.Println(tweetlist)
//	for i := 0; i < len(tweetlist); i++ {
//		fmt.Println("hello?")
//		fmt.Println(tweetlist[i].Text)
//	}
}
