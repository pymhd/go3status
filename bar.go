package main

import (
	"os"
	"fmt"
	"time"
	"sync"
	"reflect"
	"encoding/json"
	"go3status/modules"
)

type ClickEvent struct {
	Name     string `json:"name"`
	Instance string `json:"instance"`
	Button   int    `json:"button"`
	XCoord   int    `json:"x"`
	YCoord   int    `json:"y"`
}


type StatusLine struct {
	sync.Mutex
	Header  string
	Refresh	chan bool
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
		
		msg  := value.Interface().(modules.ModuleOutput)

		//Lock to update Statsuses field
		sl.Lock()
		sl.Blocks[ch_num] = msg
		sl.Unlock()
	}
}

func (sl *StatusLine) Render() {
	// ...
	fmt.Println(sl.Header)
	fmt.Printf("[[]\n,")
	
	ticker := time.NewTicker(cfg.Global.Interval)
	enc := json.NewEncoder(os.Stdout)
	for { 
		select { 
		case <- ticker.C:
			sl.render(enc)
		case <- sl.Refresh:
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
	return sl
}
