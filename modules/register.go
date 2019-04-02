package modules

import (
	"github.com/pymhd/go-logging"
)

var log logger.Logger
var RegisteredFuncs = make(map[string]func(*ModuleOutput, ModuleConfig))

func RegisterLogger(l logger.Logger) {
	log = l
}
