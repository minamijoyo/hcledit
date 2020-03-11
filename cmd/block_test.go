package cmd

import (
	"testing"
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
			cmd := newMockCmd(runBlockGetCmd, src)

			err := runBlockGetCmd(cmd, tc.args)
			stderr := mockErr(cmd)
			if tc.ok && err != nil {
				t.Fatalf("unexpected err = %s, stderr: \n%s", err, stderr)
			}

			stdout := mockOut(cmd)
			if !tc.ok && err == nil {
				t.Fatalf("expected to return an error, but no error, stdout: \n%s", stdout)
			}

			if stdout != tc.want {
				t.Fatalf("got:\n%s\nwant:\n%s", stdout, tc.want)
			}
		})
	}
}

func TestBlockMv(t *testing.T) {
	src := `resource "aws_security_group" "test1" {
  name = "tfedit-test1"
}

resource "aws_security_group" "test2" {
  name = "tfedit-test2"
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
			args: []string{"resource.aws_security_group.test1", "resource.aws_security_group.test3"},
			ok:   true,
			want: `resource "aws_security_group" "test3" {
  name = "tfedit-test1"
}

resource "aws_security_group" "test2" {
  name = "tfedit-test2"
}
`,
		},
		{
			name: "no match",
			args: []string{"resource.aws_security_group.test", "resource.aws_security_group.test3"},
			ok:   true,
			want: src,
		},
		{
			name: "no args",
			args: []string{},
			ok:   false,
			want: "",
		},
		{
			name: "1 arg",
			args: []string{"hoge"},
			ok:   false,
			want: "",
		},
		{
			name: "too many args",
			args: []string{"hoge", "fuga", "piyo"},
			ok:   false,
			want: "",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cmd := newMockCmd(runBlockGetCmd, src)

			err := runBlockMvCmd(cmd, tc.args)
			stderr := mockErr(cmd)
			if tc.ok && err != nil {
				t.Fatalf("unexpected err = %s, stderr: \n%s", err, stderr)
			}

			stdout := mockOut(cmd)
			if !tc.ok && err == nil {
				t.Fatalf("expected to return an error, but no error, stdout: \n%s", stdout)
			}

			if stdout != tc.want {
				t.Fatalf("got:\n%s\nwant:\n%s", stdout, tc.want)
			}
		})
	}
}

func TestBlockList(t *testing.T) {
	src := `terraform {
  required_version = "0.12.18"
}

provider "aws" {
  version = "2.43.0"
  region  = "ap-northeast-1"
}

resource "aws_security_group" "hoge" {
  name = "hoge"
  egress {
    from_port = 0
    to_port   = 0
    protocol  = -1
  }
}

resource "aws_security_group" "fuga" {
  name = "fuga"
  egress {
    from_port = 0
    to_port   = 0
    protocol  = -1
  }
}
`

	cases := []struct {
		name string
		ok   bool
		want string
	}{
		{
			name: "simple",
			ok:   true,
			want: `terraform
provider.aws
resource.aws_security_group.hoge
resource.aws_security_group.fuga
`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cmd := newMockCmd(runBlockGetCmd, src)

			args := []string{}
			err := runBlockListCmd(cmd, args)
			stderr := mockErr(cmd)
			if tc.ok && err != nil {
				t.Fatalf("unexpected err = %s, stderr: \n%s", err, stderr)
			}

			stdout := mockOut(cmd)
			if !tc.ok && err == nil {
				t.Fatalf("expected to return an error, but no error, stdout: \n%s", stdout)
			}

			if stdout != tc.want {
				t.Fatalf("got:\n%s\nwant:\n%s", stdout, tc.want)
			}
		})
	}
}

func TestBlockRm(t *testing.T) {
	src := `data "aws_security_group" "hoge" {
  name = "hoge"
}

data "aws_security_group" "fuga" {
  name = "fuga"
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
			args: []string{"data.aws_security_group.hoge"},
			ok:   true,
			want: `data "aws_security_group" "fuga" {
  name = "fuga"
}
`,
		},
		{
			name: "no match",
			args: []string{"hoge"},
			ok:   true,
			want: src,
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
			cmd := newMockCmd(runBlockGetCmd, src)

			err := runBlockRmCmd(cmd, tc.args)
			stderr := mockErr(cmd)
			if tc.ok && err != nil {
				t.Fatalf("unexpected err = %s, stderr: \n%s", err, stderr)
			}

			stdout := mockOut(cmd)
			if !tc.ok && err == nil {
				t.Fatalf("expected to return an error, but no error, stdout: \n%s", stdout)
			}

			if stdout != tc.want {
				t.Fatalf("got:\n%s\nwant:\n%s", stdout, tc.want)
			}
		})
	}
}
