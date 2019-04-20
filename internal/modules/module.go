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
	layout  int
	//mute    int
	//short	int
	Refresh chan bool
}

func (m *Module) Run(f func(*ModuleOutput, ModuleConfig)) {
	//create module output to send
	mo := new(ModuleOutput)
	m.preloadOutput(mo)

	//run func on startup
	f(mo, m.Cfg)
	m.sendOutput(mo)
	if  m.Cfg.ShortFormat {
		m.layout = 1
	}
	ticker := time.NewTicker(m.Cfg.Interval)
	for {
		select {
		case <-ticker.C:
			switch (m.layout +3 ) % 3 {
			case 0:
				f(mo, m.Cfg)
			case 1:
				f(mo, m.Cfg)
				mo.FullText = mo.ShortText
			case 2:
				m.muteOutput(mo)
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
			switch (m.layout +3 ) % 3 {
                        case 0:
                                f(mo, m.Cfg)
                        case 1:
                                f(mo, m.Cfg)
                                mo.FullText = mo.ShortText
                        case 2:
                                m.muteOutput(mo)
                        }
			m.postLoadOutput(mo)
			m.sendOutput(mo)
		}
	}
}

func (m *Module) HandleClickEvent(ce *ClickEvent) {
	switch ce.Button {
	// middle, reserved, change layout panel and force refresh
	case 2:
		m.layout++
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
	mo.ShortText = m.Cfg.Prefix
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
