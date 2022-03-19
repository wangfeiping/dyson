package config

import (
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"

	"github.com/wangfeiping/log"
)

type Executor struct {
	Expression []string `json:"expr,omitempty" yaml:"expr,omitempty"`
	Command    string   `json:"command" yaml:"command"`
	Parser     []string `json:"parser,omitempty" yaml:"parser,omitempty"`
	Exporter   []string `json:"exporter,omitempty" yaml:"exporter,omitempty"`
}

var mux sync.RWMutex
var execs []*Executor

func Reload() {
	mux.Lock()
	defer mux.Unlock()

	if err := viper.UnmarshalKey(ConfigKey, &execs); err != nil {
		log.Errorf("Load config error: %v", err)
		return
	}
}

func GetAll() []*Executor {
	mux.RLock()
	defer mux.RUnlock()

	return execs
}

func Load() {
	Reload()

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		Reload()
		log.Info("Config file changed:", e.Name)
	})
}

func Save() {
	expr := []string{"asd", "asd2", "asd3"}
	command := "./test cli ..."
	parser := []string{"aaa", "bbb", "ccc"}
	exporter := []string{"eee"}

	exec := &Executor{
		Expression: expr,
		Command:    command,
		Parser:     parser,
		Exporter:   exporter}
	execs = append(execs, exec)

	v := viper.New()
	v.SetConfigFile("test.yml") //viper.GetString(FlagConfig))
	v.Set(ConfigKey, execs)
	err := v.WriteConfig()
	if err != nil {
		log.Errorf("Failed: write config file error: %v", err)
	} else {
		log.Info("Success: config file saved.")
	}
}

func Check(data string, val ...string) string {
	return val[0] + data[6:]
}
