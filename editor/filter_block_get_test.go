package editor

import (
	"testing"
)

func TestBlockGetFilter(t *testing.T) {
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
		{
			name: "nested block",
			src: `
b1 {
  a1 = v1
  b2 {
    a2 = v2
  }
}
`,
			address: "b1.b2",
			ok:      true,
			want: `b2 {
  a2 = v2
}
`,
		},
		{
			name: "nested block (extra labels)",
			src: `
b1 "l1" {
  a1 = v1
  b2 {
    a2 = v2
  }
}
`,
			address: "b1.b2",
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
			address: "b1.b2",
			ok:      true,
			want: `b1 "b2" {
  a1 = v1
  b2 {
    a1 = v2
  }
}
`,
		},
		{
			name: "multi level nested block",
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
			address: "b1.b2.b3",
			ok:      true,
			want: `b3 {
  a3 = v3
}
`,
		},
		{
			name: "escaped address",
			src: `
b1 "b.2" {
  a1 = v1
  b2 {
    a1 = v2
  }
}
`,
			address: `b1.b\.2`,
			ok:      true,
			want: `b1 "b.2" {
  a1 = v1
  b2 {
    a1 = v2
  }
}
`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			o := NewEditOperator(NewBlockGetFilter(tc.address))
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
