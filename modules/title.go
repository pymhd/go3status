package modules

import (
	"sync/atomic"
	"time"
	"github.com/mdirkse/i3ipc"
)

var i3socket *i3ipc.IPCSocket

type Title struct {
	name    string
	refresh chan bool
}

func (t Title) Name() string {
	return t.name
}

func (t Title) Run(c chan ModuleOutput, cfg ModuleConfig) {
	//open one socket 
	var err error 
	i3socket, err = i3ipc.GetIPCSocket()
	if err != nil  {
		//FIXME
		panic(err)
	}

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

func (t Title) run(c chan ModuleOutput, cfg ModuleConfig) {
	output := ModuleOutput{}
	output.Name = t.name
	output.Refresh = true
	output.Markup = "pango"
	output.FullText = cfg.Prefix

	if x := atomic.LoadInt32(Mute[t.name]); x == -1 {
                output.FullText += "..."
        } else {
                output.FullText += getFocusedTitle()
        }

	c <- output
}

func (t Title) HandleClickEvent(ce *ClickEvent, cfg ModuleConfig) {
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

func (t Title) Mute() {
	atomic.StoreInt32(Mute[t.name], ^atomic.LoadInt32(Mute[t.name]))
}


func getFocusedTitle() string {
        node, _ := i3socket.GetTree()
        focused := node.FindFocused()
        return focused.Window_Properties.Title

}

func init() {
	c := make(chan bool)
	t := Title{name: "title", refresh: c}

	//register plugin to be avail in modele exported map variable Modules
	selfRegister(t)
}

