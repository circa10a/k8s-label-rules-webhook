package main

import (
	"flag"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func flags() {
	filePath := flag.String("file", "", "Path to yaml file with ruleset")
	flag.Parse()
	if *filePath == "" {
		log.Fatal("No file provided")
	}
}

func main() {
	// Validate command line arguments
	flags()
	// Test rules
	rules := getRules(*filePath)
	log.Info(rules[0])
	// match, _ := regexp.MatchString(rules[0].Value.Regex, "318-305-6964")
	// fmt.Println(match)
	// Initialize gin engine
	r := gin.Default()
	// Create prometheus registry named "gin"
	p := NewPrometheusRegistry("gin")
	// Pass gin to inject prometheus middleware
	p.Use(r)
	// Initialize paths and handlers in routes.go
	routes(r)
	// Start web server on 8080
	r.Run(":8080")
}
