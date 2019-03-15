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
var Mute  []*int32

func Register(key string) {
	mod, ok := availableModules[key]
	if !ok {
		msg := fmt.Sprintf("Module: %s unavailable", key)
		panic(msg)
	}
	if !alreadyRegistered(key) {
		Modules = append(Modules, mod)
	} else {
		newRefreshChan := make(chan bool)
		mod.refresh = newRefreshChan
		Modules = append(Modules, mod)
	}
	//mute atomic per module 
	m := int32(0)
	Mute = append(Mute, &m)	
}


func alreadyRegistered(m string) bool {
	for _, mod := range Modules {
		if mod.Name() == m {
			return true
		}
	}
	return false
}

