package commands

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/yalp/jsonpath"

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
	cache := config.GetCache()
	execs := config.GetAll()

	for _, exe := range execs {
		log.Info("Command: ", exe.Command)
		params := strings.Split(exe.Command, " ")
		cmd := exec.Command(params[0], params[1:]...)
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
		for {
			jsonStr, err2 := reader.ReadString('\n')
			if err2 != nil || io.EOF == err2 {
				log.Debug("Command stdout EOF!")
				break
			}
			log.Debug("Result: ", jsonStr)

			for _, parser := range exe.Parser {
				log.Debug("Parser: ", parser)
				// if strings.ContainsRune(parser, '=') {
				i := strings.LastIndex(parser, ".")
				name := parser[i+1:]
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

	// do export
	doExport()

	cache.Clear()
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

func doExport() {
	for _, executor := range executors {
		executor.DoExport()
	}
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

	cmd.Flags().Int64P(config.FlegDuration, "d", 30, "cycle time of the execute task")
	cmd.Flags().StringP(config.FlagListen, "l", ":9900", "listening address(ip:port) of exporter")
	return cmd
}
