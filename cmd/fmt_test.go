package cmd

import (
	"testing"
)

func TestFmt(t *testing.T) {

	cases := []struct {
		name string
		src  string
		args []string
		ok   bool
		want string
	}{
		{
			name: "unformatted",
			src: `
resource "foo" "bar" {
  attr1 = "val1"
  attr2="val2"
}
`,
			args: []string{},
			ok:   true,
			want: `
resource "foo" "bar" {
  attr1 = "val1"
  attr2 = "val2"
}
`,
		},
		{
			name: "syntax error",
			src: `
resource "foo" "bar" {
  attr1 = "val1"
`,
			args: []string{},
			ok:   false,
			want: "",
		},
		{
			name: "too many args",
			src:  "",
			args: []string{"foo"},
			ok:   false,
			want: "",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cmd := newMockCmd(newFmtCmd(), tc.src)
			assertMockCmd(t, cmd, tc.args, tc.ok, tc.want)
		})
	}
}
