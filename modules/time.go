package modules

import (
	"time"
)

const (
	defaultLayout = "02/01 15:04"
)

func getTime(mo *ModuleOutput, cfg ModuleConfig) {
	var layout string
	v, ok := cfg.Extra["format"]
	if !ok {
		layout = defaultLayout
	} else {
		layout, _ = v.(string)
	}
	now := time.Now().Format(layout)
	mo.FullText += now
}

func init() {
	RegisteredFuncs["time"] = getTime
}
