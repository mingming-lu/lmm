package config

import (
	"encoding/json"
	"io/ioutil"
)

type config struct {
	Domain struct {
		App     string `yaml:"app"`
		API     string `yaml:"api"`
		Image   string `yaml:"image"`
		Manager string `yaml:"manager"`
	}
}

type domain struct {
	App     string `json:"app"`
	API     string `json:"api"`
	Image   string `json:"image"`
	Manager string `json:"manager"`
}

var (
	DomainApp     string
	DomainAPI     string
	DomainImage   string
	DomainManager string
)

func Init(mode string) {
	conf := &config{}
	if mode == "dev" {
		conf = newConf("config_dev.json")
	} else if mode == "prod" {
		conf = newConf("config.json")
	} else {
		panic("unknown mode: " + mode)
	}
	DomainApp = conf.Domain.App
	DomainAPI = conf.Domain.API
	DomainImage = conf.Domain.Image
	DomainManager = conf.Domain.Manager
}

func newConf(filename string) *config {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	conf := config{}
	err = json.Unmarshal(b, &conf)
	if err != nil {
		panic(err)
	}
	return &conf
}
