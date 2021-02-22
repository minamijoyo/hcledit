package editor

import (
	"bytes"
	"testing"
)

func TestBlockGet(t *testing.T) {
	cases := []struct {
		name    string
		src     string
		address string
		ok      bool
		want    string
	}{
		{
			name: "simple",
			src: `
a0 = v0
b1 {
  a2 = v2
}

b2 l1 {
}
`,
			address: "b1",
			ok:      true,
			want: `b1 {
  a2 = v2
}
`,
		},
		{
			name:    "no match",
			address: "hoge",
			ok:      true,
			want:    "",
		},
		{
			name:    "empty",
			address: "",
			ok:      false,
			want:    "",
		},
		{
			name: "unformatted",
			src: `
  b1   {
}
`,
			address: "b1",
			ok:      true,
			want: `b1 {
}
`,
		},
		{
			name: "no label",
			src: `
b1 {
}

b1 l1 {
}
`,
			address: "b1",
			ok:      true,
			want: `b1 {
}
`,
		},
		{
			name: "with label",
			src: `
b1 {
}

b1 l1 {
}
`,
			address: "b1.l1",
			ok:      true,
			want: `b1 l1 {
}
`,
		},
		{
			name: "multi blocks",
			src: `
b1 {
}

b1 l1 {
}

b1 l1 {
}
`,
			address: "b1.l1",
			ok:      true,
			want: `b1 l1 {
}

b1 l1 {
}
`,
		},
		{
			name: "get a given block type and any labels",
			src: `
b1 {
}

b1 l1 {
}

b1 l2 {
}
`,
			address: "b1.*",
			ok:      true,
			want: `b1 l1 {
}

b1 l2 {
}
`,
		},
		{
			name: "get a given block type and prefixed labels",
			src: `
b1 {
}

b1 l1 {
}

b1 l1 l2 {
}

b1 l1 l3 {
}
`,
			address: "b1.l1.*",
			ok:      true,
			want: `b1 l1 l2 {
}

b1 l1 l3 {
}
`,
		},
		{
			name: "preserve comments",
			src: `// before block
b1 {
  // before attr
  attr = val // inline
}
// after block
`,
			address: "b1",
			ok:      true,
			want: `// before block
b1 {
  // before attr
  attr = val // inline
}
`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			inStream := bytes.NewBufferString(tc.src)
			outStream := new(bytes.Buffer)
			err := GetBlock(inStream, outStream, "test", tc.address)
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
