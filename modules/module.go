package modules

import (
	"fmt"
	"strconv"
	"time"
)

type Module struct {
	Name    string
	Update  chan ModuleOutput
	Cfg     ModuleConfig
	mute    int
	short	int
	Refresh chan bool
}

func (m *Module) Run(f func(*ModuleOutput, ModuleConfig)) {
	//create module output to send
	mo := new(ModuleOutput)
	m.preloadOutput(mo)

	//run func on startup
	f(mo, m.Cfg)
	m.sendOutput(mo)
	ticker := time.NewTicker(m.Cfg.Interval)
	for {
		select {
		case <-ticker.C:
			if m.short == -1 {
				m.Cfg.Extra["format"] = "short"
			} else {
				m.Cfg.Extra["format"] = "long"
			}
			if m.mute == -1 {
				m.muteOutput(mo)
			} else {
				f(mo, m.Cfg)
			}
			cacheKey := fmt.Sprintf("result:%d", m.Cfg.Id)
			previousValue, _ := cache.Get(cacheKey).(string)
			currentValue := mo.FullText
			if currentValue != previousValue {
				m.postLoadOutput(mo)
	                        m.sendOutput(mo)
			}
			cache.Add(cacheKey, currentValue, "1h")
			m.flushOutput(mo)
		case <-m.Refresh:
			if m.short == -1 {
				m.Cfg.Extra["format"] = "short"
			} else {
				m.Cfg.Extra["format"] = "long"
			}

			if m.mute == -1 {
				m.muteOutput(mo)
			} else {
				f(mo, m.Cfg)
			}
			m.postLoadOutput(mo)
			m.sendOutput(mo)
		}
	}
}

func (m *Module) HandleClickEvent(ce *ClickEvent) {
	switch ce.Button {
	// middle, reserved, shrink panel and force refresh
	case 2:
		if len(ce.Mod) > 0 {
			if ce.Mod[0] == "Shift" {
				m.refresh()
				break
			}
			if ce.Mod[0] == "Control" {
				m.short = ^m.short
				m.refresh()
				break
			}
		}
		m.mute = ^m.mute
		m.refresh()
	// any other
	default:
		buttonNumber := ce.Button
		buttonText := clickMap[buttonNumber]
		cmd, ok := m.Cfg.ClickEvents[buttonText]
		if !ok {
			//if no cmd specified in config file
			break
		}
		execute(cmd, time.Duration(500*time.Millisecond))
		m.refresh()
	}
}

func (m *Module) preloadOutput(mo *ModuleOutput) {
	mo.Name = m.Name
	mo.Instance = strconv.Itoa(m.Cfg.Id)
	mo.Markup = "pango"
	mo.FullText = m.Cfg.Prefix
}

func (m *Module) postLoadOutput(mo *ModuleOutput) {
	mo.FullText += m.Cfg.Postfix
}

func (m *Module) muteOutput(mo *ModuleOutput) {
	mo.FullText += "..."
}
/*
func (m *Module) shortOutput(mo *ModuleOutput) {
	mo.FullText += mo.ShortText
}
*/
func (m *Module) flushOutput(mo *ModuleOutput) {
	mo.FullText = m.Cfg.Prefix
}

func (m *Module) sendOutput(mo *ModuleOutput) {
	m.Update <- *mo
	m.flushOutput(mo)
}

func (m *Module) refresh() {
	m.Refresh <- true
}
