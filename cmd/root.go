package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

	RootCmd.PersistentFlags().BoolP("debug", "", false, "Enable debug mode")
	viper.BindPFlag("debug", RootCmd.PersistentFlags().Lookup("debug"))
}

func setDefaultStream(cmd *cobra.Command) {
	cmd.SetIn(os.Stdin)
	cmd.SetOut(os.Stdout)
	cmd.SetErr(os.Stderr)
}
