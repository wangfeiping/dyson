package config

type ExporterConfig struct {
	Name   string
	Help   string
	Labels []string
}

type ExporterMetricConfig struct {
	Name   string
	Labels []string
	Value  string
	// valueType string
}
