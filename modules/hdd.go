package modules

import (
	"bytes"
	"strconv"
	"sync/atomic"
	"syscall"
	"text/template"
	"time"
)

const (
	BaseTemplate   = `{{printf "%.1f" $used}}/{{printf "%.1f" $total}} (printf "%.0f" $percentage)`
	TemplatePrefix = `{{$path := .Path}}{{$used := .Used}}{{$total := .Total}}{{$avail := .Avail}}{{$percentage := .Percentage}}`
)

type filesystem struct {
	Path       string
	Used       float64
	Avail      float64
	Total      float64
	Percentage float64
}

type HDD struct {
	name string
}

func (h HDD) Name() string {
	return h.name
}

func (h HDD) Run(c chan ModuleOutput, cfg ModuleConfig) {
	//to run by start
	h.run(c, cfg)
	// to run periodically
	ticker := time.NewTicker(cfg.Interval)
	for {
		select {
		case <-ticker.C:
			h.run(c, cfg)
		case <-RefreshChans[cfg.Id]:
			h.run(c, cfg)
		}
	}
}

func (h HDD) run(c chan ModuleOutput, cfg ModuleConfig) {
	output := ModuleOutput{}
	output.Name = h.name
	output.Instance = strconv.Itoa(cfg.Id)
	output.Refresh = true
	output.Markup = "pango"
	output.FullText = cfg.Prefix

	mp, ok := cfg.Extra["mountPoint"]
	if !ok {
		output.FullText += "N/A"
		c <- output
		return
	}
	mountPoint, ok := mp.(string)
	if !ok {
		output.FullText += "N/A"
		c <- output
		return
	}

	fs := filesystem{}
	fs.Used, fs.Total = getMpStats(mountPoint)
	fs.Avail = fs.Total - fs.Used
	fs.Path = mountPoint
	fs.Percentage = fs.Used * 100 / fs.Total
	var T string
	tpl, ok := cfg.Extra["format"]
	if !ok {
		T = BaseTemplate
	} else {
		T, _ = tpl.(string)
		if len(T) == 0 {
			//empty format key
			T = BaseTemplate
		}
	}
	//to register vars
	T = TemplatePrefix + T
	t := template.Must(template.New(mountPoint).Parse(T))

	var o bytes.Buffer
	t.Execute(&o, fs)
	if x := atomic.LoadInt32(Mute[cfg.Id]); x == -1 {
		output.FullText += "..." + cfg.Postfix
	} else {
		output.FullText += o.String() + cfg.Postfix
	}
	for lvl, val := range cfg.Levels {
		if inRange(fs.Percentage, val) {
			output.Color = cfg.Colors[lvl]
		}
	}
	c <- output

}

func (h HDD) HandleClickEvent(ce *ClickEvent, cfg ModuleConfig) {
	switch ce.Button {
	// middle, reserved, shrink panel and force refresh
	case 2:
		h.Mute(cfg.Id)
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

func (h HDD) Mute(id int) {
	atomic.StoreInt32(Mute[id], ^atomic.LoadInt32(Mute[id]))
}

func getMpStats(path string) (float64, float64) {
	stats := syscall.Statfs_t{}
	syscall.Statfs(path, &stats)
	avail := float64(stats.Bavail) * float64(stats.Bsize)
	total := float64(stats.Blocks) * float64(stats.Bsize)
	return (total - avail) / 1024 / 1024 / 1024, total / 1024 / 1024 / 1024
}

func init() {
	hdd := HDD{name: "hdd"}

	selfRegister(hdd)
}
