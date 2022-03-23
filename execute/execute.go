package execute

import (
	"github.com/wangfeiping/dyson/config"
	"github.com/wangfeiping/dyson/exporter"
	"github.com/wangfeiping/log"
)

type Executor struct {
	config    *config.ExecutorConfig
	exporters []*exporter.Exporter
}

func NewExecutor(executorConfig *config.ExecutorConfig) *Executor {
	executor := &Executor{config: executorConfig}

	if len(executor.config.Exporter) > 0 {
		for _, exporterConfig := range executor.config.Exporter {
			exp := exporter.NewExporter("proposal", "proposal on the blockchain",
				[]string{"chain_id", "start", "end"})
			var msc []*config.ExporterMetricConfig
			metricConfig := &config.ExporterMetricConfig{
				Name: "proposal",
				Labels: []string{"testnet",
					"${voting_start_time}",
					"${voting_end_time)"},
				Value: "${proposal_id}"}
			msc = append(msc, metricConfig)

			exp.SetMetricConfigs(msc)
			executor.exporters = append(executor.exporters, exp)

			log.Debug("exporter: ", exporterConfig, "; metrics: ", len(msc))
		}
	}

	return executor
}

func (e *Executor) GetExporters() []*exporter.Exporter {
	return e.exporters
}

func (e *Executor) DoExport() {
	if len(e.exporters) > 0 {
		for _, exp := range e.exporters {
			exp.SetMetrics(exp.DoExport())
		}
	}
}
