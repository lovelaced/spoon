package main

type FeedBuffer struct {
	CurrSelected int
	CurrPrinted  int
	Items        []FeedItem
	LastUpdate   string
	Name         string
}
