package editor

import (
	"testing"
)

func TestBlockNewFilter(t *testing.T) {
	cases := []struct {
		name    string
		src     string
		address string
		newline bool
		want    string
	}{
		{
			name: "block with blockType and 2 labels, resource with newline",
			src: `
variable "var1" {
  type        = string
  default     = "foo"
  description = "example variable"
}
`,
			address: "resource.aws_instance.example",
			newline: true,
			want: `
variable "var1" {
  type        = string
  default     = "foo"
  description = "example variable"
}

resource "aws_instance" "example" {
}
`,
		},
		{
			name: "block with blockType and 1 label, module without newline",
			src: `
variable "var1" {
  type        = string
  default     = "foo"
  description = "example variable"
}
`,
			address: "module.example",
			newline: false,
			want: `
variable "var1" {
  type        = string
  default     = "foo"
  description = "example variable"
}
module "example" {
}
`,
		},
		{
			name: "block with blockType and 0 labels, locals without newline",
			src: `
variable "var1" {
  type        = string
  default     = "foo"
  description = "example variable"
}
`,
			address: "locals",
			newline: false,
			want: `
variable "var1" {
  type        = string
  default     = "foo"
  description = "example variable"
}
locals {
}
`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			o := NewEditOperator(NewBlockNewFilter(tc.address, tc.newline))
			output, err := o.Apply([]byte(tc.src), "test")
			if err != nil {
				t.Fatalf("unexpected err = %s", err)
			}

			got := string(output)
			if got != tc.want {
				t.Fatalf("got:\n%s\nwant:\n%s", got, tc.want)
			}
		})
	}
}
