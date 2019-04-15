package modules

import (
	"time"
)

func execCmd(mo *ModuleOutput, cfg ModuleConfig) {
	_, ok := cfg.Extra["cmd"]
	if !ok {
		mo.FullText += "Provide command"
		return
	}
	cmd, _ := cfg.Extra["cmd"].(string)
	_, ok = cfg.Extra["color"]
	if ok {
		mo.Color, _ = cfg.Extra["color"].(string)
	}
	timeout := time.Duration(500 * time.Millisecond)
	tmt, ok := cfg.Extra["timeout"]
	if ok {
		ts, ok := tmt.(string)
		if ok {
			t, err := time.ParseDuration(ts)
			if err == nil {
				timeout = t
			}
		}
	}
	if cfg.ShortFormat {
		mo.FullText = cfg.Prefix
	} else {
		mo.FullText += execute(cmd, timeout)
	}
}

func init() {
	RegisteredFuncs["exec"] = execCmd
}
