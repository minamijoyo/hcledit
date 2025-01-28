package cmd

import (
	"errors"
	"fmt"

	"github.com/minamijoyo/hcledit/editor"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	RootCmd.AddCommand(newBlockCmd())
}

func newBlockCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "block",
		Short: "Edit block",
		RunE: func(cmd *cobra.Command, _ []string) error {
			return cmd.Help()
		},
	}

	cmd.AddCommand(
		newBlockGetCmd(),
		newBlockMvCmd(),
		newBlockListCmd(),
		newBlockRmCmd(),
		newBlockAppendCmd(),
		newBlockNewCmd(),
	)

	return cmd
}

func newBlockGetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <ADDRESS>",
		Short: "Get block",
		Long: `Get matched blocks at a given address

Arguments:
  ADDRESS          An address of block to get.
`,
		RunE: runBlockGetCmd,
	}

	return cmd
}

func runBlockGetCmd(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("expected 1 argument, but got %d arguments", len(args))
	}

	address := args[0]
	file := viper.GetString("file")
	update := viper.GetBool("update")

	filter := editor.NewBlockGetFilter(address)
	c := newDefaultClient(cmd)
	return c.Edit(file, update, filter)
}

func newBlockMvCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mv <FROM_ADDRESS> <TO_ADDRESS>",
		Short: "Move block (Rename block type and labels)",
		Long: `Move block (Rename block type and labels)

Arguments:
  FROM_ADDRESS     An old address of block.
  TO_ADDRESS       A new address of block.
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
	file := viper.GetString("file")
	update := viper.GetBool("update")

	filter := editor.NewBlockRenameFilter(from, to)
	c := newDefaultClient(cmd)
	return c.Edit(file, update, filter)
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

	file := viper.GetString("file")
	update := viper.GetBool("update")
	if update {
		return errors.New("The update flag is not allowed")
	}

	sink := editor.NewBlockListSink()
	c := newDefaultClient(cmd)
	return c.Derive(file, sink)
}

func newBlockRmCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rm <ADDRESS>",
		Short: "Remove block",
		Long: `Remove matched blocks at a given address

Arguments:
  ADDRESS          An address of block to remove.
`,
		RunE: runBlockRmCmd,
	}

	return cmd
}

func runBlockRmCmd(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("expected 1 argument, but got %d arguments", len(args))
	}

	address := args[0]
	file := viper.GetString("file")
	update := viper.GetBool("update")

	filter := editor.NewBlockRemoveFilter(address)
	c := newDefaultClient(cmd)
	return c.Edit(file, update, filter)
}

func newBlockAppendCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "append <PARENT_ADDRESS> <CHILD_ADDRESS>",
		Short: "Append block",
		Long: `Append a new child block to matched blocks at a given parent block address

Arguments:
  PARENT_ADDRESS      A parent block address to be appended.
  CHILD_ADDRESS       A new child block relative address.
`,
		RunE: runBlockAppendCmd,
	}

	flags := cmd.Flags()
	flags.Bool("newline", false, "Append a new line before a new child block")
	_ = viper.BindPFlag("block.append.newline", flags.Lookup("newline"))

	return cmd
}

func runBlockAppendCmd(cmd *cobra.Command, args []string) error {
	if len(args) != 2 {
		return fmt.Errorf("expected 2 argument, but got %d arguments", len(args))
	}

	parent := args[0]
	child := args[1]
	newline := viper.GetBool("block.append.newline")
	file := viper.GetString("file")
	update := viper.GetBool("update")

	filter := editor.NewBlockAppendFilter(parent, child, newline)
	c := newDefaultClient(cmd)
	return c.Edit(file, update, filter)
}

func newBlockNewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "new <BLOCK_TYPE> <LABEL>...",
		Short: "Create a new block",
		Long: `Create a new block

Arguments:
  BLOCK_TYPE       A block type to be created.
  LABEL            A label of block to be created.
`,
		RunE: runBlockNewCmd,
	}
	return cmd
}

func runBlockNewCmd(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("expected at least 1 argument, but got %d arguments", len(args))
	}

	blockType := args[0]
	labels := args[1:]
	file := viper.GetString("file")
	update := viper.GetBool("update")

	filter := editor.NewBlockNewFilter(blockType, labels)
	c := newDefaultClient(cmd)
	return c.Edit(file, update, filter)
}
