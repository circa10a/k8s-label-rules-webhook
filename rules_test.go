package main

import (
	"io/ioutil"
	"os"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func createInvalidYamlFile(t *testing.T, path string) {
	// Create invalid yaml file
	invalidYaml := `
	rules:
	- name:`
	invalidYamlFile := []byte(invalidYaml)
	// Break test if error writing file
	fileErr := ioutil.WriteFile(path, invalidYamlFile, 0644)
	if fileErr != nil {
		t.Error("Error writing test yaml file")
	}
}

func createInvalidRulesFile(t *testing.T, path string) {
	// Create invalid yaml file
	invalidRules := rules{
		Rules: []rule{
			{
				Name: "require-phone-number",
				Key:  "phone-number",
				Value: value{
					Regex: "[",
				},
			},
		},
	}
	data, _ := yaml.Marshal(invalidRules)
	// Break test if error writing file
	fileErr := ioutil.WriteFile(path, data, 0644)
	if fileErr != nil {
		t.Error("Error writing test yaml file")
	}
}

func TestLoadValid(t *testing.T) {
	// Load valid yaml
	// New struct
	r := &rules{}
	// Init map
	r.CompiledRegexs = make(map[string]*regexp.Regexp)
	// Load valid yaml
	err := r.load("rules.yaml")
	assert.NoError(t, err)
}

func TestLoadInvalid(t *testing.T) {
	// Load invalid yaml
	// New struct
	r := &rules{}
	// Init map
	r.CompiledRegexs = make(map[string]*regexp.Regexp)
	// Create invalid yaml file
	invalidYamlFile := "test.yaml"
	createInvalidYamlFile(t, invalidYamlFile)
	// Clean file after test
	defer os.Remove(invalidYamlFile)
	// Ensure error due to invalid yaml
	assert.Error(t, r.load(invalidYamlFile))
}

func TestDefaultCompiledRegex(t *testing.T) {
	_, err := defaultCompiledRegex()
	if err != nil {
		t.Error("Failed compiling default regex")
	}
}

func TestCompileRegexValid(t *testing.T) {
	// New struct
	r := &rules{}
	// Init map
	r.CompiledRegexs = make(map[string]*regexp.Regexp)
	// Load valid yaml
	err := r.load("rules.yaml")
	if err != nil {
		t.Error("Error loading yaml")
	}
	// Ensure no rule errors
	assert.Equal(t, 0, len(r.compileRegex(false)), "Compiling regex should not return any errors")
}

func TestCompileRegexValidStore(t *testing.T) {
	// New struct
	r := &rules{}
	// Init map
	r.CompiledRegexs = make(map[string]*regexp.Regexp)
	// Load valid yaml
	err := r.load("rules.yaml")
	if err != nil {
		t.Error("Error loading yaml")
	}
	// Store in map
	r.compileRegex(true)
	_, phoneKeyPresent := r.CompiledRegexs["require-phone-number"]
	_, numberKeyPresent := r.CompiledRegexs["require-number"]
	assert.True(t, phoneKeyPresent)
	assert.True(t, numberKeyPresent)
}

func TestCompileRegexInvalid(t *testing.T) {
	invalidRulesFile := "test.yaml"
	// Create invalid yaml file
	createInvalidRulesFile(t, invalidRulesFile)
	// Clean file after test
	defer os.Remove(invalidRulesFile)
	// New struct
	r := &rules{}
	// Init map
	r.CompiledRegexs = make(map[string]*regexp.Regexp)
	err := r.load(invalidRulesFile)
	if err != nil {
		t.Error("Error loading yaml")
	}
	// Ensure 1 rule error
	assert.Equal(t, 1, len(r.compileRegex(false)), "Compiling regex should return 1 error")
}

func TestValidateAllRulesRegex(t *testing.T) {
	// New struct
	r := &rules{}
	// Init map
	r.CompiledRegexs = make(map[string]*regexp.Regexp)
	// Load valid yaml
	err := r.load("rules.yaml")
	if err != nil {
		t.Error("Error loading yaml")
	}
	// Ensure no rule errors
	assert.Equal(t, 0, len(r.validateAllRulesRegex()), "Compiling regex should not return any errors")
}

func TestEnsureLabelsContainRulesValid(t *testing.T) {
	// New struct
	r := &rules{}
	// Init map
	r.CompiledRegexs = make(map[string]*regexp.Regexp)
	// Load valid yaml
	err := r.load("rules.yaml")
	if err != nil {
		t.Error("Error loading yaml")
	}
	// Simulate labels from k8s request
	labels := make(map[string]interface{})
	labels["phone-number"] = "555-555-5555"
	labels["number"] = "0"
	assert.NoError(t, r.ensureLabelsContainRules(labels))
}

func TestEnsureLabelsContainRulesInvalid(t *testing.T) {
	// New struct
	r := &rules{}
	// Init map
	r.CompiledRegexs = make(map[string]*regexp.Regexp)
	// Load valid yaml
	err := r.load("rules.yaml")
	if err != nil {
		t.Error("Error loading yaml")
	}
	// Simulate labels from k8s request
	labels := make(map[string]interface{})
	labels["phone-number"] = "555-555-5555"
	labels["num"] = "0"
	assert.Error(t, r.ensureLabelsContainRules(labels))
}
