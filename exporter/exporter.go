package exporter

import (
	"strconv"
	"sync"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/wangfeiping/dyson/config"
	"github.com/wangfeiping/log"
)

type ExporterMetric struct {
	Name   string
	Labels []string
	Value  float64
	// valueType string
}

type Exporter struct {
	desc          *prometheus.Desc
	metricConfigs []*config.ExporterMetricConfig
	metrics       []*ExporterMetric

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

func (e *Exporter) SetMetricConfigs(metrics []*config.ExporterMetricConfig) {
	e.metricConfigs = metrics
}

func (e *Exporter) DoExport() []*ExporterMetric {
	cache := config.GetCache()

	var metrics []*ExporterMetric
	// var value int
	value, err := strconv.Atoi(cache.Get("proposal_id"))
	if err != nil {
		log.Error("Convert value err: ", err)
		return metrics
	}
	metric := &ExporterMetric{
		Name: "proposal",
		Labels: []string{"testnet",
			cache.Get("voting_start_time"),
			cache.Get("voting_end_time")},
		Value: float64(value)}
	metrics = append(metrics, metric)
	return metrics
}
