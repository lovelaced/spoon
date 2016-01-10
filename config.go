package main

import (
	"encoding/json"
	"io"
	"launchpad.net/go-xdg"
	"log"
	"os"
)

type Config struct {
	Twitter struct {
		Api_key       string
		Api_secret    string
		Access_token  string
		Access_secret string
	}

	RSS struct {
		Feeds []RSSfeed
	}
}

//TODO: Flesh out RSS feed options.
type RSSfeed struct {
	Name string
	Url  string
}

func loadConfig() Config {
	conf := Config{}
	var file io.Reader
	path, err := xdg.Config.Find("spoon/config.json")
	if err == nil {
		file, _ = os.Open(path)
	} else {
                os.
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&conf)
	if err != nil {
		log.Fatal(err)
	}
	return conf
}
