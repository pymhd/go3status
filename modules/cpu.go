package modules

import (
	"os"
	"fmt"
	"time"
	"strconv"
	"sync/atomic"
)

var ( 
	s string
	//those ints need to store previous values of cpu time
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
	//to run by start
	cpu.run(c, cfg)

	// to run periodically
	ticker := time.NewTicker(cfg.Interval)
	for {
		select {
		case <-ticker.C:
			cpu.run(c, cfg)
		case <-cpu.refresh:
			cpu.run(c, cfg)
		}
	}
}

func (cpu CPU) run(c chan ModuleOutput, cfg ModuleConfig) {
	output := ModuleOutput{}
	output.Name = cpu.name
	output.Instance = strconv.Itoa(cfg.Id)
	output.Refresh = true
	output.Markup = "pango"
	output.FullText = cfg.Prefix
	
	percentage := getCpuPercentage()
	for lvl, val := range cfg.Levels {
		if inRange(percentage, val) {
			output.Color = cfg.Colors[lvl]
		}
	}
	if x := atomic.LoadInt32(Mute[cfg.Id]); x == -1 {
		output.FullText += " ..."
	} else {
		output.FullText += fmt.Sprintf(" %.2f%%%s", percentage, cfg.Postfix)
	}
	c <- output
}

func (cpu CPU) HandleClickEvent(ce *ClickEvent, cfg ModuleConfig) {
	switch ce.Button {
	// middle, reserved, shrink panel and force refresh
	case 2:
		cpu.Mute(cfg.Id)
		cpu.refresh <- true
	// any other
	default:
		buttonNumber := ce.Button
		buttonText := clickMap[buttonNumber]
		cmd, ok := cfg.ClickEvents[buttonText]
                if !ok {
                	//if no cmd specified in config file
                        break
                }
                execute(cmd)

	}
}

func (cpu CPU) Mute(id int) {
	atomic.StoreInt32(Mute[id], ^atomic.LoadInt32(Mute[id]))
}


func getCpuPercentage() float64 {
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
        return cpu
}


func init() {
	c := make(chan bool)
	cpu := CPU{name: "cpu", refresh: c}

	//register plugin to be avail in modele exported map variable Modules
	selfRegister(cpu)
}

