package main

import ginprometheus "github.com/zsais/go-gin-prometheus"

// NewPrometheusRegistry create new prometheus registry
func NewPrometheusRegistry(registryName string) *ginprometheus.Prometheus {
	return ginprometheus.NewPrometheus(registryName)
}
