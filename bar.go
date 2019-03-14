package main

import (
	"fmt"
	"encoding/json"
	"go3status/modules"
	"os"
	"reflect"
	"sync"
	"time"
)


type StatusLine struct {
	sync.Mutex
	Header  string
	Refresh chan bool
	Modules []modules.Module
	Blocks  []modules.ModuleOutput
	cases   []reflect.SelectCase
}

func (sl *StatusLine) Start() {
	for n, module := range sl.Modules {
		c := make(chan modules.ModuleOutput)
		sl.cases[n] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(c)}
		go module.Run(c, cfg.Modules[module.Name()])
	}
}

func (sl *StatusLine) Run() {
	for {
		ch_num, value, _ := reflect.Select(sl.cases)

		mo, ok := value.Interface().(modules.ModuleOutput)
		if !ok {
			// why
			fmt.Println()
		}
		//Lock to update Statsuses field
		sl.Lock()
		sl.Blocks[ch_num] = mo
		sl.Unlock()
		// refresh panel if urgent, value alredy injected
		// use with caution, it is not rate-limited
		if mo.Refresh {
			//fmt.Println("Trying to refresh")
			sl.Refresh <- true
			//fmt.Println("PUSHED REFRESH")
		}
	}
}

func (sl *StatusLine) Render() {
	// ...
	fmt.Println(sl.Header)
	fmt.Printf("[[]\n,")

	enc := json.NewEncoder(os.Stdout)
	for {
		select {
		// regular output
		case <- time.After(cfg.Global.Interval):
			//fmt.Println("refresh by regulat time after")
			sl.render(enc)
		//need to refresh now
		case <- sl.Refresh:
			//fmt.Println("refresh by CHANEL")
			sl.render(enc)
		}
	}
}

func (sl *StatusLine) render(e *json.Encoder) {
	sl.Lock()
	defer sl.Unlock()

	e.Encode(sl.Blocks)
	fmt.Printf(",")
}

func NewStatusLine() *StatusLine {
	sl := new(StatusLine)
	sl.Header = `{"version": 1, "click_events": true, "stop_signal": 20}`
	for _, mod := range modules.Modules {
		sl.Modules = append(sl.Modules, mod)
	}
	sl.Blocks = make([]modules.ModuleOutput, 1)
	sl.cases = make([]reflect.SelectCase, 1)
	sl.Refresh = make(chan bool, 0)
	return sl
}
