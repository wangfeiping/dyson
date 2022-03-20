package commands

import (
	"context"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/wangfeiping/dyson/config"
	"github.com/wangfeiping/log"
)

// Runner is command call function
type Runner func() (context.CancelFunc, error)

// NewRootCommand returns root command
func NewRootCommand(versioner Runner) *cobra.Command {
	root := &cobra.Command{
		Use:   config.CmdRoot,
		Short: config.ShortDescription,
		Run: func(cmd *cobra.Command, args []string) {
			if viper.GetBool(config.FlagVersion) {
				versioner()
				return
			}
			cmd.Help()
		},
		PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
			if err := viper.BindPFlags(cmd.Flags()); err != nil {
				log.Error("bind flags error: ", err)
				return err
			}

			if strings.EqualFold(cmd.Use, config.CmdRoot) ||
				strings.EqualFold(cmd.Use, config.CmdVersion) {
				// doesn't need init config & log
				return nil
			}

			initConfig()

			if !strings.EqualFold(cmd.Use, config.CmdStart) {
				// doesn't need init log
				return nil
			}

			initLogger()

			return
		},
	}

	root.Flags().BoolP(config.FlagVersion, "v", false, "show version info")
	root.PersistentFlags().StringP(config.FlagConfig, "c", "./config.yml", "config file path")

	return root
}

func initConfig() error {
	viper.SetConfigFile(viper.GetString(config.FlagConfig))

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		// stderr, so if we redirect output to json file, this doesn't appear
		// fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
		config.Load()
	} else if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
		// ignore not found error, return other errors
		return err
	}

	return nil
}

func initLogger() {
	log.Config(log.RollingFileConfig())
	log.Infof("Starting at %s", getExecPath())
}

// getExecPath returns the execution path
func getExecPath() (execPath string) {
	file, _ := exec.LookPath(os.Args[0])
	execFile := filepath.Base(file)
	execPath, _ = filepath.Abs(file)
	if len(execPath) > 1 {
		rs := []rune(execPath)
		execPath = string(rs[0:(len(execPath) - len(execFile))])
	}
	return
}

func commandRunner(run Runner, isKeepRunning bool) error {
	cancel, err := run()
	if err != nil {
		log.Error("Running error: ", err.Error())
		return err
	}
	if isKeepRunning {
		keepRunning(func(sig os.Signal) {
			defer log.Flush()
			if cancel != nil {
				cancel()
			}
			log.Debug("Stopped by signal: ", sig)
		})
	}
	return nil
}

func keepRunning(callback func(sig os.Signal)) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)

	select {
	case s, ok := <-signals:
		log.Infof("System signal [%v] %t, trying to run callback...", s, ok)
		if !ok {
			break
		}
		if callback != nil {
			callback(s)
		}
		log.Flush()
		os.Exit(1)
	}
}
