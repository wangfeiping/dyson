package main

import (
	"context"
	"fmt"

	"github.com/wangfeiping/dyson/config"
)

// nolint
var (
	Version   = "0.0.0"
	GitCommit string
	GoVersion string
	BuidDate  string
)

var versioner = func() (context.CancelFunc, error) {

	s := `Dyson - %s
version:	%s
revision:	%s
compile:	%s
go version:	%s
`

	fmt.Printf(s, config.ShortDescription, Version, GitCommit, BuidDate, GoVersion)

	return nil, nil
}
