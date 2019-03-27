package main

import (
	"encoding/json"
	"fmt"
	"go3status/modules"
	"os"
	"reflect"
	"sync"
)

type StatusLine struct {
	sync.Mutex
	Header  string
	Refresh chan bool
	Modules []*modules.Module
	Blocks  []modules.ModuleOutput
	cases   []reflect.SelectCase
}

func (sl *StatusLine) Start() {
	for n, module := range sl.Modules {
		sl.cases[n] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(module.Update)}
		f, ok := modules.RegisteredFuncs[module.Name]
		if ok {
			go module.Run(f)
			log.Infof("Module %q started\n", module.Name)
		} else {
			log.Errorf("Module %q not found (Check if source code registered)\n", module.Name)
		}
	}
}

func (sl *StatusLine) Run() {
	for {
		chNum, value, _ := reflect.Select(sl.cases)

		mo, ok := value.Interface().(modules.ModuleOutput)
		if !ok {
			// why
			fmt.Println()
		}
		//Lock to update Statsuses field
		sl.Lock()
		sl.Blocks[chNum] = mo
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
		<-sl.Refresh
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

	sl.Modules = make([]*modules.Module, len(cfg.Modules))
	for n, moduleCm := range cfg.Modules {
		for name, mcfg := range moduleCm {
			upd := make(chan modules.ModuleOutput)
			rfsh := make(chan bool)
			//used later as Instance attr in module output to distinct same modules
			mcfg.Id = n
			
			m := new(modules.Module)
			m.Name = name
			m.Update = upd
			m.Refresh = rfsh
			m.Cfg = mcfg
			
			sl.Modules[n] = m
			log.Infof("Configured module %q with update interval: %s\n", name, mcfg.Interval)
		}
	}

	sl.Blocks = make([]modules.ModuleOutput, len(cfg.Modules))
	sl.cases = make([]reflect.SelectCase, len(cfg.Modules))
	sl.Refresh = make(chan bool, 0)
	//fmt.Println(modules.Mute)
	return sl
}
