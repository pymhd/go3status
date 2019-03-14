package modules

import ()

var Modules = make(map[string]Module, 0)
var Mute = make(map[string]*int32, 0)

func Register(module Module) {
	Modules[module.Name()] = module
	//to store mute state and change it using atomic operations
	var m int32
	Mute[module.Name()] = &m
}
