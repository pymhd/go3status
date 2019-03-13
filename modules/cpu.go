package modules

import (
	"time"
	"encoding/json"
)

const (
	base = iota
	norm = iota + 1
	warn
	crit
)

type CPU struct {
	name	string
}

func (cpu CPU) Run(c chan []byte, cfg ModuleConfig) {
	// to run on Start()
	cpu.run(c, cfg)
	
	// to run periodically
	ticker := time.NewTicker(cfg.Interval)
	for range ticker.C {
		cpu.run(c, cfg)
	}
}

func (cpu CPU) Name() string {
	return cpu.name
}

func (cpu CPU) run (c chan []byte, cfg ModuleConfig) {
	output := ModuleOutput{}
	//do some stuff
	// ...
	// ...
	output.FullText = "27%"
	output.Color = cfg.Colors["good"]
        //data, _  := json.Marshal(output)
        json.NewEncoder(CW{c}).Encode(output)
        //c <- []byte(data)
}

func init() {
	cpu := CPU{"cpu"}
	Register("cpu", cpu)
}


type CW struct {
	ch (chan []byte)
}


func (c CW) Write(b []byte) (n int, err error) {
	c.ch <- b 
	return 1, nil		
}

