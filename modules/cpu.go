package modules

import (
	"time"
)


type CPU struct {
	name	string
}

func (cpu CPU) Name() string {
        return cpu.name
}

func (cpu CPU) Run(c chan ModuleOutput, cfg ModuleConfig) {
	//w := ChanWriter{Chan: c}
	cpu.run(c, cfg)
	
	// to run periodically
	ticker := time.NewTicker(cfg.Interval)
	for range ticker.C {
		cpu.run(c, cfg)
	}
}

func (cpu CPU) run (c chan ModuleOutput, cfg ModuleConfig ) {
	output := ModuleOutput{}
	
	output.FullText = "27%"
	output.Color = cfg.Colors["good"]
	
	c <- output
}

func init() {
	cpu := CPU{"cpu"}
	Register("cpu", cpu)
}
