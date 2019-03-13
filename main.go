package main

import (
	"fmt"
	"time"
	"sync"
	"reflect"
	"encoding/json"
	"go3status/modules"
)

type ModuleOutput struct {
	Align     string `json:"align,omitempty"`
	Color     string `json:"color,omitempty"`
	FullText  string `json:"full_text"`
	Instance  string `json:"instance,omitempty"`
	MinWidth  string `json:"min_width,omitempty"`
	Name      string `json:"name,omitempty"`
	ShortText string `json:"short_text,omitempty"`
	Separator bool   `json:"separator"`
	Urgent    bool   `json:"urgent"`
}

type ClickEvent struct {
	Name     string `json:"name"`
	Instance string `json:"instance"`
	Button   int    `json:"button"`
	XCoord   int    `json:"x"`
	YCoord   int    `json:"y"`
}


type StatusLine struct {
	sync.Mutex
	Header   string
	Modules  []modules.Module
	Statuses []ModuleOutput
	cases    []reflect.SelectCase
}


func (sl *StatusLine) Start() {
	for n, module := range sl.Modules {
		c := make(chan []byte)
		sl.cases[n] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(c)}
		go module.Run(c)
	}
}

func (sl *StatusLine) Run() {
	for {
		ch_num, value, _ := reflect.Select(sl.cases)
		
		mo := new(ModuleOutput)
		msg  := value.Bytes()
		
		if err := json.Unmarshal(msg, mo); err != nil {
			fmt.Println(err)
		}
		//Lock to update Statsuses field
		sl.Lock()
		sl.Statuses[ch_num] = *mo
		sl.Unlock()
	}
}

func (sl *StatusLine) Render() {
	for {
		sl.Lock()
		fmt.Printf("%+v\n", sl.Statuses)
		sl.Unlock()
		time.Sleep(1 * time.Second)
	}
}


func NewStatusLine() *StatusLine {
	sl := new(StatusLine)
	sl.Header = `{"version": 1, "click": true}`
	for _, mod := range modules.Modules {
		sl.Modules = append(sl.Modules, mod)
	}
	sl.Statuses = make([]ModuleOutput, 1)
	sl.cases = make([]reflect.SelectCase, 1)
	return sl
}

func main() {
	s := NewStatusLine()
	s.Start()
	go s.Run()
	s.Render()
}

