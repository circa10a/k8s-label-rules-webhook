package main

import (
	"io/ioutil"
	"os"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
	testYamlFile := "test.yaml"
	// Clean file after test
	defer os.Remove(testYamlFile)
	invalidYaml := `
	rules:
	- name:`
	invalidYamlFile := []byte(invalidYaml)
	// Break test if error writing file
	fileErr := ioutil.WriteFile(testYamlFile, invalidYamlFile, 0644)
	if fileErr != nil {
		t.Error("Error writing test yaml file")
	}
	// Ensure error due to invalid yaml
	assert.Error(t, r.load(testYamlFile))
}
