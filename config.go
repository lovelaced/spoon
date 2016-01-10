package main

import (
    "io"
    "launchpad.net/go-xdg"
    "os"
    "log"
    "encoding/json"
)

type Config struct {
	Twitter struct {
                Api_key      string
                Api_secret   string
                Access_token string
                Access_secret string
	}
}

func loadConfig() Config{
    conf := Config{}
    var file io.Reader
    path, err := xdg.Config.Ensure("spoon/config.json")
    if err == nil {
        file, _ = os.Open(path)
    } else {
        log.Fatal(err)
    }
    decoder := json.NewDecoder(file)
    err = decoder.Decode(&conf)
    if err != nil {
        log.Fatal(err)
    }
    return conf
}
