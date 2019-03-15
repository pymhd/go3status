package modules

import (
	"time"
	"sync/atomic"
)

type TimeModule struct {
	name	string
	refresh chan bool
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
                case <-t.refresh:
                        t.run(c, cfg)
                }
        }
}

func (t TimeModule) run(c chan ModuleOutput, cfg ModuleConfig) {
        var layout string
        i, ok  := cfg.Extra["format"]
        if !ok {
                layout = "01-02-2006 15:04:05"
        } else {
                layout, _ = i.(string)
        }
        now := time.Now().Format(layout)
        
        output := ModuleOutput{}
        output.Name = t.name
        output.Refresh = true
        output.Markup = "pango"
        output.FullText = cfg.Prefix
        
        if x := atomic.LoadInt32(Mute[t.name]); x == -1 {
                output.FullText += "..."
        } else {
                output.FullText += now
        }

        c <- output

}


func (t TimeModule) HandleClickEvent(ce *ClickEvent, cfg ModuleConfig) {
        switch ce.Button {
        // middle, reserved, shrink panel and force refresh
        case 2:
                t.Mute()
                t.refresh <- true
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


func (t TimeModule) Mute() {
        atomic.StoreInt32(Mute[t.name], ^atomic.LoadInt32(Mute[t.name]))
}


func init() {
	c := make(chan bool)
	tm := TimeModule{name: "time", refresh: c}
	
	selfRegister(tm)
}
