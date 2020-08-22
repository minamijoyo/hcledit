package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	// Version is version number which automatically set on build.
	Version = "0.1.0"
)

func init() {
	RootCmd.AddCommand(newVersionCmd())
}

func newVersionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print version",
		RunE:  runVersionCmd,
	}

	return cmd
}

func runVersionCmd(cmd *cobra.Command, args []string) error {
	_, err := fmt.Fprintln(cmd.OutOrStdout(), Version)
	return err
}
