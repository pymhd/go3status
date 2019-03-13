package modules 

import (

)

var Modules = make(map[string]Module, 0)


func Register(name string, module Module) {
        Modules[name] = module        
}
