package editor

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestAttributeGetSink(t *testing.T) {
	cases := []struct {
		name         string
		src          string
		address      string
		withComments bool
		ok           bool
		want         string
	}{
		{
			name: "simple top level attribute",
			src: `
a0 = v0
a1 = v1
`,
			address:      "a0",
			withComments: false,
			ok:           true,
			want:         "v0\n",
		},
		{
			name: "quoted literal is as it is and should not be unquoted",
			src: `
a0 = "v0"
`,
			address:      "a0",
			withComments: false,
			ok:           true,
			want:         "\"v0\"\n",
		},
		{
			name: "not found",
			src: `
a0 = v0
a1 = v1
`,
			address:      "hoge",
			withComments: false,
			ok:           true,
			want:         "",
		},
		{
			name: "attribute with comments",
			src: `
// attr comment
a0 = v0 // inline comment
a1 = v1
`,
			address:      "a0",
			withComments: false,
			ok:           true,
			want:         "v0\n",
		},
		{
			name: "multiline attribute with comments",
			src: `
// attr comment
a0 = v0
a1 = [
  "val1",
  "val2", // inline comment
  "val3", # another comment
  "val4",
  # a ocmment line
  "val5",
]
a2 = v2
`,
			address:      "a1",
			withComments: false,
			ok:           true,
			want: `[
  "val1",
  "val2",
  "val3",
  "val4",

  "val5",
]
`,
		},
		{
			name: "duplicated attributes should be error",
			src: `
a0 = v0
a0 = v1
`,
			address:      "a0",
			withComments: false,
			ok:           false,
			want:         "",
		},
		{
			name: "attribute in block",
			src: `
b1 {
  a1 = v1
}
`,
			address:      "b1.a1",
			withComments: false,
			ok:           true,
			want:         "v1\n",
		},
		{
			name: "attribute in block with a label",
			src: `
b1 "l1" {
  a1 = v1
}
`,
			address:      "b1.l1.a1",
			withComments: false,
			ok:           true,
			want:         "v1\n",
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
			address:      "b1.l1.l2.a1",
			withComments: false,
			ok:           true,
			want:         "v2\n",
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
			address:      "b1.b2.a2",
			withComments: false,
			ok:           true,
			want:         "v2\n",
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
			address:      "b1.b2.a2",
			withComments: false,
			ok:           true,
			want:         "",
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
			address:      "b1.b2.a1",
			withComments: false,
			ok:           true,
			want:         "v1\n",
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
			address:      "b1.b2.b3.a3",
			withComments: false,
			ok:           true,
			want:         "v3\n",
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
			address:      "b1.b2.b3.a2",
			withComments: false,
			ok:           true,
			want:         "v2\n",
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
			address:      "b1.l1.l2.a1",
			withComments: false,
			ok:           true,
			want:         "v1\n",
		},
		{
			name: "attribute in block with a escaped address",
			src: `
b1 "l.1" {
  a1 = v1
}
`,
			address:      `b1.l\.1.a1`,
			withComments: false,
			ok:           true,
			want:         "v1\n",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			o := NewDeriveOperator(NewAttributeGetSink(tc.address, tc.withComments))
			output, err := o.Apply([]byte(tc.src), "test")
			if tc.ok && err != nil {
				t.Fatalf("unexpected err = %s", err)
			}

			got := string(output)
			if !tc.ok && err == nil {
				t.Fatalf("expected to return an error, but no error, outStream: \n%s", got)
			}

			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatalf("got:\n%s\nwant:\n%s\ndiff(-want +got):\n%v", got, tc.want, diff)
			}
		})
	}
}
