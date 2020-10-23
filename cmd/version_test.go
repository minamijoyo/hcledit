package cmd

import (
	"testing"
)

func TestVersion(t *testing.T) {
	cases := []struct {
		name string
		args []string
		ok   bool
		want string
	}{
		{
			name: "simple",
			args: []string{},
			ok:   true,
			want: Version + "\n",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cmd := newMockCmd(newVersionCmd(), "")
			assertMockCmd(t, cmd, tc.args, tc.ok, tc.want)
		})
	}
}
