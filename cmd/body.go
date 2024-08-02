package cmd

import (
	"fmt"

	"github.com/minamijoyo/hcledit/editor"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	RootCmd.AddCommand(newBodyCmd())
}

func newBodyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "body",
		Short: "Edit body",
		RunE: func(cmd *cobra.Command, _ []string) error {
			return cmd.Help()
		},
	}

	cmd.AddCommand(
		newBodyGetCmd(),
	)

	return cmd
}

func newBodyGetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <ADDRESS>",
		Short: "Get body",
		Long: `Get body of first matched block at a given address

Arguments:
  ADDRESS          An address of block to get.
`,
		RunE: runBodyGetCmd,
	}

	return cmd
}

func runBodyGetCmd(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("expected 1 argument, but got %d arguments", len(args))
	}

	address := args[0]
	file := viper.GetString("file")
	update := viper.GetBool("update")

	filter := editor.NewBodyGetFilter(address)
	c := newDefaultClient(cmd)
	return c.Edit(file, update, filter)
}
