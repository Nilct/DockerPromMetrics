package main

import (
	"strings"
	"github.com/prometheus/client_golang/prometheus"
)

func addMetrics() map[string]*prometheus.GaugeVec {
	gaugeVecs := make(map[string]*prometheus.GaugeVec)

	// Stack Metrics
	gaugeVecs["imageDuration"] = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "codingame",
			Subsystem: "my_computer",
			Name:      "docker_image_duration",
			Help:      "Docker image existence duration (in be decided)",
		}, []string{"repository", "tag"})
	gaugeVecs["stacksState"] = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "codingame",
			Subsystem: "my_computer",
			Name:      "docker_image_size",
			Help:      "Docker image size (in Mo)",
		}, []string{"repository", "tag"})
	return gaugeVecs
}

// setMetrics - Logic to set the state of a system as a gauge metric
func (e *Exporter) setMetrics() error {

	imageinfo:= listEvents()

	for _, ii:= range imageinfo {
		e.gaugeVecs["imageDuration"].With(prometheus.Labels{"repository": repository, "tag": tag}).Set(float64(ii.rawDuration))
		e.gaugeVecs["imageSize"].With(prometheus.Labels{"repository": repository, "tag": tag}).Set(float64(in.rawSize / 1000000))

	}
	return nil
}
