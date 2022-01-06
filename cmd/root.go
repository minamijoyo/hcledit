package cmd

import (
	"os"

	"github.com/minamijoyo/hcledit/editor"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// RootCmd is a top level command instance
var RootCmd = &cobra.Command{
	Use:           "hcledit",
	Short:         "A command line editor for HCL",
	SilenceErrors: true,
	SilenceUsage:  true,
}

func init() {
	// set global flags
	flags := RootCmd.PersistentFlags()
	flags.StringP("file", "f", "-", "A path of input file")
	flags.BoolP("update", "u", false, "Update files in-place")
	viper.BindPFlag("file", flags.Lookup("file"))     // nolint: errcheck
	viper.BindPFlag("update", flags.Lookup("update")) // nolint: errcheck

	setDefaultStream(RootCmd)
}

func setDefaultStream(cmd *cobra.Command) {
	cmd.SetIn(os.Stdin)
	cmd.SetOut(os.Stdout)
	cmd.SetErr(os.Stderr)
}

func newDefaultClient(cmd *cobra.Command) editor.Client {
	o := &editor.Option{
		InStream:  cmd.InOrStdin(),
		OutStream: cmd.OutOrStdout(),
		ErrStream: cmd.ErrOrStderr(),
	}
	return editor.NewClient(o)
}
