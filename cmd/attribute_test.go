package cmd

import (
	"testing"
)

func TestAttributeGet(t *testing.T) {
	src := `terraform {
  backend "s3" {
    region = "ap-northeast-1"
    bucket = "minamijoyo-hcledit"
    key    = "services/hoge/dev/terraform.tfstate"
  }
}
locals {
  map = {
    # comment
    attribute = "bar"
  }
  attribute = "foo" # comment
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
			args: []string{"terraform.backend.s3.key"},
			ok:   true,
			want: "\"services/hoge/dev/terraform.tfstate\"\n",
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
		{
			name: "with comments",
			args: []string{"--with-comments", "locals.map"},
			ok:   true,
			want: `{
    # comment
    attribute = "bar"
  }
`,
		},
		{
			name: "without comments",
			args: []string{"locals.map"},
			ok:   true,
			want: `{

    attribute = "bar"
  }
`,
		},
		{
			name: "single with comments",
			args: []string{"--with-comments", "locals.attribute"},
			ok:   true,
			want: `"foo" # comment
`,
		},
		{
			name: "single without comments",
			args: []string{"locals.attribute"},
			ok:   true,
			want: `"foo"
`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cmd := newMockCmd(newAttributeGetCmd(), src)
			assertMockCmd(t, cmd, tc.args, tc.ok, tc.want)
		})
	}
}

func TestAttributeSet(t *testing.T) {
	src := `terraform {
  backend "s3" {
    region = "ap-northeast-1"
    bucket = "minamijoyo-hcledit"
    key    = "services/hoge/dev/terraform.tfstate"
  }
}
module "hoge" {
  source = "./hoge"
  env    = "dev"
}
`

	cases := []struct {
		name string
		args []string
		ok   bool
		want string
	}{
		{
			name: "string literal",
			args: []string{"terraform.backend.s3.key", `"services/fuga/dev/terraform.tfstate"`},
			ok:   true,
			want: `terraform {
  backend "s3" {
    region = "ap-northeast-1"
    bucket = "minamijoyo-hcledit"
    key    = "services/fuga/dev/terraform.tfstate"
  }
}
module "hoge" {
  source = "./hoge"
  env    = "dev"
}
`,
		},
		{
			name: "string literal to variable reference",
			args: []string{"module.hoge.env", "var.env"},
			ok:   true,
			want: `terraform {
  backend "s3" {
    region = "ap-northeast-1"
    bucket = "minamijoyo-hcledit"
    key    = "services/hoge/dev/terraform.tfstate"
  }
}
module "hoge" {
  source = "./hoge"
  env    = var.env
}
`,
		},
		{
			name: "no match",
			args: []string{"hoge", "fuga"},
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
			cmd := newMockCmd(newAttributeSetCmd(), src)
			assertMockCmd(t, cmd, tc.args, tc.ok, tc.want)
		})
	}
}

func TestAttributeMv(t *testing.T) {
	src := `locals {
  foo1 = "bar1"
  foo2 = "bar2"
}

resource "foo" "bar" {
  foo3 = "bar3"
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
			args: []string{"locals.foo1", "locals.foo3"},
			ok:   true,
			want: `locals {
  foo3 = "bar1"
  foo2 = "bar2"
}

resource "foo" "bar" {
  foo3 = "bar3"
}
`,
		},
		{
			name: "no match",
			args: []string{"locals.foo3", "locals.foo4"},
			ok:   true,
			want: src,
		},
		{
			name: "duplicated",
			args: []string{"locals.foo1", "locals.foo2"},
			ok:   false,
			want: "",
		},
		{
			name: "move an attribute accross blocks",
			args: []string{"locals.foo1", "resource.foo.bar.foo1"},
			ok:   false,
			want: "",
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
			cmd := newMockCmd(newAttributeMvCmd(), src)
			assertMockCmd(t, cmd, tc.args, tc.ok, tc.want)
		})
	}
}

func TestAttributeRm(t *testing.T) {
	src := `locals {
  service = "hoge"
  env     = "dev"
  region  = "ap-northeast-1"
}`

	cases := []struct {
		name string
		args []string
		ok   bool
		want string
	}{
		{
			name: "remove an unused local variable",
			args: []string{"locals.region"},
			ok:   true,
			want: `locals {
  service = "hoge"
  env     = "dev"
}`,
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
			cmd := newMockCmd(newAttributeRmCmd(), src)
			assertMockCmd(t, cmd, tc.args, tc.ok, tc.want)
		})
	}
}

func TestAttributeAppend(t *testing.T) {
	src := `terraform {
  required_version = "0.13.5"

  backend "s3" {
    region = "ap-northeast-1"
    bucket = "foo"
    key    = "bar/terraform.tfstate"
  }

  required_providers {
  }
}
`

	cases := []struct {
		name string
		args []string
		ok   bool
		want string
	}{
		{
			name: "map literal",
			args: []string{"terraform.required_providers.aws", `{
  source = "hashicorp/aws"
  version = "3.11.0"
}`},
			ok: true,
			want: `terraform {
  required_version = "0.13.5"

  backend "s3" {
    region = "ap-northeast-1"
    bucket = "foo"
    key    = "bar/terraform.tfstate"
  }

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "3.11.0"
    }
  }
}
`,
		},
		{
			name: "no match",
			args: []string{"foo.bar", "baz"},
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
			args: []string{"terraform.required_providers.aws"},
			ok:   false,
			want: "",
		},
		{
			name: "too many args",
			args: []string{"terraform.required_providers.aws", "foo", "var"},
			ok:   false,
			want: "",
		},
		{
			name: "map literal (newline)",
			args: []string{
				"terraform.required_providers.aws",
				`{
  source = "hashicorp/aws"
  version = "3.11.0"
}`,
				"--newline"},
			ok: true,
			want: `terraform {
  required_version = "0.13.5"

  backend "s3" {
    region = "ap-northeast-1"
    bucket = "foo"
    key    = "bar/terraform.tfstate"
  }

  required_providers {

    aws = {
      source  = "hashicorp/aws"
      version = "3.11.0"
    }
  }
}
`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cmd := newMockCmd(newAttributeAppendCmd(), src)
			assertMockCmd(t, cmd, tc.args, tc.ok, tc.want)
		})
	}
}
