package commands

import (
	"github.com/spf13/cobra"

	"github.com/wangfeiping/dyson/config"
)

// NewVersionCommand returns version command
func NewVersionCommand(run Runner) *cobra.Command {
	cmd := &cobra.Command{
		Use:   config.CmdVersion,
		Short: "Show version info",
		RunE: func(cmd *cobra.Command, args []string) error {
			run()
			return nil
		},
	}
	return cmd
}
