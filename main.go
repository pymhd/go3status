package main

import (
	_ "flag"
)
var (
	cfg	*Config
)

func main() {
	cfg = ParseConfig("./config.yaml")
	s := NewStatusLine()
	s.Start()
	go s.Run()
	s.Render()
}
