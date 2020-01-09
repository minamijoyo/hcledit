package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// RootCmd is a top level command instance
var RootCmd = &cobra.Command{
	Use:           "hcledit",
	Short:         "A stream editor for HCL",
	SilenceErrors: true,
	SilenceUsage:  true,
}

func init() {
	setDefaultStream(RootCmd)
}

func setDefaultStream(cmd *cobra.Command) {
	cmd.SetIn(os.Stdin)
	cmd.SetOut(os.Stdout)
	cmd.SetErr(os.Stderr)
}
