package main 


import (
        "time"
        "io/ioutil"
        "gopkg.in/yaml.v2"
        "go3status/modules"
)

type Config struct {
        Global	struct {
            Interval	time.Duration `yaml:"interval"`
            Color	string	`yaml:"color"`
        } `yaml:"global"`
        Modules        map[string]modules.ModuleConfig `yaml:"modules"` 
}


func ParseConfig(filename string) *Config {
	cfg := new(Config)
	
	fb, err := ioutil.ReadFile(filename)
	must(err)

	must(yaml.Unmarshal(fb, cfg))
	return cfg
}
