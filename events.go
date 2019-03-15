package main

import (
	"os"
	_ "fmt"
	"strconv"
	"bufio"
	"encoding/json"
	"go3status/modules"
)

func RunClickEventsHandler() {
	scanner := bufio.NewScanner(os.Stdin)
	//cache := make(map[string]int)
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
			instance := ce.Instance
			//cached_n, ok := cache[name]
			//if ok {
			//	mc := cfg.Modules[cached_n][name]
			//	mc.Id = cached_n
			//	modules.Modules[cached_n].HandleClickEvent(ce, mc)
			//	continue
			//}
			for n, _ := range cfg.Modules {
				if strconv.Itoa(n) == instance {
						//cache[name] = n
						mc := cfg.Modules[n][name]
						mc.Id = n
						modules.Modules[n].HandleClickEvent(ce, mc)
				}
			}
		} 
	}
}
