package editor

import (
	"testing"
)

func TestAttributeRenameFilter(t *testing.T) {
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
			src: `
a0 = v0
a1 = v1
`,
			from: "a0",
			to:   "a2",
			ok:   true,
			want: `
a2 = v0
a1 = v1
`,
		},
		{
			name: "with comments",
			src: `
# before attr
a0 = "v0" # inline
a1 = "v1"
`,
			from: "a0",
			to:   "a2",
			ok:   true,
			want: `
# before attr
a2 = "v0" # inline
a1 = "v1"
`,
		},
		{
			name: "attribute in block",
			src: `
a0 = v0
b1 "l1" {
  a1 = v1
}
`,
			from: "b1.l1.a1",
			to:   "b1.l1.a2",
			ok:   true,
			want: `
a0 = v0
b1 "l1" {
  a2 = v1
}
`,
		},
		{
			name: "not found",
			src: `
a0 = v0
`,
			from: "a1",
			to:   "a2",
			ok:   true,
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
			from: "b1.l1.a2",
			to:   "b1.l1.a3",
			ok:   true,
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
			from: "b2.l1.a1",
			to:   "b2.l1.a2",
			ok:   true,
			want: `
a0 = v0
b1 "l1" {
  a1 = v1
}
`,
		},
		{
			name: "attribute already exists",
			src: `
a0 = v0
a1 = v1
`,
			from: "a0",
			to:   "a1",
			ok:   false,
			want: "",
		},
		{
			name: "moving an attribute across blocks",
			src: `
a0 = v0
b1 "l1" {
  a1 = v1
}
`,
			from: "a0",
			to:   "b1.l1.a0",
			ok:   false,
			want: "",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			o := NewEditOperator(NewAttributeRenameFilter(tc.from, tc.to))
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
