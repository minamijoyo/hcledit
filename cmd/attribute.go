package cmd

import (
	"fmt"

	"github.com/minamijoyo/hcledit/editor"
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
		newAttributeSetCmd(),
	)

	return cmd
}

func newAttributeGetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <NAME>",
		Short: "Get attribute",
		Long: `Get matched attribute at a given address

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

	return editor.GetAttribute(cmd.InOrStdin(), cmd.OutOrStdout(), "-", address)
}

func newAttributeSetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set <NAME> <VALUE>",
		Short: "Set attribute",
		Long: `Set a value of matched attribute at a given address

Arguments:
  NAME             An address of attribute to set.
  VALUE            A new value of attribute.
                   The value is set literally, even if references or expressions.
                   Thus, if you want to set a string literal "hoge", be sure to
                   escape double quotes so that they are not discarded by your shell.
                   e.g.) hcledit attribute set aaa.bbb.ccc '"hoge"'
`,
		RunE: runAttributeSetCmd,
	}

	return cmd
}

func runAttributeSetCmd(cmd *cobra.Command, args []string) error {
	if len(args) != 2 {
		return fmt.Errorf("expected 2 argument, but got %d arguments", len(args))
	}

	address := args[0]
	value := args[1]

	return editor.SetAttribute(cmd.InOrStdin(), cmd.OutOrStdout(), "-", address, value)
}
