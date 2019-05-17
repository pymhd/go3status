package modules

import (
	"time"
)

const (
	spawned = ":spawned"
)

type execExtraConf struct {
	cmd	string
	color	string
	timeout time.Duration
	update  time.Duration 
	cacheEnabled bool
	
}

func execCmd(mo *ModuleOutput, cfg ModuleConfig) {
	ecfg := execExtraConf{timeout: time.Duration(500 * time.Millisecond), update: time.Duration(500 * time.Millisecond)}
	for k, v := range cfg.Extra {
		switch k {
		case "cmd":
			ecfg.cmd, _ = v.(string)
			if len(ecfg.cmd) == 0 {
				mo.FullText += "conf err"
				return
			}
		case "color":
			ecfg.color, _ = v.(string)
		case "timeout":
			vs, _ := v.(string)
			td, err := time.ParseDuration(vs)
			if err != nil {
                                mo.FullText += "conf err"
                                return
                        }
			ecfg.timeout = td
		case "update":
			vs, _ := v.(string)
                        td, err := time.ParseDuration(vs)
                        if err != nil {
                        	mo.FullText += "conf err"
                        	return
                        }
                        ecfg.update = td
		case "cache":
			ecfg.cacheEnabled, _ = v.(bool)
		}
	}
	mo.Color = ecfg.color
	if ecfg.cacheEnabled {
		//if it is first module run
		if workerSpawned := cache.Get(ecfg.cmd + spawned); workerSpawned == nil {
			updateTicker := time.NewTicker(ecfg.update)
			go func() {
				cache.Add(ecfg.cmd+spawned, true, "365d")
				//exec ecfg.cmd right now then periodically
				o := execute(ecfg.cmd, ecfg.timeout)
				cache.Add(ecfg.cmd, o, "1h")

				for range updateTicker.C {
					o := execute(ecfg.cmd, ecfg.timeout)
					cache.Add(ecfg.cmd, o, "24h")
				}
			}()

		}
		// if worker already was spawned then we wiil wait for latest cache value
		for {
			output, _ := cache.Get(ecfg.cmd).(string)
			if len(output) > 0 {
				mo.FullText += output
				return
			}
			time.Sleep(100 * time.Millisecond)
		}
	}
	//if cache is not used for this module
	mo.FullText += execute(ecfg.cmd, ecfg.timeout)
	
}

func init() {
	RegisteredFuncs["exec"] = execCmd
}
