package editor

import (
	"testing"
)

func TestAttributeRemoveFilter(t *testing.T) {
	cases := []struct {
		name    string
		src     string
		address string
		ok      bool
		want    string
	}{
		{
			name: "simple top level attribute",
			src: `
a0 = v0
a1 = v1
`,
			address: "a0",
			ok:      true,
			want: `
a1 = v1
`,
		},
		{
			name: "simple top level attribute (with comments)",
			src: `
// before attr
a0 = "v0" // inline
a1 = "v1"
`,
			address: "a0",
			ok:      true,
			want: `
a1 = "v1"
`, // Unfortunately we can't keep the before attr comment.
		},
		{
			name: "attribute in block",
			src: `
a0 = v0
b1 "l1" {
  a1 = v1
  a2 = v2
}
`,
			address: "b1.l1.a1",
			ok:      true,
			want: `
a0 = v0
b1 "l1" {
  a2 = v2
}
`,
		},
		{
			name: "top level attribute not found",
			src: `
a0 = v0
`,
			address: "a1",
			ok:      true,
			want: `
a0 = v0
`,
		},
		{
			name: "attribute not found in block",
			src: `
a0 = v0
b1 "l1" {
  a1 = v1
}
`,
			address: "b1.l1.a2",
			ok:      true,
			want: `
a0 = v0
b1 "l1" {
  a1 = v1
}
`,
		},
		{
			name: "block not found",
			src: `
a0 = v0
b1 "l1" {
  a1 = v1
}
`,
			address: "b2.l1.a1",
			ok:      true,
			want: `
a0 = v0
b1 "l1" {
  a1 = v1
}
`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			o := NewEditOperator(NewAttributeRemoveFilter(tc.address))
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
