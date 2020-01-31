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
		newBlockMvCmd(),
		newBlockListCmd(),
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

	return hclwritex.GetBlock(cmd.InOrStdin(), cmd.OutOrStdout(), "-", address)
}

func newBlockMvCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mv <FROM> <TO>",
		Short: "Move block (Rename block type and labels)",
		Long: `Move block (Rename block type and labels)

Arguments:
  FROM             An old address of block.
  TO               A new address of block.
`,
		RunE: runBlockMvCmd,
	}

	return cmd
}

func runBlockMvCmd(cmd *cobra.Command, args []string) error {
	if len(args) != 2 {
		return fmt.Errorf("expected 2 argument, but got %d arguments", len(args))
	}

	from := args[0]
	to := args[1]

	return hclwritex.RenameBlock(cmd.InOrStdin(), cmd.OutOrStdout(), "-", from, to)
}

func newBlockListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List block",
		RunE:  runBlockListCmd,
	}

	return cmd
}

func runBlockListCmd(cmd *cobra.Command, args []string) error {
	if len(args) != 0 {
		return fmt.Errorf("expected 0 argument, but got %d arguments", len(args))
	}

	return hclwritex.ListBlock(cmd.InOrStdin(), cmd.OutOrStdout(), "-")
}
