package main

import (
	"flag"
	"go3status/modules"
	
	"github.com/pymhd/go-logging"
	"github.com/pymhd/go-logging/handlers"
)

var (
	cfg *Config
	log *logger.Logger
)

func main() {
	// provide some additional flags...
	cf := flag.String("config", "~/go/src/go3status/config.yaml", "cofig file")
	flag.Parse()
	//Parse config file
	cfg = ParseConfig(*cf)
	
	if len(cfg.Global.LogFile) > 0 {
		severity := 3 - cfg.Global.LogLevel
		log = logger.New(handlers.NewFileHandler(cfg.Global.LogFile), severity, logger.OLEVEL|logger.OFILE|logger.OTIME)
	} else {
		log = logger.New(handlers.NullHandler{}, logger.ERROR, 0)
	}
        modules.RegisterLogger(log)

	s := NewStatusLine()
	//Start will create chans and select cases to handle, and run modules
	s.Start()
	//Run is main handler that updates receives updates from modules and pushes them to internal blocks slice.
	go s.Run()
	//This is stdin scanner for JSON formatted click events
	go RunClickEventsHandler(s)
	//This is method that prints internal blocks slice to stdout for i3bar
	s.Render()
}
