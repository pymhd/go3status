package modules

import (
        "go-logging"
)

var log  *logger.Logger
var RegisteredFuncs = make(map[string]func(*ModuleOutput, ModuleConfig))

func RegisterLogger(l *logger.Logger) {
    log = l
}
