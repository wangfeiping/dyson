package config

// nolint
const (
	CmdRoot          = "dyson"
	CmdStart         = "start"
	CmdAdd           = "add"
	CmdCall          = "call"
	CmdConfig        = "config"
	CmdVersion       = "version"
	CmdHelp          = "help"
	ShortDescription = "Command wrapper capable of being invoked remotely or automatically executed"
)

// nolint
const (
	FlagLog      = "log"
	FlagConfig   = "config"
	FlagListen   = "listen"
	FlagAlias    = "alias"
	FlagURL      = "url"
	FlagBody     = "body"
	FlagMethod   = "method"
	FlagRegex    = "regex"
	FlegDuration = "duration"
	FlagVersion  = CmdVersion

	ConfigKey = "executor"
)
