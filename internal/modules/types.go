package modules

import (
	"time"
)

type ModuleConfig struct {
	Id          int
	Name        string                 `yaml:"name"`
	Interval    time.Duration          `yaml:"interval"`
	Prefix      string                 `yaml:"prefix"`
	Postfix     string                 `yaml:"postfix"`
	Colors      map[string]string      `yaml:"colors"`
	Levels      map[string]string      `yaml:"levels"`
	ClickEvents map[string]string      `yaml:"clickEvents"`
	ShortFormat bool                   `yaml:"short"`
	Extra       map[string]interface{} `yaml:"extra"`
}

type ModuleOutput struct {
	Align      string `json:"align,omitempty"`
	Color      string `json:"color,omitempty"`
	FullText   string `json:"full_text"`
	Instance   string `json:"instance,omitempty"`
	MinWidth   string `json:"min_width,omitempty"`
	Name       string `json:"name,omitempty"`
	ShortText  string `json:"short_text,omitempty"`
	Separator  bool   `json:"separator"`
	Urgent     bool   `json:"urgent"`
	Background string `json:"background,omitempty"`
	//pango
	Markup string `json:"markup"`
	// do i need to refresh
	refresh bool
}

type ClickEvent struct {
	Name     string   `json:"name"`
	Instance string   `json:"instance"`
	Button   int      `json:"button"`
	X        int      `json:"x"`
	Y        int      `json:"y"`
	Mod      []string `json:"modifiers"`
}
