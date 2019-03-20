package modules

import (
	"github.com/mdirkse/i3ipc"
	"unicode/utf8"
	"strconv"
	"sync/atomic"
	"time"
)

var i3socket *i3ipc.IPCSocket

type Title struct {
	name string
}

func (t Title) Name() string {
	return t.name
}

func (t Title) Run(c chan ModuleOutput, cfg ModuleConfig) {
	//open one socket
	var err error
	i3socket, err = i3ipc.GetIPCSocket()
	if err != nil {
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
		case <-RefreshChans[cfg.Id]:
			t.run(c, cfg)
		}
	}
}

func (t Title) run(c chan ModuleOutput, cfg ModuleConfig) {
	output := ModuleOutput{}
	output.Name = t.name
	output.Instance = strconv.Itoa(cfg.Id)
	output.Refresh = true
	output.Markup = "pango"
	output.FullText = cfg.Prefix

	var max int
	m, ok := cfg.Extra["maxChars"]
	if ok {
		max, _ = m.(int)
	}
	if x := atomic.LoadInt32(Mute[cfg.Id]); x == -1 {
		output.FullText += "..." + cfg.Postfix
	} else {
		output.FullText += getFocusedTitle(max) + cfg.Postfix
	}

	c <- output
}

func (t Title) HandleClickEvent(ce *ClickEvent, cfg ModuleConfig) {
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
		execute(cmd, time.Duration(500 * time.Millisecond))
		RefreshChans[cfg.Id] <- true

	}
}

func (t Title) Mute(id int) {
	atomic.StoreInt32(Mute[id], ^atomic.LoadInt32(Mute[id]))
}

func getFocusedTitle(max int) string {
	node, _ := i3socket.GetTree()
	focused := node.FindFocused()
	name := focused.Window_Properties.Title
	length := utf8.RuneCountInString(name)
	if max == 0 || length <= max {
		return name[:length]
	}
	return name[:]

}

func init() {
	t := Title{name: "title"}

	//register plugin to be avail in modele exported map variable Modules
	selfRegister(t)
}
