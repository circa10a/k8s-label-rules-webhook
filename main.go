package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	// Test rules
	// rules := getRules("./rules.yaml")
	// log.Info(rules[0])
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
