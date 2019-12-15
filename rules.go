package main

import (
	"errors"
	"fmt"

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
	Regex string `yaml:"regex" json:"regex"`
}

func (r *rules) load(path string) error {
	rulesData := readFile(path)
	err := yaml.Unmarshal([]byte(rulesData), &r)
	if err != nil {
		log.Error(err)
	}
	return err
}

func ensureLabelsContainRules(labels map[string]interface{}) error {
	for _, rule := range R.Rules {
		// Ensure labels contains rule
		if _, ok := labels[rule.Key]; ok {
			// If rule is found, match regex
		} else {
			// If rule is not found, reject
			errStr := fmt.Sprintf("%v not in labels", rule.Key)
			return errors.New(errStr)
		}
	}
	return nil
}
