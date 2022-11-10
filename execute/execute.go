package execute

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/yalp/jsonpath"

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
		for _, exporterStr := range executor.config.Exporter {
			exporterConfig, metricConfig, err := parseExporter(exporterStr)
			if err != nil {
				log.Error(err.Error())
				return nil
			}
			exp := exporter.NewExporter(exporterConfig.Name,
				exporterConfig.Help,
				exporterConfig.Labels)
			var msc []*config.ExporterMetricConfig
			msc = append(msc, metricConfig)

			exp.SetMetricConfigs(msc)
			executor.exporters = append(executor.exporters, exp)

			log.Debug("exporter: ", exporterConfig, "; metrics: ", len(msc))
		}
	}

	return executor
}

func (e *Executor) Execute() {
	cache := config.GetCache()

	log.Info("Command: ", e.config.Command)
	// params := strings.Split(e.config.Command, " ")
	// cmd := exec.Command(params[0], params[1:]...)

	cmd := exec.Command("bash", "-c", e.config.Command)

	// fmt.Println("exec ", cmd.Args)
	// StdoutPipe方法返回一个在命令Start后与命令标准输出关联的管道。
	// Wait方法获知命令结束后会关闭这个管道，一般不需要显式的关闭该管道。
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Error("Command stdout pipe err: ", err)
		return
	}

	cmd.Stderr = os.Stderr
	// cmd.Dir = dir
	err = cmd.Start()
	// err = cmd.Run()
	if err != nil {
		log.Error("Command start err: ", err)
		return
	}

	//创建一个流来读取管道内内容这里逻辑是通过一行一行的读取的
	reader := bufio.NewReader(stdout)
	//实时循环读取输出流中的一行内容
	var jsonStr string
	for {
		jsonStr, err = reader.ReadString('\n')
		if err != nil || io.EOF == err {
			log.Debug("Command stdout EOF!")
			break
		}
		log.Debug("Result: ", jsonStr)

		for _, parser := range e.config.Parser {
			var name string
			if strings.ContainsRune(parser, '=') {
				i := strings.Index(parser, "=")
				name = parser[:i]
				parser = parser[i+1:]
			} else {
				i := strings.LastIndex(parser, ".")
				name = parser[i+1:]
			}
			log.Debug("Parser: ", parser)
			filter, err := jsonpath.Prepare(parser)
			if err != nil {
				log.Error("jsonpath prepare err: ", err)
				return
			}
			var data interface{}
			if err = json.Unmarshal([]byte(jsonStr), &data); err != nil {
				log.Error(err)
				return
			}
			out, err := filter(data)
			if err != nil {
				log.Error("jsonpath filter err: ", err)
				return
			}
			cache.Put(name, fmt.Sprint(out))
			log.Debug("Cached name: ", name, "; value: ", cache.Get(name))
		}
	}

	err = cmd.Wait()
	if err != nil {
		log.Error("Command wait err: ", err)
		return
	}
}

func (e *Executor) GetExporters() []*exporter.Exporter {
	return e.exporters
}

func (e *Executor) Export() {
	if len(e.exporters) > 0 {
		for _, exp := range e.exporters {
			exp.SetMetrics(exp.DoExport())
		}
	}
}

func parseExporter(exporterStr string) (*config.ExporterConfig, *config.ExporterMetricConfig, error) {
	i := strings.Index(exporterStr, "{")
	if i < 0 {
		return nil, nil, errors.New("exporter config error")
	}
	headers := strings.Split(exporterStr[0:i], ":")
	name := strings.Trim(headers[0], " ")
	if len(name) == 0 {
		return nil, nil, errors.New("exporter config error: no name")
	}
	exporterConfig := &config.ExporterConfig{
		Name: name}
	if len(headers) > 1 {
		exporterConfig.Help = headers[1]
	}

	metricConfig := &config.ExporterMetricConfig{
		Name: exporterConfig.Name}

	exporterStr = exporterStr[i:]
	i = strings.LastIndex(exporterStr, "${")
	labelStr := exporterStr[0:i]
	valueStr := exporterStr[i:]

	metricConfig.Value = strings.Trim(valueStr, " ")

	exporterConfig.Labels, metricConfig.Labels = parseLabels(labelStr)

	log.Debug("exporter config: name=", exporterConfig.Name)
	log.Debug("exporter config: help=", exporterConfig.Help)
	log.Debug("exporter config: labels=", exporterConfig.Labels)

	log.Debug("metric config: name=", metricConfig.Name)
	log.Debug("metric config: labels=", metricConfig.Labels)
	log.Debug("metric config: value=", metricConfig.Value)
	return exporterConfig, metricConfig, nil
}

func parseLabels(labelStr string) ([]string, []string) {
	var exporterLabels []string
	var metricLabels []string
	labelStr = strings.Trim(labelStr, " {}")
	labels := strings.Split(labelStr, ",")
	for _, labelStr := range labels {
		label := strings.Split(labelStr, ":")
		exporterLabels = append(exporterLabels, strings.Trim(label[0], " \""))
		metricLabels = append(metricLabels, strings.Trim(label[1], " \""))
	}
	return exporterLabels, metricLabels
}
