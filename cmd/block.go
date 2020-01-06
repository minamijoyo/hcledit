package cmd

import (
	"fmt"

	"github.com/minamijoyo/hcledit/hclwritex"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(newBlockCmd())
}

func newBlockCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "block",
		Short: "Edit block",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	cmd.AddCommand(
		newBlockGetCmd(),
	)

	return cmd
}

func newBlockGetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get block",
		RunE:  runBlockGetCmd,
	}

	return cmd
}

func runBlockGetCmd(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("expected 1 argument, but got %d arguments", len(args))
	}

	address := args[0]

	return hclwritex.GetBlocks(cmd.InOrStdin(), cmd.OutOrStdout(), "-", address)
}
