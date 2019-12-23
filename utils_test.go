package main

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadFile(t *testing.T) {
	data, _ := readFile("rules.yaml")
	assert.True(t, len(data) > 0, "Ensure file was read into byte array")
}

func TestStr2Bool(t *testing.T) {
	assert.True(t, str2bool("test"), "Converts nonempty string to true")
	assert.False(t, str2bool(""), "Converts empty string to false")
}

func TestErrToStr(t *testing.T) {
	assert.Equal(t, errToStr(errors.New("Test error")), "Test error", "Converts error to string")
	assert.Equal(t, errToStr(nil), "", "Converts nil error to empty string")
}

func TestGetEnv(t *testing.T) {
	testEnvVar := "TESTING_ENV_VAR"
	os.Setenv(testEnvVar, "test val")
	defer os.Unsetenv(testEnvVar)

	assert.Equal(t, getEnv("TESTING_ENV_VAR", "default"), "test val", "Ensure env var is returned")
	assert.Equal(t, getEnv("UNSET_TESTING_ENV_VAR", "default"), "default", "Ensure default value is returned")
}
