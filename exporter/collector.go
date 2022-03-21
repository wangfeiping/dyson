package exporter

import (
	"github.com/prometheus/client_golang/prometheus"

	"github.com/wangfeiping/log"
)

var collector *dysonCollector

func init() {
	collector = &dysonCollector{}
}

type dysonCollector struct {
}

// Collector returns a collector
// which exports metrics about status code of network service response
func Collector() prometheus.Collector {
	return collector
}

// Describe returns all descriptions of the collector.
func (c *dysonCollector) Describe(ch chan<- *prometheus.Desc) {
	log.Debug("dysonCollector.Describe")
}

// Collect returns the current state of all metrics of the collector.
func (c *dysonCollector) Collect(ch chan<- prometheus.Metric) {
	log.Debug("dysonCollector.Collect")
}
