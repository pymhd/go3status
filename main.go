package main

import (
	_ "flag"
)
var (
	cfg	*Config
)

func main() {
	
	cfg = ParseConfig("/home/mhd/go/src/go3status/config.yaml")
	s := NewStatusLine()
	s.Start()
	go s.Run()
	s.Render()
}
