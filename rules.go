package main

import (
	"errors"
	"fmt"
	"regexp"

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

func validRuleRegex(r rule) error {
	_, err := regexp.Compile(r.Value.Regex)
	if err != nil {
		errStr := fmt.Sprintf("Rule: %v contains invalid regex", r.Name)
		return errors.New(errStr)
	}
	return nil
}

func validateAllRulesRegex(r rules) []string {
	// To send back every rule that has invalid regex
	// Instead of just one at a time
	var errArr []string
	for _, rule := range r.Rules {
		_, err := regexp.Compile(rule.Value.Regex)
		if err != nil {
			errStr := fmt.Sprintf("Rule: %v contains invalid regex", rule.Name)
			log.Errorf(errStr)
			errArr = append(errArr, errStr)
		}
	}
	if len(errArr) > 0 {
		return errArr
	}
	return nil
}

func (r *rules) load(path string) error {
	rulesData := readFile(path)
	err := yaml.Unmarshal([]byte(rulesData), &r)
	if err != nil {
		log.Error(err)
	}
	validateAllRulesRegex(*r)
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

func ensureLabelsMatchRules(labels map[string]interface{}) error {

	// Test label values against rules

	return nil
}
