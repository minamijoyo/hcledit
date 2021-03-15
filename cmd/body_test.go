package cmd

import (
	"testing"
)

func TestBodyGet(t *testing.T) {
	src := `resource "foo" "bar" {
  attr1 = "val1"
  nested {
    attr2 = "val2"
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
			name: "simple",
			args: []string{"resource.foo.bar"},
			ok:   true,
			want: `attr1 = "val1"
nested {
  attr2 = "val2"
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
			cmd := newMockCmd(newBodyGetCmd(), src)
			assertMockCmd(t, cmd, tc.args, tc.ok, tc.want)
		})
	}
}
