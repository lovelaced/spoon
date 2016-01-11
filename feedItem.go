package main

import (
	"time"
)

type FeedItem struct {
	Time time.Time
	Name string
	Body string
	Read bool
}
