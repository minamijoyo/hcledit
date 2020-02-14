package cmd

import (
	"regexp"
	"testing"
)

func TestVersion(t *testing.T) {
	cases := []struct {
		name   string
		args   []string
		ok     bool
		wantRe *regexp.Regexp
	}{
		{
			name:   "simple",
			args:   []string{},
			ok:     true,
			wantRe: regexp.MustCompile(`[0-9]+(\.[0-9]+)*(-.*)*`),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cmd := newMockCmd(runVersionCmd, "")

			err := runVersionCmd(cmd, tc.args)
			stderr := mockErr(cmd)
			if tc.ok && err != nil {
				t.Fatalf("unexpected err = %s, stderr: \n%s", err, stderr)
			}

			stdout := mockOut(cmd)
			if !tc.ok && err == nil {
				t.Fatalf("expected to return an error, but no error, stdout: \n%s", stdout)
			}

			if !tc.wantRe.MatchString(stdout) {
				t.Fatalf("got:\n%s\nwantRe:\n%s", stdout, tc.wantRe)
			}
		})
	}
}
