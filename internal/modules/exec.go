package modules

import (
	"time"
)

const (
	spawned = ":spawned"
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
	var cacheEnabled bool
	_, ok = cfg.Extra["cache"]
	if ok {
		cacheEnabled, _ = cfg.Extra["cache"].(bool)	
	}
	if cacheEnabled {
		update := time.Duration(10 * time.Second)
		upd, ok := cfg.Extra["update"]
		if ok {
			ts, ok := upd.(string)
			if ok {
				t, err := time.ParseDuration(ts)
				if err == nil {
					update = t
				}
			}
		}
		//if it is first module run
		if workerSpawned := cache.Get(cmd + spawned); workerSpawned == nil {
			updateTicker := time.NewTicker(update)
			go func() {
				cache.Add(cmd+spawned, true, "365d")
				//exec cmd right now then periodically
				o := execute(cmd, timeout)
				cache.Add(cmd, o, "1h")

				for range updateTicker.C {
					o := execute(cmd, timeout)
					cache.Add(cmd, o, "24h")
				}
			}()

		}
		// if worker already was spawned then we wiil wait for latest cache value
		for {
			output, _ := cache.Get(cmd).(string)
			if len(output) > 0 {
				mo.FullText += output
				return
			}
			time.Sleep(100 * time.Millisecond)
		}
	}
	//if cache is not used fir this module
	mo.FullText += execute(cmd, timeout)
	
}

func init() {
	RegisteredFuncs["exec"] = execCmd
}
