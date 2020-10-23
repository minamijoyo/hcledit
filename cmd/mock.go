package cmd

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/spf13/cobra"
)

// newMockCmd is a helper function which returns a *cobra.Command
// whose in/out/err streams are mocked for testing.
func newMockCmd(cmd *cobra.Command, input string) *cobra.Command {
	inStream := bytes.NewBufferString(input)
	outStream := new(bytes.Buffer)
	errStream := new(bytes.Buffer)

	cmd.SetIn(inStream)
	cmd.SetOut(outStream)
	cmd.SetErr(errStream)

	return cmd
}

// assertMockCmd is a high-level test helper to run a given mock command with
// arguments and check if an error and its stdout are expected.
func assertMockCmd(t *testing.T, cmd *cobra.Command, args []string, ok bool, want string) {
	err := runMockCmd(cmd, args)

	stderr := mockErr(cmd)
	if ok && err != nil {
		t.Fatalf("unexpected err = %s, stderr: \n%s", err, stderr)
	}

	stdout := mockOut(cmd)
	if !ok && err == nil {
		t.Fatalf("expected to return an error, but no error, stdout: \n%s", stdout)
	}

	if stdout != want {
		t.Fatalf("got:\n%s\nwant:\n%s", stdout, want)
	}
}

// runMockCmd is a helper function which parses flags and invokes a given mock
// command.
func runMockCmd(cmd *cobra.Command, args []string) error {
	cmdFlags := cmd.Flags()
	if err := cmdFlags.Parse(args); err != nil {
		return fmt.Errorf("failed to parse arguments: %s", err)
	}

	return cmd.RunE(cmd, cmdFlags.Args())
}

// mockErr is a helper function which returns a string written to mocked err stream.
func mockErr(cmd *cobra.Command) string {
	return cmd.ErrOrStderr().(*bytes.Buffer).String()
}

// mockOut is a helper function which returns a string written to mocked out stream.
func mockOut(cmd *cobra.Command) string {
	return cmd.OutOrStdout().(*bytes.Buffer).String()
}
