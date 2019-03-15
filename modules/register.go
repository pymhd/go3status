package modules

import (
	"fmt"
)

//var Modules = make(map[string]Module, 0)

var availableModules = make(map[string]Module, 0)
func selfRegister(module Module) {
	availableModules[module.Name()] = module
}

var Modules []Module
var Mute = make(map[int]*int32)
var RefreshChans = make(map[int]chan bool)

func Register(id int, key string) {
	mod, ok := availableModules[key]
	if !ok {
		msg := fmt.Sprintf("Module: %s unavailable", key)
		panic(msg)
	}
	//put module in map
	Modules = append(Modules, mod)
	//put refresh chan for modules in map
	c := make(chan bool)
	RefreshChans[id] = c
	//put mute detector in map
	m := int32(0)
	Mute[id] = &m
}
