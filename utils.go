package main

import (
	"io/ioutil"

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
