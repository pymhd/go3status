package main

import (
	"os"
	_ "fmt"
	"bufio"
	"encoding/json"
	"go3status/modules"
)

func RunClickEventsHandler() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		b := scanner.Bytes()
		ce := new(modules.ClickEvent)
		
		if json.Valid(b) {
			json.Unmarshal(b, ce)
		} else {
			json.Unmarshal(b[1:], ce)
		}
		//not strange stdin
		if len(ce.Name) > 0 {
			modules.Modules[ce.Name].HandleClickEvent(ce, cfg.Modules[ce.Name])
		} 
	}
}
