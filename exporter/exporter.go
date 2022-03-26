package exporter

import (
	"strconv"
	"strings"
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

	var metrics []*ExporterMetric

	for _, mConfig := range e.metricConfigs {

		value, err := strconv.Atoi(checkVariable(mConfig.Value))
		if err != nil {
			log.Error("Convert value err: ", err)
			return metrics
		}

		var labelValues []string
		for _, lValue := range mConfig.Labels {
			lValue = checkVariable(lValue)
			labelValues = append(labelValues, lValue)
		}

		metric := &ExporterMetric{
			Name:   mConfig.Name,
			Labels: labelValues,
			Value:  float64(value)}
		metrics = append(metrics, metric)
	}
	return metrics
}

func checkVariable(str string) string {
	cache := config.GetCache()

	if strings.HasPrefix(str, "$") {
		value := strings.Trim(str, " ${}()")
		return cache.Get(value)
	}

	return str
}
