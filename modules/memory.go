package modules

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

type Memory struct {
	name string
}

func (m Memory) Name() string {
	return m.name
}

func (m Memory) Run(c chan ModuleOutput, cfg ModuleConfig) {
	//to run by start
	m.run(c, cfg)

	// to run periodically
	ticker := time.NewTicker(cfg.Interval)
	for {
		select {
		case <-ticker.C:
			m.run(c, cfg)
		case <-RefreshChans[cfg.Id]:
			m.run(c, cfg)
		}
	}
}

func (m Memory) run(c chan ModuleOutput, cfg ModuleConfig) {
	output := ModuleOutput{}
	output.Name = m.name
	output.Instance = strconv.Itoa(cfg.Id)
	output.Refresh = true
	output.Markup = "pango"
	output.FullText = cfg.Prefix

	memUsed, memTotal := getMemory()
	percentage := 100 * (memUsed / memTotal)

	for lvl, val := range cfg.Levels {
		if inRange(percentage, val) {
			output.Color = cfg.Colors[lvl]
		}
	}
	memoryRepr := fmt.Sprintf("%.1f/%.1f (%.0f%%)%s", memUsed/1048576, memTotal/1048576, percentage, cfg.Postfix)

	if x := atomic.LoadInt32(Mute[cfg.Id]); x == -1 {
		output.FullText += "..." + cfg.Postfix
	} else {
		output.FullText += memoryRepr
	}

	c <- output
}

func (m Memory) HandleClickEvent(ce *ClickEvent, cfg ModuleConfig) {
	switch ce.Button {
	// middle, reserved, shrink panel and force refresh
	case 2:
		m.Mute(cfg.Id)
		RefreshChans[cfg.Id] <- true
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
		RefreshChans[cfg.Id] <- true

	}
}

func (m Memory) Mute(id int) {
	atomic.StoreInt32(Mute[id], ^atomic.LoadInt32(Mute[id]))
}

func getMemory() (float64, float64) {
	meminfo, _ := os.Open("/proc/meminfo")
	defer meminfo.Close()

	var total, avail string
	var done int
	scanner := bufio.NewScanner(meminfo)
	for scanner.Scan() {
		line := scanner.Text()
		sl := strings.Fields(line)
		switch sl[0] {
		case "MemTotal:":
			total = sl[1]
			done += 1
		case "MemAvailable:":
			avail = sl[1]
			done += 1
		}
		if done == 2 {
			break
		}
	}
	totalF, _ := strconv.ParseFloat(total, 64)
	availF, _ := strconv.ParseFloat(avail, 64)

	return totalF - availF, totalF
}

func init() {
	m := Memory{name: "memory"}

	//register plugin to be avail in modele exported map variable Modules
	selfRegister(m)
}
