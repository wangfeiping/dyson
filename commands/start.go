package commands

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/wangfeiping/dyson/config"
	"github.com/wangfeiping/dyson/execute"
	"github.com/wangfeiping/dyson/exporter"
	"github.com/wangfeiping/log"
)

var executors []*execute.Executor
var exporters []*exporter.Exporter

var starter = func() (cancel context.CancelFunc, err error) {
	log.Info("Start...")

	running := true
	t := time.NewTicker(time.Duration(
		viper.GetInt64(config.FlegDuration)) * time.Second)
	var wg sync.WaitGroup
	cancel = func() {
		running = false
		t.Stop()
		wg.Wait()
		log.Info("Stopped.")
	}

	initExecutors()

	go func() {
		wg.Add(1)
		doJob()
		for running {
			select {
			case <-t.C:
				{
					doJob()
				}
			default:
				{
					time.Sleep(100 * time.Millisecond)
				}
			}
		}
		log.Info("Done")
		wg.Done()
	}()

	prometheus.MustRegister(exporter.Collector())

	http.Handle("/metrics", promhttp.Handler())
	listen := viper.GetString(config.FlagListen)
	err = http.ListenAndServe(listen, nil)
	log.Error(err)

	return
}

func doJob() {

	for _, executor := range executors {
		// execute
		executor.Execute()
		// export
		executor.Export()
	}

	config.GetCache().Clear()
}

func initExecutors() {
	execs := config.GetAll()

	for _, exe := range execs {
		executor := execute.NewExecutor(exe)
		exporters = append(exporters, executor.GetExporters()...)
		executors = append(executors, executor)
	}

	exporter.SetExporters(exporters)
}

// NewStartCommand 创建 start/服务启动 命令
func NewStartCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   config.CmdStart,
		Short: "Start",
		RunE: func(cmd *cobra.Command, args []string) error {
			return commandRunner(starter, true)
		},
	}

	cmd.Flags().Int64P(config.FlegDuration, "d", 3600, "cycle time of the execute task")
	cmd.Flags().StringP(config.FlagListen, "l", ":25559", "listening address(ip:port) of exporter")
	return cmd
}
