package main

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/gavv/httpexpect"
)

func init() {
	// set defaults to global pointers for config
	rulesFile := "rules.yaml"
	FilePath = &rulesFile
	// Init map to store compiled regexs
	R.CompiledRegexs = make(map[string]*regexp.Regexp)
	// Load initial rules into memory
	R.load(*FilePath)
	// load hadnlers into gin engine
	routes(G)
}

func TestRulesEndpoint(t *testing.T) {
	//start()
	// run server using httptest
	server := httptest.NewServer(G)
	defer server.Close()

	// create httpexpect instance
	e := httpexpect.New(t, server.URL)
	// is it working?
	response := e.GET("/rules").
		Expect().
		Status(http.StatusOK).JSON().Array()

	response.Length().Equal(2)
	// 1st rule
	response.Element(0).Object().ValueEqual("name", "require-phone-number")
	response.Element(0).Object().ValueEqual("key", "phone-number")
	response.Element(0).Object().Value("value").Object().ValueEqual("regex", "[0-9]{3}-[0-9]{3}-[0-9]{4}")
	// 2nd rule
	response.Element(1).Object().ValueEqual("name", "require-number")
	response.Element(1).Object().ValueEqual("key", "number")
	response.Element(1).Object().Value("value").Object().ValueEqual("regex", "[0-1]{1}")
}

func TestReloadEndpoint(t *testing.T) {
	//start()
	// run server using httptest
	server := httptest.NewServer(G)
	defer server.Close()

	// create httpexpect instance
	e := httpexpect.New(t, server.URL)
	// is it working?
	response := e.POST("/reload").
		Expect().
		Status(http.StatusOK).JSON()

	response.Object().ValueEqual("reloaded", true)
}
