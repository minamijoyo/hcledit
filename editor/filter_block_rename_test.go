package editor

import (
	"testing"
)

func TestBlockRenameFilter(t *testing.T) {
	cases := []struct {
		name string
		src  string
		from string
		to   string
		ok   bool
		want string
	}{
		{
			name: "simple",
			src: `a0 = v0
b1 "l1" {
  a2 = v2
}

b2 "l2" {
}
`,
			from: "b1.l1",
			to:   "b1.l2",
			ok:   true,
			want: `a0 = v0
b1 "l2" {
  a2 = v2
}

b2 "l2" {
}
`,
		},

		{
			name: "escaped address",
			src: `a0 = v0
b1 "l.1" {
  a2 = v2
}

b2 "l2" {
}
`,
			from: `b1.l\.1`,
			to:   `b1.l\.2`,
			ok:   true,
			want: `a0 = v0
b1 "l.2" {
  a2 = v2
}

b2 "l2" {
}
`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			o := NewEditOperator(NewBlockRenameFilter(tc.from, tc.to))
			output, err := o.Apply([]byte(tc.src), "test")
			if tc.ok && err != nil {
				t.Fatalf("unexpected err = %s", err)
			}

			got := string(output)
			if !tc.ok && err == nil {
				t.Fatalf("expected to return an error, but no error, outStream: \n%s", got)
			}

			if got != tc.want {
				t.Fatalf("got:\n%s\nwant:\n%s", got, tc.want)
			}
		})
	}
}
