package editor

import (
	"testing"
)

func TestAttributeReplaceFilter(t *testing.T) {
	cases := []struct {
		name    string
		src     string
		address string
		toName  string
		toValue string
		ok      bool
		want    string
	}{
		{
			name: "simple",
			src: `
a0 = v0
a1 = v1
`,
			address: "a0",
			toName:  "a2",
			toValue: "v2",
			ok:      true,
			want: `
a2 = v2
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
			address: "a0",
			toName:  "a2",
			toValue: `"v2"`,
			ok:      true,
			want: `
# before attr
a2 = "v2" # inline
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
			address: "b1.l1.a1",
			toName:  "a2",
			toValue: "v2",
			ok:      true,
			want: `
a0 = v0
b1 "l1" {
  a2 = v2
}
`,
		},
		{
			name: "not found",
			src: `
a0 = v0
`,
			address: "a1",
			toName:  "a2",
			toValue: "v2",
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
			toName:  "a3",
			toValue: "v3",
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
			toName:  "a2",
			toValue: "v2",
			ok:      true,
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
			address: "a0",
			toName:  "a1",
			toValue: "v2",
			ok:      false,
			want:    "",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			o := NewEditOperator(NewAttributeReplaceFilter(tc.address, tc.toName, tc.toValue))
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
