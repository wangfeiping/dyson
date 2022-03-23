package exporter

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

type ExporterMetric struct {
	Name   string
	Labels []string
	Value  float64
	// valueType string
}

type Exporter struct {
	desc    *prometheus.Desc
	metrics []*ExporterMetric

	mux sync.RWMutex
}

func NewExporter(name string, help string, labels []string) *Exporter {
	exp := &Exporter{
		desc: prometheus.NewDesc(name, help, labels, nil)}
	return exp
}

func (e *Exporter) GetDesc() *prometheus.Desc {
	return e.desc
}

func (e *Exporter) GetMetrics() []*ExporterMetric {
	e.mux.RLock()
	defer e.mux.RUnlock()

	return e.metrics
}

func (e *Exporter) SetMetrics(metrics []*ExporterMetric) {
	e.mux.Lock()
	defer e.mux.Unlock()

	e.metrics = metrics
}
