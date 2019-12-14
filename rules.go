package main

import (
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

// Rules array from yaml
type rules struct {
	Rules []rule `yaml:"rules" json:"rules"`
}

// Individual rule within rules array
type rule struct {
	Name  string `yaml:"name" json:"name"`
	Key   string `yaml:"key" json:"key"`
	Value value  `yaml:"value" json:"value"`
}

// Value struct within each rule
type value struct {
	DataType string `yaml:"type" json:"type"`
	Regex    string `yaml:"regex" json:"regex"`
}

func (r *rules) load(path string) error {
	rulesData := readFile(path)
	err := yaml.Unmarshal([]byte(rulesData), &r)
	if err != nil {
		log.Error(err)
	}
	return err
}
