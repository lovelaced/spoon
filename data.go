package main

import (
	"launchpad.net/go-xdg"
)

func initData() error {
	_, err := xdg.Data.Ensure("/spoon/")
	return err
}
