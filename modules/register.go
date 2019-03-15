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
var Mute = make(map[string]*int32, 0)

func Register(id int, key string) {
	mod, ok := availableModules[key]
	if !ok {
		msg := fmt.Sprintf("Module: %s unavailable", key)
		panic(msg)
	}
	Modules = append(Modules, mod)
	m := int32(0)
	Mute[mod.Name()] = &m	
}
