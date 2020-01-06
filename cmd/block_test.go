package cmd

import (
	"bytes"
	"testing"

	"github.com/spf13/cobra"
)

func setMockStream(cmd *cobra.Command) {
	inStream := new(bytes.Buffer)
	outStream := new(bytes.Buffer)
	errStream := new(bytes.Buffer)
	cmd.SetIn(inStream)
	cmd.SetOut(outStream)
	cmd.SetErr(errStream)
}

func TestBlockGet(t *testing.T) {
	src := `terraform {
  required_version = "0.12.18"
}

provider "aws" {
  version = "2.43.0"
  region  = "ap-northeast-1"
}
`

	cases := []struct {
		name string
		args []string
		ok   bool
		want string
	}{
		{
			name: "simple",
			args: []string{"terraform"},
			ok:   true,
			want: `terraform {
  required_version = "0.12.18"
}
`,
		},
		{
			name: "no match",
			args: []string{"hoge"},
			ok:   true,
			want: "",
		},
		{
			name: "no args",
			args: []string{},
			ok:   false,
			want: "",
		},
		{
			name: "too many args",
			args: []string{"hoge", "fuga"},
			ok:   false,
			want: "",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cmd := &cobra.Command{
				RunE: runBlockGetCmd,
			}

			setMockStream(cmd)
			inStream := cmd.InOrStdin().(*bytes.Buffer)
			inStream.Write([]byte(src))

			err := runBlockGetCmd(cmd, tc.args)
			stderr := cmd.ErrOrStderr().(*bytes.Buffer).String()
			if tc.ok && err != nil {
				t.Fatalf("unexpected err = %s, stderr: \n%s", err, stderr)
			}

			stdout := cmd.OutOrStdout().(*bytes.Buffer).String()
			if !tc.ok && err == nil {
				t.Fatalf("expected to return an error, but no error, stdout: \n%s", stdout)
			}

			if stdout != tc.want {
				t.Fatalf("got:\n%s\nwant:\n%s", stdout, tc.want)
			}
		})
	}
}
