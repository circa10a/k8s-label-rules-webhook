package main

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/gavv/httpexpect"
)

func start() {
	// set defaults to global pointers for config
	flags()
	// Init map to store compiled regexs
	R.CompiledRegexs = make(map[string]*regexp.Regexp)
	// Load initial rules into memory
	R.load("rules.yaml")
	// load hadnlers into gin engine
	routes(G)
}

func TestRulesEndpoint(t *testing.T) {
	start()
	// run server using httptest
	server := httptest.NewServer(G)
	defer server.Close()

	// create httpexpect instance
	e := httpexpect.New(t, server.URL)
	// is it working?
	e.GET("/rules").
		Expect().
		Status(http.StatusOK).JSON().Array()
}
