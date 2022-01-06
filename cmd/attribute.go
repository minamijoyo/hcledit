package cmd

import (
	"errors"
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
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
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
	file := viper.GetString("file")
	update := viper.GetBool("update")
	if update {
		return errors.New("The update flag is not allowed")
	}

	sink := editor.NewAttributeGetSink(address)
	c := newDefaultClient(cmd)
	return c.Derive(file, sink)
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
	file := viper.GetString("file")
	update := viper.GetBool("update")

	filter := editor.NewAttributeSetFilter(address, value)
	c := newDefaultClient(cmd)
	return c.Edit(file, update, filter)
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
	file := viper.GetString("file")
	update := viper.GetBool("update")

	filter := editor.NewAttributeRemoveFilter(address)
	c := newDefaultClient(cmd)
	return c.Edit(file, update, filter)
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
	_ = viper.BindPFlag("attribute.append.newline", flags.Lookup("newline"))

	return cmd
}

func runAttributeAppendCmd(cmd *cobra.Command, args []string) error {
	if len(args) != 2 {
		return fmt.Errorf("expected 2 argument, but got %d arguments", len(args))
	}

	address := args[0]
	value := args[1]
	newline := viper.GetBool("attribute.append.newline")
	file := viper.GetString("file")
	update := viper.GetBool("update")

	filter := editor.NewAttributeAppendFilter(address, value, newline)
	c := newDefaultClient(cmd)
	return c.Edit(file, update, filter)
}
