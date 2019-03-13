package modules

import (
	"time"
)

type Module interface {
	Name() string
	Run(c chan ModuleOutput, cfg ModuleConfig)
}

type ModuleConfig struct {
	Name     string                 `yaml:"name"`
	Interval time.Duration          `yaml:"interval"`
	Prefix   string                 `yaml:"prefix"`
	Postfix  string                 `yaml:"postfix"`
	Colors   map[string]string      `yaml:"colors"`
	Extra    map[string]interface{} `yaml:"extra"`
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
}
