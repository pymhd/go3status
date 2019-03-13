package modules

import (
	"time"
	"encoding/json"
)


type CPU struct {
	name	string
}

func (cpu CPU) Name() string {
        return cpu.name
}

func (cpu CPU) Run(c chan []byte, cfg ModuleConfig) {
	w := ChanWriter{Chan: c}
	enc := json.NewEncoder(w)
	
	cpu.run(c, cfg, enc)
	
	// to run periodically
	ticker := time.NewTicker(cfg.Interval)
	for range ticker.C {
		cpu.run(c, cfg, enc)
	}
}

func (cpu CPU) run (c chan []byte, cfg ModuleConfig, e *json.Encoder) {
	output := ModuleOutput{}
	
	output.FullText = "27%"
	output.Color = cfg.Colors["good"]

        e.Encode(output)
}

func init() {
	cpu := CPU{"cpu"}
	Register("cpu", cpu)
}
