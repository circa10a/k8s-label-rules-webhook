package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialise gin engine
	r := gin.Default()
	// Create prometheus registry named "gin"
	p := NewPrometheusRegistry("gin")
	// Pass gin to inject prometheus middleware
	p.Use(r)
	// Initialise paths and handlers in routes.go
	routes(r)
	// Start web server on 8080
	r.Run(":8080")
}
