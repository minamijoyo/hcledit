package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"
	"testing"
)

const (
	VarRunMainForTesting          = "RUN_MAIN_FOR_TESTING"
	VarRunMainForTestingArgPrefix = "RUN_MAIN_FOR_TESTING_ARG_"
)

func TestHCLEditMain(t *testing.T) {
	if os.Getenv(VarRunMainForTesting) == "1" {
		os.Args = append([]string{"hcledit"}, envToArgs(os.Environ())...)

		// We DO call hcledit's main() here. So this looks like a normal `hcledit` process.
		main()

		// If `main()` did not call os.Exit(0) explicitly, we assume there was no error hence it's safe to call os.Exit(0)
		// on behalf of go runtime.
		os.Exit(0)

		// As main() or this block calls os.Exit, we never reach this line.
		// But the test called this block of code catches and verifies the exit code.
		return
	}

	testcases := []struct {
		subject      string
		input        string
		args         []string
		wantStdout   string
		wantStderr   string
		wantExitCode int
	}{
		{
			subject: "set existing nested attribute",
			input: `resource "foo" "bar" {
  attr1 = "val1"
  nested {
    attr2 = "val2"
  }
}
`,
			args: []string{
				"attribute",
				"set",
				"resource.foo.bar.nested.attr2",
				"\"val3\"",
			},
			wantStdout: `resource "foo" "bar" {
  attr1 = "val1"
  nested {
    attr2 = "val3"
  }
}
`,
		},
		{
			subject: "set with insufficient args",
			input: `resource "foo" "bar" {
  attr1 = "val1"
  nested {
    attr2 = "val2"
  }
}
`,
			args: []string{
				"attribute",
				"set",
				"resource.foo.bar.nested.attr2",
			},
			wantStderr: `expected 2 argument, but got 1 arguments
`,
			wantExitCode: 1,
		},
	}

	for i := range testcases {
		tc := testcases[i]

		t.Run(tc.subject, func(t *testing.T) {
			// Do a second run of this specific test(TestHCLEditMain) with RUN_MAIN_FOR_TESTING=1 set,
			// So that the second run is able to run main() and this first run can verify the exit status returned by that.
			//
			// This technique originates from https://talks.golang.org/2014/testing.slide#23.
			cmd := exec.Command(os.Args[0], "-test.run=TestHCLEditMain") // nolint: gosec
			cmd.Env = append(
				cmd.Env,
				os.Environ()...,
			)
			cmd.Env = append(
				cmd.Env,
				VarRunMainForTesting+"=1",
			)
			cmd.Env = append(
				cmd.Env,
				argsToEnv(tc.args)...,
			)

			stdin := strings.NewReader(tc.input)
			stdout := &bytes.Buffer{}
			stderr := &bytes.Buffer{}

			cmd.Stdin = stdin
			cmd.Stdout = stdout
			cmd.Stderr = stderr

			err := cmd.Run()

			if got := stdout.String(); got != tc.wantStdout {
				t.Errorf("Unexpected stdout: want %q, got %q", tc.wantStdout, got)
			}

			if got := stderr.String(); got != tc.wantStderr {
				t.Errorf("Unexpected stderr: want %q, got %q", tc.wantStderr, got)
			}

			if tc.wantExitCode != 0 {
				exiterr, ok := err.(*exec.ExitError)

				if !ok {
					t.Fatalf("Unexpected error returned by os.Exit: %T", err)
				}

				if got := exiterr.ExitCode(); got != tc.wantExitCode {
					t.Errorf("Unexpected exit code: want %d, got %d", tc.wantExitCode, got)
				}
			}
		})
	}
}

func argsToEnv(args []string) []string {
	var env []string

	for i, arg := range args {
		env = append(env, fmt.Sprintf("%s%d=%s", VarRunMainForTestingArgPrefix, i, arg))
	}

	return env
}

func envToArgs(env []string) []string {
	var envvars []string

	for _, kv := range os.Environ() {
		if strings.HasPrefix(kv, VarRunMainForTestingArgPrefix) {
			envvars = append(envvars, kv)
		}
	}

	sort.Strings(envvars)

	var args []string

	for _, kv := range envvars {
		args = append(args, strings.Split(kv, "=")[1])
	}

	return args
}
