package command

import (
	"fmt"
	"os"
	"strings"

	"github.com/minamijoyo/hcledit/hclwritex"
	flag "github.com/spf13/pflag"
)

// BlockGetCommand is a command which gets blocks matching a given address.
type BlockGetCommand struct {
	Meta
	address string
}

// Run runs the procedure of this command.
func (c *BlockGetCommand) Run(args []string) int {
	cmdFlags := flag.NewFlagSet("block get", flag.ContinueOnError)

	if err := cmdFlags.Parse(args); err != nil {
		c.UI.Error(fmt.Sprintf("failed to parse arguments: %s", err))
		return 1
	}

	if len(cmdFlags.Args()) != 1 {
		c.UI.Error(fmt.Sprintf("The command expects 1 argument, but got %d", len(cmdFlags.Args())))
		c.UI.Error(c.Help())
		return 1
	}

	c.address = cmdFlags.Arg(0)
	err := hclwritex.GetBlocks(os.Stdin, os.Stdout, "-", c.address)
	if err != nil {
		c.UI.Error(fmt.Sprintf("failed to get blocks: %s", err))
		c.UI.Error(c.Help())
		return 1
	}

	return 0
}

// Help returns long-form help text.
func (c *BlockGetCommand) Help() string {
	helpText := `
Usage: hcledit block get [options] <ADDRESS>

Arguments
  ADDRESS            An address of blocks to get.
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns one-line help text.
func (c *BlockGetCommand) Synopsis() string {
	return "Get blocks matching a given address"
}
