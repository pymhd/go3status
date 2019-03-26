package modules

var RegisteredFuncs = make(map[string]func(*ModuleOutput, ModuleConfig))
