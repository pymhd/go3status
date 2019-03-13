package modules

import (
	"os/exec"
	"time"
)

type CPU struct {
	name string
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

func (cpu CPU) run(c chan ModuleOutput, cfg ModuleConfig) {
	output := ModuleOutput{}

	output.FullText = "27% to run periodically ChanWriter{Chan:"
	output.ShortText = "27%"
	output.Color = cfg.Colors["good"]
	output.Name = cpu.name
	//output.Markup = "pango"
	//output.Background = "#ffffff"

	c <- output
}

func (cpu CPU) HandleClickEvent(ce *ClickEvent) {
	cmd := exec.Command("urxvt", "-name", "__scratchpad", "-e", "htop")	
	err := cmd.Start()
	if err != nil {
		panic(err)
	}
}

func init() {
	cpu := CPU{"cpu"}
	Register("cpu", cpu)
}
