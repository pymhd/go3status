package main

import (
	"os"
	_ "fmt"
	"bufio"
	"encoding/json"
)

type ClickEvent struct {
	Name     string   `json:"name"`
	Instance string   `json:"instance"`
	Button   int      `json:"button"`
	X        int      `json:"x"`
	Y        int      `json:"y"`
	Mod      []string `json:"modifiers"`
}


func RunClickEventsHandler() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		b := scanner.Bytes()
		ce := new(ClickEvent)
		
		f, _  := os.OpenFile("/tmp/i3ce", os.O_APPEND|os.O_WRONLY, 0644)
		enc := json.NewEncoder(f)
		if json.Valid(b) {
			json.Unmarshal(b, ce)
			enc.Encode(b)
		} else {
			json.Unmarshal(b[1:], ce)
		}
		enc.Encode(ce)
		f.Close()
	}
}
