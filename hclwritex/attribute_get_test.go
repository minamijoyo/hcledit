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
		{
			name: "attribute in block",
			src: `
b1 {
  a1 = v1
}
`,
			address: "b1.a1",
			ok:      true,
			want:    "v1",
		},
		{
			name: "attribute in block with a label",
			src: `
b1 "l1" {
  a1 = v1
}
`,
			address: "b1.l1.a1",
			ok:      true,
			want:    "v1",
		},
		{
			name: "attribute in block with multiple labels",
			src: `
b1 {
  a1 = v0
}
b1 "l1" {
  a1 = v1
}
b1 "l1" "l2" {
  a1 = v2
}
b1 "l1" "l2" "l3" {
  a1 = v3
}
`,
			address: "b1.l1.l2.a1",
			ok:      true,
			want:    "v2",
		},
		{
			name: "attribute in nested block",
			src: `
b1 {
  a1 = v1
  b2 {
    a2 = v2
  }
}
`,
			address: "b1.b2.a2",
			ok:      true,
			want:    "v2",
		},
		{
			name: "attribute in nested block (extra labels)",
			src: `
b1 "l1" {
  a1 = v1
  b2 {
    a2 = v2
  }
}
`,
			address: "b1.b2.a2",
			ok:      true,
			want:    "",
		},
		{
			name: "labels take precedence over nested blocks",
			src: `
b1 "b2" {
  a1 = v1
  b2 {
    a1 = v2
  }
}
`,
			address: "b1.b2.a1",
			ok:      true,
			want:    "v1",
		},
		{
			name: "attribute in multi level nested block",
			src: `
b1 {
  a1 = v1
  b2 {
    a2 = v2
    b3 {
      a3 = v3
    }
  }
}
`,
			address: "b1.b2.b3.a3",
			ok:      true,
			want:    "v3",
		},
		{
			name: "attribute in nested block with labels",
			src: `
b1 {
  a1 = v1
  b2 "b3" {
    a2 = v2
    b3 {
      a2 = v3
    }
  }
}
`,
			address: "b1.b2.b3.a2",
			ok:      true,
			want:    "v2",
		},
		{
			name: "attribute in duplicated blocks",
			src: `
b1 "l1" "l2" {
  a1 = v1
}
b1 "l1" "l2" {
  a1 = v2
}
`,
			address: "b1.l1.l2.a1",
			ok:      true,
			want:    "v1",
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
