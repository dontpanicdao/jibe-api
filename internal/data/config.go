package data

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

var (
	Conf *Config
)

type Config struct {
	JibeUser     string `yaml:"jibeUser"`
	JibeHost     string `yaml:"jibeHost"`
	JibeDB       string `yaml:"jibeDB"`
	JibePassword string `yaml:"jibePassword"`
}

func InitConfig() {
	yamlFile, err := ioutil.ReadFile("conf.yaml")
	if err != nil {
		log.Panic(err.Error())
	}
	err = yaml.Unmarshal(yamlFile, &Conf)
	if err != nil {
		log.Panic(err.Error())
	}
}
