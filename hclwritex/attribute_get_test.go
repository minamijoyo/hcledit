package hclwritex

import (
	"bytes"
	"testing"
)

func TestAttributeGet(t *testing.T) {
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
			want:    "v0",
		},
		{
			name: "quoted literal should be unquoted",
			src: `
a0 = "v0"
`,
			address: "a0",
			ok:      true,
			want:    "v0",
		},
		{
			name: "not found",
			src: `
a0 = v0
a1 = v1
`,
			address: "hoge",
			ok:      true,
			want:    "",
		},
		{
			name: "attribute with comments",
			src: `
// attr comment
a0 = v0 // inline comment
a1 = v1
`,
			address: "a0",
			ok:      true,
			want:    "v0",
		},
		{
			name: "duplicated attributes should be error",
			src: `
a0 = v0
a0 = v1
`,
			address: "a0",
			ok:      false,
			want:    "",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			inStream := bytes.NewBufferString(tc.src)
			outStream := new(bytes.Buffer)
			err := GetAttribute(inStream, outStream, "test", tc.address)
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
