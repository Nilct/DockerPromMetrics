package main

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

// see : https://github.com/infinityworksltd/prometheus-rancher-exporter/blob/master/exporter.go

// Exporter Sets up all the runtime and metrics
type Exporter struct {
	// can add authorisation, etc. 
	gaugeVecs  map[string]*prometheus.GaugeVec
}

// NewExporter creates the metrics we wish to monitor
func newExporter() *Exporter {
	gaugeVecs := addMetrics()
	return &Exporter{
		gaugeVecs:  gaugeVecs,
	}
}