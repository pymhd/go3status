package main

import (
	_ "flag"
)

var (
	cfg *Config
)

func main() {

	cfg = ParseConfig("/home/kgs/go/src/go3status/config.yaml")
	s := NewStatusLine()
	s.Start()
	go s.Run()
	go RunClickEventsHandler()
	s.Render()
}
