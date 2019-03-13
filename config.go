package main 


import (
        "io/ioutil"
        "go3status/modules"
)

func ParseConfig(f string) {


}

type Config struct {
        Global	struct {
            Interval	string `yaml:"interval"`
            Color	string	`yaml:"color"`
        } `yaml:"global"`
        Modules	map[string]modules.ModuleConfig `yaml:"modules"` 
}


func ParseConfig(filename string) *Config {
	cfg := new(Config)
	
	fb, err := ioutil.ReadFile(filename)
	must(err)

	must(yaml.Unmarshal(fb, cfg))
	return cfg
}
