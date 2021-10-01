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
	// FilePath String pointer of path to rules yaml file
	FilePath *string
	// TLS Enable TLS
	TLS *bool
	// TLSCert Path to certificate
	TLSCert *string
	// TLSKey Path to key
	TLSKey *string
	// TLSPort TLS listening port
	TLSPort *int
	// R main rules struct to hold current ruleset
	R rules
	// Version is used to output the version of the application
	Version string
	// G default gin engine
	G = gin.Default()
)

// Read flags from command line args and set defaults
func flags() {
	// --file arg
	FilePath = flag.String("file", defaultRulesFile, "Path to yaml file with ruleset")
	// --metrics arg
	metrics := flag.Bool("metrics", str2bool(os.Getenv("METRICS")), "Enable prometheus endpoint at /metrics")
	// --tls arg
	TLS = flag.Bool("tls", str2bool(os.Getenv("TLS_ENABLED")), "Enable TLS")
	// --tls-cert arg
	TLSCert = flag.String("tls-cert", os.Getenv("TLS_CERT"), "Path to TLS certificate")
	// --tls-key arg
	TLSKey = flag.String("tls-key", os.Getenv("TLS_KEY"), "Path to TLS key")
	// --tls-port arg
	TLSPort = flag.Int("tls-port", defaultTLSPort, "TLS listening port")

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
// @version 0.2.15
// @description A kubernetes webhook to standardize labels on resources

// @contact.name GitHub
// @contact.url https://github.com/circa10a/k8s-label-rules-webhook/

// @license.name MIT
// @license.url https://github.com/circa10a/k8s-label-rules-webhook/blob/master/LICENSE
func main() {
	// Output version of application
	log.Infof("Version: %v", Version)
	// Validate command line arguments
	flags()
	// Instantiate map to cache regex compilations in
	R.CompiledRegexs = make(map[string]*regexp.Regexp)
	// Load initial rules into memory
	err := R.load(*FilePath)
	if err != nil {
		log.Error(err)
	}
	// Initialize paths and handlers in routes.go
	routes(G)
	// Listen via https if TLS enabled
	if *TLS {
		err = G.RunTLS(fmt.Sprintf(":%v", *TLSPort), *TLSCert, *TLSKey)
		if err != nil {
			log.Fatal(err)
		}
	}
	// Else listen on http
	// Defaults to port 8080, can be overridden via PORT env var.
	// Example: export PORT=3000
	err = G.Run()
	if err != nil {
		log.Fatal(err)
	}
}
