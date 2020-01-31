package hclwritex

import (
	"bytes"
	"testing"
)

func TestBlockRename(t *testing.T) {
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
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			inStream := bytes.NewBufferString(tc.src)
			outStream := new(bytes.Buffer)
			err := RenameBlock(inStream, outStream, "test", tc.from, tc.to)
			if tc.ok && err != nil {
				t.Fatalf("unexpected err = %s", err)
			}

			got := outStream.String()
			if !tc.ok && err == nil {
				t.Fatalf("expected to return an error, but no error, outStream: \n%s", got)
			}

			if got != tc.want {
				t.Fatalf("got:\n%s\nwant:\n%s", got, tc.want)
			}
		})
	}
}
