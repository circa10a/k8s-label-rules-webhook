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
	Rules          []rule `yaml:"rules" json:"rules"`
	CompiledRegexs map[string]*regexp.Regexp
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

func (r *rules) load(path string) error {
	rulesData, _ := readFile(path)
	err := yaml.Unmarshal([]byte(rulesData), &r)
	if err != nil {
		log.Error(err)
	}
	r.compileRegex(true)
	// Should return nil
	return err
}

// In case of invalid regex provided
func defaultCompiledRegex() (*regexp.Regexp, error) {
	wildcard, err := regexp.Compile(".*")
	return wildcard, err
}

// Insert into map to prevent recompiling for every call
func (r *rules) compileRegex(storeInMap bool) []ruleError {
	var errArr []ruleError
	for _, rule := range r.Rules {
		compiled, err := regexp.Compile(rule.Value.Regex)
		if err != nil {
			log.Errorf("rule: %v, err: %v", rule.Name, err.Error())
			errArr = append(errArr, ruleError{RuleName: rule.Name, Err: err.Error()})
			// In the event of invalid regex, anything for that rule is allowed
			defaultRegex, _ := defaultCompiledRegex()
			r.CompiledRegexs[rule.Name] = defaultRegex
		} else {
			// Store user supplied compiled regex if no error
			if storeInMap {
				// Update/Store in map
				r.CompiledRegexs[rule.Name] = compiled
			}
		}
	}
	// If any regex compilation errors detected, return
	if len(errArr) > 0 {
		return errArr
	}
	// No errors
	return nil
}

func (r *rules) validateAllRulesRegex() []ruleError {
	// To send back every rule that has invalid regex
	return r.compileRegex(false)
}

func (r *rules) ensureLabelsContainRules(labels map[string]interface{}) error {
	for _, rule := range r.Rules {
		// Ensure labels contains rule
		if _, ok := labels[rule.Key]; !ok {
			// If rule is not found, reject
			errStr := fmt.Sprintf("%v not in labels", rule.Key)
			return errors.New(errStr)
		}
	}
	return nil
}

func (r *rules) ensureLabelsMatchRules(labels map[string]interface{}) error {
	for _, rule := range r.Rules {
		// Force all values to strings to prevent panic from interface conversion
		labelVal := fmt.Sprintf("%v", labels[rule.Key])
		regex, _ := r.CompiledRegexs[rule.Name]
		if !regex.MatchString(labelVal) {
			errStr := fmt.Sprintf("Value for label '%v' does not match expression '%v'", rule.Key, rule.Value.Regex)
			return errors.New(errStr)
		}
	}
	return nil
}
