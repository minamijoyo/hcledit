package command

import (
	"bytes"
	"testing"

	"github.com/mitchellh/cli"
	"github.com/spf13/afero"
)

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
		name     string
		args     []string
		exitCode int
		want     string
	}{
		{
			name:     "simple",
			args:     []string{"terraform"},
			exitCode: 0,
			want: `terraform {
  required_version = "0.12.18"
}

`,
		},
		{
			name:     "no match",
			args:     []string{"hoge"},
			exitCode: 0,
			want:     "",
		},
		{
			name:     "no args",
			args:     []string{},
			exitCode: 1,
			want:     "",
		},
		{
			name:     "too many args",
			args:     []string{"hoge", "fuga"},
			exitCode: 1,
			want:     "",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			input := new(bytes.Buffer)
			input.Write([]byte(src))

			ui := cli.NewMockUi()
			meta := Meta{
				UI:    ui,
				Input: input,
				Fs:    afero.NewMemMapFs(),
			}

			c := &BlockGetCommand{
				Meta: meta,
			}

			if exitCode := c.Run(tc.args); exitCode != tc.exitCode {
				t.Fatalf("unexpected exit code: got = %d, but want = %d, stderr: \n%s", exitCode, tc.exitCode, ui.ErrorWriter.String())
			}

			if got := ui.OutputWriter.String(); got != tc.want {
				t.Fatalf("got: %s\nwant: %s", got, tc.want)
			}
		})
	}
}
