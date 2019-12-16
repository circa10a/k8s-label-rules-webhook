package main

import (
	"flag"
	"os"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	ginprometheus "github.com/zsais/go-gin-prometheus"
)

var (
	// FilePath String pointer of path to rules yaml file
	FilePath *string
	// R main rules struct to hold current ruleset
	R rules
	// G default gin engine
	G = gin.Default()
)

func flags() {
	// --file arg
	FilePath = flag.String("file", os.Getenv("FILE"), "Path to yaml file with ruleset")
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
		p := ginprometheus.NewPrometheus("gin")
		// Pass gin to inject prometheus middleware
		p.Use(G)
	}
}

func main() {
	// Validate command line arguments
	flags()
	// Load initial rules into memory
	R.load(*FilePath)
	// Initialize paths and handlers in routes.go
	routes(G)
	// Start web server
	// Defaults to port 8080, can be overridden via PORT env var.
	// Example: export PORT=3000
	G.Run()
}
