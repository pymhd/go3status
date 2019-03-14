package modules

import (
	_ "fmt"
	_ "os/exec"
	"sync/atomic"
	"time"
)

type CPU struct {
	name    string
	refresh chan bool
}

func (cpu CPU) Name() string {
	return cpu.name
}

func (cpu CPU) Run(c chan ModuleOutput, cfg ModuleConfig) {
	//w := ChanWriter{Chan: c}
	cpu.run(c, cfg, false)

	// to run periodically
	ticker := time.NewTicker(cfg.Interval)
	for {
		select {
		case <-ticker.C:
			//fmt.Println("ticker")
			cpu.run(c, cfg, false)
		case <-cpu.refresh:
			//fmt.Println("by refresh")
			cpu.run(c, cfg, true)
		}
	}
}

func (cpu CPU) run(c chan ModuleOutput, cfg ModuleConfig, urgent bool) {
	output := ModuleOutput{}

	output.FullText = "27% to run periodically ChanWriter{Chan:"
	if x := atomic.LoadInt32(Mute[cpu.name]); x == -1 {
		output.FullText = "33%"
	}
	output.ShortText = "27%"
	output.Color = cfg.Colors["good"]
	output.Name = cpu.name
	output.Refresh = urgent
	output.Markup = "pango"
	//output.Background = "#ffffff"

	c <- output
}

func (cpu CPU) HandleClickEvent(ce *ClickEvent) {
	cpu.Mute()
	cpu.refresh <- true
}

func (cpu CPU) Mute() {
	atomic.StoreInt32(Mute[cpu.name], ^atomic.LoadInt32(Mute[cpu.name]))
}

func init() {
	rch := make(chan bool)
	cpu := CPU{"cpu", rch}
	Register("cpu", cpu)
}
