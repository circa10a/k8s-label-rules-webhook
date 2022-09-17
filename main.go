package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

const (
	// Default filename
	defaultRulesFile string = "rules.yaml"
	// Default TLS port
	defaultTLSPort int = 8443
)

var (
	// filePath String pointer of path to rules yaml file
	filePath *string
	// tls Enable TLS
	tls *bool
	// tlsCert Path to certificate
	tlsCert *string
	// tlsKey Path to key
	tlsKey *string
	// tlsPort TLS listening port
	tlsPort *int
	// r main rules struct to hold current ruleset
	r rules
	// version is used to output the version of the application
	version string
	// G default gin engine
	g = gin.Default()
)

// Read flags from command line args and set defaults
func flags() {
	// --file arg
	filePath = flag.String("file", defaultRulesFile, "Path to yaml file with ruleset")
	// --metrics arg
	metrics := flag.Bool("metrics", str2bool(os.Getenv("METRICS")), "Enable prometheus endpoint at /metrics")
	// --tls arg
	tls = flag.Bool("tls", str2bool(os.Getenv("TLS_ENABLED")), "Enable TLS")
	// --tls-cert arg
	tlsCert = flag.String("tls-cert", os.Getenv("TLS_CERT"), "Path to TLS certificate")
	// --tls-key arg
	tlsKey = flag.String("tls-key", os.Getenv("TLS_KEY"), "Path to TLS key")
	// --tls-port arg
	tlsPort = flag.Int("tls-port", defaultTLSPort, "TLS listening port")

	flag.Parse()
	// Input file validation
	if *filePath == "" {
		flag.PrintDefaults()
		log.Fatal("No file provided")
	}
	// Metrics flag validation
	if *metrics {
		// Create prometheus registry named "gin"
		p := newRegistry("gin")
		// Pass gin to inject prometheus middleware
		p.Use(g)
	}
}

// @title k8s-label-rules-webhook
// @description A kubernetes webhook to standardize labels on resources

// @contact.name GitHub
// @contact.url https://github.com/circa10a/k8s-label-rules-webhook/

// @license.name MIT
// @license.url https://github.com/circa10a/k8s-label-rules-webhook/blob/main/LICENSE
func main() {
	// Output version of application
	log.Infof("Version: %s", version)
	// Validate command line arguments
	flags()
	// Instantiate map to cache regex compilations in
	r.compiledRegexs = make(map[string]*regexp.Regexp)
	// Load initial rules into memory
	err := r.load(*filePath)
	if err != nil {
		log.Error(err)
	}
	// Initialize paths and handlers in routes.go
	routes(g)
	// Listen via https if TLS enabled
	if *tls {
		err = g.RunTLS(fmt.Sprintf(":%d", *tlsPort), *tlsCert, *tlsKey)
		if err != nil {
			log.Fatal(err)
		}
	}
	// Else listen on http
	// Defaults to port 8080, can be overridden via PORT env var.
	// Example: export PORT=3000
	err = g.Run()
	if err != nil {
		log.Fatal(err)
	}
}
