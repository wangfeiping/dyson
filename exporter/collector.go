package exporter

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/wangfeiping/log"
)

var collector *dysonCollector

func init() {
	collector = &dysonCollector{}
}

type dysonCollector struct {
	exporters []*Exporter
	mux       sync.RWMutex
}

// Collector returns a collector
// which exports metrics about status code of network service response
func Collector() prometheus.Collector {
	return collector
}

// Describe returns all descriptions of the collector.
func (c *dysonCollector) Describe(ch chan<- *prometheus.Desc) {
	c.mux.RLock()
	defer c.mux.RUnlock()

	for _, exp := range c.exporters {
		ch <- exp.GetDesc()
	}

	log.Debug("dysonCollector.Describe")
}

// Collect returns the current state of all metrics of the collector.
func (c *dysonCollector) Collect(ch chan<- prometheus.Metric) {
	c.mux.RLock()
	defer c.mux.RUnlock()

	for _, exp := range c.exporters {
		for _, metric := range exp.GetMetrics() {
			ch <- prometheus.MustNewConstMetric(
				exp.GetDesc(),
				prometheus.GaugeValue,
				metric.Value, metric.Labels...)
		}
	}

	log.Debug("dysonCollector.Collect")
}

func (c *dysonCollector) setExporters(exporters []*Exporter) {
	c.mux.Lock()
	defer c.mux.Unlock()

	c.exporters = exporters
}

// SetExporters set exporters
func SetExporters(exporters []*Exporter) {
	collector.setExporters(exporters)
}
