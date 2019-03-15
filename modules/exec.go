package modules

import (
	"strconv"
	"sync/atomic"
	"time"
)


type Exec struct {
	name    string
	refresh chan bool
}

func (e Exec) Name() string {
	return e.name
}

func (e Exec) Run(c chan ModuleOutput, cfg ModuleConfig) {
	//to run by start
	e.run(c, cfg)

	// to run periodically
	ticker := time.NewTicker(cfg.Interval)
	for {
		select {
		case <-ticker.C:
			e.run(c, cfg)
		case <-e.refresh:
			e.run(c, cfg)
		}
	}
}

func (e Exec) run(c chan ModuleOutput, cfg ModuleConfig) {
	output := ModuleOutput{}
	output.Name = e.name
	output.Instance = strconv.Itoa(cfg.Id)
	output.Refresh = true
	output.Markup = "pango"
	output.FullText = cfg.Prefix
	
	s, ok := cfg.Extra["cmd"]
	if !ok {
		output.FullText = "Provide command"
		output.Color = "#7f0909"
		c <- output
		return			
	}
	cmd, ok := s.(string)
	if !ok {
                output.FullText = "Wrong cmd"
                output.Color = "#7f0909"
                c <- output
                return     
        }
        color, ok := cfg.Extra["color"]
        if ok {
        	c, ok := color.(string)
        	if ok {
        		output.Color = c
        	}
        }
	if x := atomic.LoadInt32(Mute[cfg.Id]); x == -1 {
                output.FullText += "..."
        } else {
                output.FullText += execute(cmd)
        }

	c <- output
}

func (e Exec) HandleClickEvent(ce *ClickEvent, cfg ModuleConfig) {
	switch ce.Button {
	// middle, reserved, shrink panel and force refresh
	case 2:
		e.Mute(cfg.Id)
		e.refresh <- true
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

func (e Exec) Mute(id int) {
	atomic.StoreInt32(Mute[id], ^atomic.LoadInt32(Mute[id]))
}


func init() {
	c := make(chan bool)
	e := Exec{name: "exec", refresh: c}

	//register plugin to be avail in modele exported map variable Modules
	selfRegister(e)
}
