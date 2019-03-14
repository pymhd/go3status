package modules

import (
	"os"
	"fmt"
	"strings"
	"os/exec"
	"sync/atomic"
	"time"
)

var ( 
	s string
	puser, pnice, psystem, pidle, pio, pirq, psoftirq, psteal, pguest, pguest_nice int
)

type CPU struct {
	name    string
	refresh chan bool
}

func (cpu CPU) Name() string {
	return cpu.name
}

func (cpu CPU) Run(c chan ModuleOutput, cfg ModuleConfig) {
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

func (cpu CPU) run(c chan ModuleOutput, cfg ModuleConfig, refresh bool) {
	output := ModuleOutput{}
	output.Name = cpu.name
	output.Refresh = refresh
	output.Markup = "pango"
	output.FullText = cfg.Prefix
	
	
	output.FullText = getCpuPercentage()
	if x := atomic.LoadInt32(Mute[cpu.name]); x == -1 {
		output.FullText = "33%"
	}
	output.Color = cfg.Colors["good"]

	c <- output
}

func (cpu CPU) HandleClickEvent(ce *ClickEvent, cfg ModuleConfig) {
	switch ce.Button {
	// left click, get from cfg
	case 1:
		cmd, ok := cfg.ClickEvents[leftClick]
		if !ok {
			break
		}
		execute(cmd)
	// middle, reserved, shrink panel and force refresh
	case 2:
		cpu.Mute()
		cpu.refresh <- true
	}
}

func (cpu CPU) Mute() {
	atomic.StoreInt32(Mute[cpu.name], ^atomic.LoadInt32(Mute[cpu.name]))
}


func execute(cmd string) {
	args := strings.Split(cmd, " ")
	c := exec.Command(args[0], args[1:]...)
	c.Start()
}

func getCpuPercentage() string {
        var user, nice, system, idle, io, irq, softirq, steal, guest, guest_nice int
        
        stat, _ := os.Open("/proc/stat")
        defer stat.Close()
        
        fmt.Fscanf(stat, "%s %d %d %d %d %d %d %d %d %d %d", &s, &user, &nice, &system, &idle, &io, &irq, &softirq, &steal, &guest, &guest_nice)
        
        PrevIdle := pidle + pio
        Idle := idle + io
        
        PrevNonIdle := puser + pnice + psystem + pirq + psoftirq + psteal
        NonIdle :=  user + nice + system + irq + softirq + steal
        
        PrevTotal := PrevIdle + PrevNonIdle
        Total := Idle + NonIdle
        
        totald := Total - PrevTotal
        idled := Idle - PrevIdle

        cpu := 100 * float64(totald - idled)/float64(totald)
        puser, pnice, psystem, pidle, pio, pirq, psoftirq, psteal, pguest, pguest_nice = user, nice, system, idle, io, irq, softirq, steal, guest, guest_nice
        return fmt.Sprintf("%.2f%%", cpu) 
}


func init() {
	rch := make(chan bool)
	cpu := CPU{"cpu", rch}
	Register("cpu", cpu)
}
