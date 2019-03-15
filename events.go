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
	cache := make(map[string]int)
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
			name := ce.Name
			cached_n, ok := cache[name]
			if ok {
				modules.Modules[name].HandleClickEvent(ce, cfg.Modules[cached_n][name])
				continue
			}
			for n, modmap := range cfg.Modules {
				for k, _ := range modmap {
					if k == name {
						cache[name] = n
						modules.Modules[name].HandleClickEvent(ce, cfg.Modules[n][name])
					}
				}
			}
		} 
	}
}
