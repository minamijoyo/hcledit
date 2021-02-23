package cmd

import (
	"fmt"

	"github.com/minamijoyo/hcledit/editor"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
		newAttributeRmCmd(),
		newAttributeAppendCmd(),
	)

	return cmd
}

func newAttributeGetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <ADDRESS>",
		Short: "Get attribute",
		Long: `Get matched attribute at a given address

Arguments:
  ADDRESS          An address of attribute to get.
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
		Use:   "set <ADDRESS> <VALUE>",
		Short: "Set attribute",
		Long: `Set a value of matched attribute at a given address

Arguments:
  ADDRESS          An address of attribute to set.
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

	o := editor.NewEditOperator(editor.NewAttributeSetFilter(address, value))
	return o.Apply(cmd.InOrStdin(), cmd.OutOrStdout(), "-")
}

func newAttributeRmCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rm <ADDRESS>",
		Short: "Remove attribute",
		Long: `Remove a matched attribute at a given address

Arguments:
  ADDRESS          An address of attribute to remove.
`,
		RunE: runAttributeRmCmd,
	}

	return cmd
}

func runAttributeRmCmd(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("expected 1 argument, but got %d arguments", len(args))
	}

	address := args[0]

	o := editor.NewEditOperator(editor.NewAttributeRemoveFilter(address))
	return o.Apply(cmd.InOrStdin(), cmd.OutOrStdout(), "-")
}

func newAttributeAppendCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "append <ADDRESS> <VALUE>",
		Short: "Append attribute",
		Long: `Append a new attribute at a given address

Arguments:
  ADDRESS          An address of attribute to append.
  VALUE            A new value of attribute.
                   The value is set literally, even if references or expressions.
                   Thus, if you want to set a string literal "hoge", be sure to
                   escape double quotes so that they are not discarded by your shell.
                   e.g.) hcledit attribute append aaa.bbb.ccc '"hoge"'
`,
		RunE: runAttributeAppendCmd,
	}

	flags := cmd.Flags()
	flags.Bool("newline", false, "Append a new line before a new attribute")
	viper.BindPFlag("attribute.append.newline", flags.Lookup("newline"))

	return cmd
}

func runAttributeAppendCmd(cmd *cobra.Command, args []string) error {
	if len(args) != 2 {
		return fmt.Errorf("expected 2 argument, but got %d arguments", len(args))
	}

	address := args[0]
	value := args[1]
	newline := viper.GetBool("attribute.append.newline")

	o := editor.NewEditOperator(editor.NewAttributeAppendFilter(address, value, newline))
	return o.Apply(cmd.InOrStdin(), cmd.OutOrStdout(), "-")
}
