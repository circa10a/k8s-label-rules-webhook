package main

import (
	"flag"
	"os"
	"regexp"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

const (
	// Default filename
	defaultRulesFile string = "rules.yaml"
)

var (
	// FilePath String pointer of path to rules yaml file
	FilePath *string
	// SwaggerAPIDocURLFlag Swagger URL Flag
	SwaggerAPIDocURLFlag *string
	// SwaggerAPIDocURLStr Swagger URL Flag
	SwaggerAPIDocURLStr string
	// R main rules struct to hold current ruleset
	R rules
	// G default gin engine
	G = gin.Default()
)

// Read flags from command line args and set defaults
func flags() {
	// --file arg
	FilePath = flag.String("file", defaultRulesFile, "Path to yaml file with ruleset")
	// --metrics arg
	metrics := flag.Bool("metrics", str2bool(os.Getenv("METRICS")), "Enable prometheus endpoint at /metrics")
	flag.Parse()
	// Input file validation
	if *FilePath == "" {
		flag.PrintDefaults()
		log.Fatal("No file provided")
	}
	// Metrics flag validation
	if *metrics {
		// Create prometheus registry named "gin"
		p := newRegistry("gin")
		// Pass gin to inject prometheus middleware
		p.Use(G)
	}
}

// @title k8s-label-rules-webhook
// @version 0.1.0
// @description A kubernetes webhook to standardize labels on resources

// @contact.name GitHub
// @contact.url https://github.com/circa10a/k8s-label-rules-webhook/

// @license.name MIT
// @license.url https://github.com/circa10a/k8s-label-rules-webhook/blob/master/LICENSE
func main() {
	// Validate command line arguments
	flags()
	// Instantiate map to cache regex compilations in
	R.CompiledRegexs = make(map[string]*regexp.Regexp)
	// Load initial rules into memory
	R.load(*FilePath)
	// Initialize paths and handlers in routes.go
	routes(G)
	// Start web server
	// Defaults to port 8080, can be overridden via PORT env var.
	// Example: export PORT=3000
	G.Run()
}
