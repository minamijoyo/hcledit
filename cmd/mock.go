package cmd

import (
	"bytes"

	"github.com/spf13/cobra"
)

// newMockCmd is a helper function which returns a *cobra.Command
// whose in/out/err streams are mocked for testing.
func newMockCmd(runE func(cmd *cobra.Command, args []string) error, input string) *cobra.Command {
	cmd := &cobra.Command{
		RunE: runE,
	}

	inStream := bytes.NewBufferString(input)
	outStream := new(bytes.Buffer)
	errStream := new(bytes.Buffer)

	cmd.SetIn(inStream)
	cmd.SetOut(outStream)
	cmd.SetErr(errStream)

	return cmd
}

// mockErr is a helper function which returns a string written to mocked err stream.
func mockErr(cmd *cobra.Command) string {
	return cmd.ErrOrStderr().(*bytes.Buffer).String()
}

// mockOut is a helper function which returns a string written to mocked out stream.
func mockOut(cmd *cobra.Command) string {
	return cmd.OutOrStdout().(*bytes.Buffer).String()
}
