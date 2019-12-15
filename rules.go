package main

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

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

type ruleError struct {
	RuleName string `json:"rulename"`
	Err      string `json:"err"`
}

func errsToString(r []ruleError) string {
	var s []string
	for _, v := range r {
		errStr := fmt.Sprintf("rule: %v, err: %v", v.RuleName, v.Err)
		s = append(s, errStr)
	}
	return strings.Join(s, " ")
}

func (r *rules) load(path string) error {
	rulesData := readFile(path)
	err := yaml.Unmarshal([]byte(rulesData), &r)
	if err != nil {
		log.Error(err)
	}
	r.validateAllRulesRegex()
	return err
}

func (r *rules) validateAllRulesRegex() []ruleError {
	// To send back every rule that has invalid regex
	// Instead of just one at a time
	var errArr []ruleError
	for _, rule := range r.Rules {
		_, err := regexp.Compile(rule.Value.Regex)
		if err != nil {
			log.Errorf("rule: %v, err: %v", rule.Name, err.Error())
			errArr = append(errArr, ruleError{RuleName: rule.Name, Err: err.Error()})
		}
	}
	if len(errArr) > 0 {
		return errArr
	}
	return nil
}

func (r *rules) ensureLabelsContainRules(labels map[string]interface{}) error {
	for _, rule := range r.Rules {
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

func (r *rules) ensureLabelsMatchRules(labels map[string]interface{}) error {

	// Test label values against rules

	return nil
}
