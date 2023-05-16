package main

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v3"
)

type conf struct {
	File string `yaml:"file"`
}

func (c *conf) ReadConfigFile() *conf {

	yamlFile, err := ioutil.ReadFile("repono.yaml")
	if err != nil {
		log.Printf("yamlFile.Get err #%v ", err)
	}

	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal %v", err)
	}

	return c
}
