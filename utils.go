package main

import (
	"io/ioutil"
	"regexp"

	log "github.com/sirupsen/logrus"
)

func readFile(path string) []byte {
	data, fileErr := ioutil.ReadFile(path)
	if fileErr != nil {
		log.Fatal(fileErr)
	}
	return data
}

func str2bool(s string) bool {
	return s != ""
}

func validRuleRegex(r rule) bool {
	_, err := regexp.Compile(r.Value.Regex)
	return err != nil
}

func validateAllRulesRegex(r rules) {
	for _, rule := range r.Rules {
		_, err := regexp.Compile(rule.Value.Regex)
		if err != nil {
			log.Errorf("Rule: %v contains invalid regex", rule.Name)
		}
	}
}
