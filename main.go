package main

import (
)


func main() {
	s := NewStatusLine()
	s.Start()
	go s.Run()
	s.Render()
}

