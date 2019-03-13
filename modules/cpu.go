package modules

import (
	"time"
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
		c <- []byte(`{"full_text": "80%", "separator": true}`)
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
