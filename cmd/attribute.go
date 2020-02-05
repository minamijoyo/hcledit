package cmd

import (
	"fmt"

	"github.com/minamijoyo/hcledit/hclwritex"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(newAttributeCmd())
}

func newAttributeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "attribute",
		Short: "Edit attribute",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	cmd.AddCommand(
		newAttributeGetCmd(),
	)

	return cmd
}

func newAttributeGetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <NAME>",
		Short: "Get attribute",
		Long: `Get matched attributes at a given address

Arguments:
  NAME             An address of attribute to get.
`,
		RunE: runAttributeGetCmd,
	}

	return cmd
}

func runAttributeGetCmd(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("expected 1 argument, but got %d arguments", len(args))
	}

	address := args[0]

	return hclwritex.GetAttribute(cmd.InOrStdin(), cmd.OutOrStdout(), "-", address)
}
