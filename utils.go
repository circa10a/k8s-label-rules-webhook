package main

import (
	"io/ioutil"
	"os"
	"strconv"
)

func readFile(path string) ([]byte, error) {
	data, fileErr := ioutil.ReadFile(path)
	return data, fileErr
}

func str2bool(s string) bool {
	if b, err := strconv.ParseBool(s); err == nil {
		return b
	}
	return s != ""
}

// Remove nil possibility
func errToStr(e error) string {
	if e != nil {
		return e.Error()
	}
	return ""
}

// Set default value for environment variable if not found
func getEnv(key, def string) string {
	val := os.Getenv(key)
	if len(val) == 0 {
		return def
	}
	return val
}
