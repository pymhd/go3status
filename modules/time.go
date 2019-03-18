package modules

import (
	"strconv"
	"sync/atomic"
	"time"
)

type TimeModule struct {
	name string
}

func (t TimeModule) Name() string {
	return t.name
}

func (t TimeModule) Run(c chan ModuleOutput, cfg ModuleConfig) {
	//to run by start
	t.run(c, cfg)
	// to run periodically
	ticker := time.NewTicker(cfg.Interval)
	for {
		select {
		case <-ticker.C:
			t.run(c, cfg)
		case <-RefreshChans[cfg.Id]:
			t.run(c, cfg)
		}
	}
}

func (t TimeModule) run(c chan ModuleOutput, cfg ModuleConfig) {
	var layout string
	i, ok := cfg.Extra["format"]
	if !ok {
		layout = "01-02-2006 15:04:05"
	} else {
		layout, _ = i.(string)
	}
	now := time.Now().Format(layout)

	output := ModuleOutput{}
	output.Name = t.name
	output.Instance = strconv.Itoa(cfg.Id)
	output.Refresh = true
	output.Markup = "pango"
	output.FullText = cfg.Prefix

	if x := atomic.LoadInt32(Mute[cfg.Id]); x == -1 {
		output.FullText += "..." + cfg.Postfix
	} else {
		output.FullText += now + cfg.Postfix
	}

	c <- output

}

func (t TimeModule) HandleClickEvent(ce *ClickEvent, cfg ModuleConfig) {
	switch ce.Button {
	// middle, reserved, shrink panel and force refresh
	case 2:
		t.Mute(cfg.Id)
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

func (t TimeModule) Mute(id int) {
	atomic.StoreInt32(Mute[id], ^atomic.LoadInt32(Mute[id]))
}

func init() {
	tm := TimeModule{name: "time"}

	selfRegister(tm)
}
