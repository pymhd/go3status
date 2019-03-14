package main

import (
	"flag"
)

var (
	cfg *Config
)

func main() {
	// provide some additional flags...
        cf := flag.String("config", "~/go/src/go3status/config.yaml", "cofig file")
	flag.Parse()
	//Parse config file
	cfg = ParseConfig(*cf)
	
	s := NewStatusLine()
	//Start will create chans and select cases to handle, and run modules
	s.Start()
	//Run is main handler that updates receives updates from modules and pushes them to internal blocks slice. 
	go s.Run()
	//This is stdin scanner for JSON formatted click events
	go RunClickEventsHandler()
	//This is method that prints internal blocks slice to stdout for i3bar
	s.Render()
}
