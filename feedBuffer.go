package main

type FeedBuffer struct {
	currSelected int
	currPrinted  int
	items        []FeedItem
	lastUpdate   string
}
