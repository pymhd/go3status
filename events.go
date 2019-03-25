package main

import (
        "os"
	"bufio"
	"strconv"
	"encoding/json"
	"go3status/modules"
)

func RunClickEventsHandler(sl *StatusLine) {
	scanner := bufio.NewScanner(os.Stdin)
	//cache := make(map[string]int)
	for scanner.Scan() {
		b := scanner.Bytes()
		ce := new(modules.ClickEvent)

		if json.Valid(b) {
			json.Unmarshal(b, ce)
		} else {
			if err := json.Unmarshal(b[1:], ce); err != nil {
				// skip, cant do nothing with
				continue
			}
		}
		//not strange stdin
		if len(ce.Name) > 0 {
			id, _ := strconv.Atoi(ce.Instance)
			go sl.Modules[id].HandleClickEvent(ce)
		}
	}
}
