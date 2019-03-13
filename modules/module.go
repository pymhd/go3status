package modules

import (
	"time"
)


type Module interface {
	Name() string
	Run(c chan []byte, cfg ModuleConfig)
}

type ModuleConfig struct {
	Name     string         `yaml:"name"`
	Interval time.Duration  `yaml:"interval"`
	Prefix   string         `yaml:"prefix"`
	Postfix  string         `yaml:"postfix"`
	Colors   map[string]string `yaml:"colors"`
	Extra	 map[string]interface{} `yaml:"extra"`	
}
