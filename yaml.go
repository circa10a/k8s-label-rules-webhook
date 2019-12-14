package main

import (
	"io/ioutil"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

func readYamlFile(path string) []byte {
	data, fileErr := ioutil.ReadFile(path)
	if fileErr != nil {
		log.Fatal(fileErr)
	}
	return data
}

func getRules(path string) []rule {
	rulesData := readYamlFile(path)
	r := rules{}
	err := yaml.Unmarshal([]byte(rulesData), &r)
	if err != nil {
		log.Fatal(err)
	}
	return r.Rules
}
