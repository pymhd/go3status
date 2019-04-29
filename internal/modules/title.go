package modules

import (
	"github.com/mdirkse/i3ipc"
	"unicode/utf8"
)

const (
	i3socketKey = "title:socket"
)

func getFocusedTitle(mo *ModuleOutput, cfg ModuleConfig) {
	var i3socket *i3ipc.IPCSocket
	cv := cache.Get(i3socketKey)
	if cv != nil {
		i3socket, _ = cv.(*i3ipc.IPCSocket)
	} else {
		var err error
		i3socket, err = i3ipc.GetIPCSocket()
		if err != nil {
			return
		}
		cache.Add(i3socketKey, i3socket, "12h")
	}
	node, _ := i3socket.GetTree()
	focused := node.FindFocused()
	name := focused.Window_Properties.Title
	var max int
	v, ok := cfg.Extra["maxChars"]
	if ok {
		max, _ = v.(int)
	}
	length := utf8.RuneCountInString(name)
	if max == 0 || length <= max {
		mo.FullText += name
	} else {
		mo.FullText += name[:max]
	}
}

func init() {
	RegisteredFuncs["title"] = getFocusedTitle
}
