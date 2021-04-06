package cmd

import (
	"fmt"

	"github.com/minamijoyo/hcledit/editor"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	RootCmd.AddCommand(newFmtCmd())
}

func newFmtCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "fmt",
		Short: "Format file",
		Long:  "Format a file to a caconical style",
		RunE:  runFmtCmd,
	}

	return cmd
}

func runFmtCmd(cmd *cobra.Command, args []string) error {
	if len(args) != 0 {
		return fmt.Errorf("expected no argument, but got %d arguments", len(args))
	}

	file := viper.GetString("file")
	update := viper.GetBool("update")

	filter := editor.NewFormatterFilter()
	c := newDefaultClient(cmd)
	return c.Edit(file, update, filter)
}
