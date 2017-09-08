package main

import (
	"fmt"
	ui "github.com/gizak/termui"
	"github.com/mmcdole/gofeed"
)

func main() {
	err := ui.Init()
	if err != nil {
		panic(err)
	}
	defer ui.Close()

	ls := ui.NewList()

	fp := gofeed.NewParser()
	feed, err := fp.ParseURL("http://www.mbl.is/feeds/fp/")
	if err != nil {
		fmt.Println(err)
	}
	for i := 0; i < len(feed.Items); i++ {
		ls.Items = append(ls.Items, feed.Items[i].Title)
	}

	ls.ItemFgColor = ui.ColorYellow
	ls.BorderLabel = "List"
	ls.Height = 57

	ui.Body.AddRows(
		ui.NewRow(
			ui.NewCol(6, 0, ls),
		),
	)
	ui.Body.Align()
	ui.Render(ui.Body)
	ui.Handle("/sys/kbd/q", func(ui.Event) {
		ui.StopLoop()
	})
	ui.Loop()
}
