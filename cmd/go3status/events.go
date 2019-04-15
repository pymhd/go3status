package main

import (
	"bufio"
	"encoding/json"
	"go3status/internal/modules"
	"os"
	"strconv"
)

func RunClickEventsHandler(sl *StatusLine) {
	scanner := bufio.NewScanner(os.Stdin)
	//cache := make(map[string]int)
	for scanner.Scan() {
		b := scanner.Bytes()
		ce := new(modules.ClickEvent)
		if json.Valid(b) {
			log.Info("New valid Click Event received")
			json.Unmarshal(b, ce)
		} else {
			if err := json.Unmarshal(b[1:], ce); err != nil {
				// skip, cant do nothing with
				log.Errorf("Could not unmarshal JSON obj (%s)\n", err)
				continue
			}
			log.Warning("Parsed json with extra char")
		}
		//not strange stdin
		if len(ce.Name) > 0 {
			id, _ := strconv.Atoi(ce.Instance)
			go sl.Modules[id].HandleClickEvent(ce)
			log.Infof("Ran %q module click event handler routine\n", sl.Modules[id].Name)
		}
	}
}
