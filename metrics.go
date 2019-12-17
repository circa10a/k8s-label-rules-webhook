package main

import ginprometheus "github.com/zsais/go-gin-prometheus"

func newRegistry(name string) *ginprometheus.Prometheus {
	return ginprometheus.NewPrometheus(name)
}
