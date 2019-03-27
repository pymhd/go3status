package main

import (
	"go3status/modules"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type ModuleConfigMap map[string]modules.ModuleConfig

type Config struct {
	Global struct {
		LogFile     string   `yaml:"logFile"`
		LogLevel    int    `yaml:"LogLevel"`
	} `yaml:"global"`
	Modules []ModuleConfigMap `yaml:"modules"`
}

func ParseConfig(filename string) *Config {
	cfg := new(Config)

	fb, err := ioutil.ReadFile(filename)
	must(err)

	must(yaml.Unmarshal(fb, cfg))
	return cfg
}
