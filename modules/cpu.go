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
	for {
		output := ModuleOutput{FullText: "80%", Name: "cpu"}
		data, err  := json.Marshal(output)
		if err != nil {
			panic(err)
		}
		c <- []byte(data)
		time.Sleep(1 * time.Second)
	}
}

func (cpu CPU) Name() string {
	return cpu.name
}

func init() {
	cpu := CPU{"cpu"}
	Register("cpu", cpu)
}
