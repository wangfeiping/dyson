package commands

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"

	"github.com/wangfeiping/dyson/config"
	"github.com/wangfeiping/log"
)

var starter = func() (cancel context.CancelFunc, err error) {
	log.Info("Start...")

	doJob()

	return
}

func doJob() {
	config.Load()
	execs := config.GetAll()

	for _, exe := range execs {
		log.Info("execute: ", exe.Command)
		params := strings.Split(exe.Command, " ")
		cmd := exec.Command(params[0], params[1:]...)
		// fmt.Println("exec ", cmd.Args)
		// StdoutPipe方法返回一个在命令Start后与命令标准输出关联的管道。
		// Wait方法获知命令结束后会关闭这个管道，一般不需要显式的关闭该管道。
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			log.Error("cmd.StdoutPipe err: ", err)
			return
		}

		cmd.Stderr = os.Stderr
		// cmd.Dir = dir
		err = cmd.Start()
		// err = cmd.Run()
		if err != nil {
			log.Error("cmd.Start err: ", err)
			return
		}

		//创建一个流来读取管道内内容这里逻辑是通过一行一行的读取的
		reader := bufio.NewReader(stdout)
		//实时循环读取输出流中的一行内容
		for {
			line, err2 := reader.ReadString('\n')
			if err2 != nil || io.EOF == err2 {
				log.Debug("Stdout EOF!")
				break
			}
			fmt.Println("Readed: ", line)
		}

		err = cmd.Wait()
		if err != nil {
			log.Error("cmd.Wait err: ", err)
		}
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

	cmd.Flags().Int64P(config.FlegDuration, "d", 30, "The cycle time of the watch task")
	cmd.Flags().StringP(config.FlagListen, "l", ":9900", "The listening address(ip:port) of exporter")
	return cmd
}
