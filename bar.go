package main

import (
	"fmt"
	"encoding/json"
	"go3status/modules"
	"os"
	"reflect"
	"sync"
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
		// mc - module config
		// id needed to mute by order
		mc := cfg.Modules[n][module.Name()]
		mc.Id = n
		go module.Run(c, mc)
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
		//Better to refresh every time we accept update
		//no need to print by ticker the same info
		sl.Refresh <- true
	}
}

func (sl *StatusLine) Render() {
	// ...
	fmt.Println(sl.Header)
	fmt.Printf("[[]\n,")

	enc := json.NewEncoder(os.Stdout)
	for {
		<- sl.Refresh
		sl.render(enc)
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
	
	sl.Modules = make([]modules.Module, len(cfg.Modules))
	for n, modmap := range cfg.Modules {
		for name, _ := range modmap {
			// thist panics if module not in avail modules
			modules.Register(n, name) 
			sl.Modules[n] = modules.Modules[n]
		}
	}
	
	sl.Blocks = make([]modules.ModuleOutput, len(cfg.Modules))
	sl.cases = make([]reflect.SelectCase, len(cfg.Modules))
	sl.Refresh = make(chan bool, 0)
	//fmt.Println(modules.Mute)
	return sl
}
