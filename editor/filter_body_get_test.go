package editor

import (
	"testing"
)

func TestBodyGetFilter(t *testing.T) {
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
			want: `a2 = v2
`,
		},
		{
			name: "no match",
			src: `
a0 = v0
b1 {
  a2 = v2
}
`,
			address: "foo",
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
			name: "with label",
			src: `
b1 {
  a1 =  v1
}

b1 l1 {
  a2 = v2
}
`,
			address: "b1.l1",
			ok:      true,
			want: `a2 = v2
`,
		},
		{
			name: "get first block",
			src: `
b1 {
  a1 =  v1
}

b1 l1 {
  a2 = v2
}

b1 l1 {
  a3 = v3
}
`,
			address: "b1.l1",
			ok:      true,
			want: `a2 = v2
`,
		},
		{
			name: "get inside block",
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
			want: `a2 = v2
`,
		},
		{
			name: "get outside block",
			src: `
b1 {
  a1 = v1
  b2 {
    a2 = v2
  }
}
`,
			address: "b1",
			ok:      true,
			want: `a1 = v1
b2 {
  a2 = v2
}
`,
		},
		{
			name: "escaped address",
			src: `
b1 {
  a1 =  v1
}

b1 "l.1" {
  a2 = v2
}
`,
			address: `b1.l\.1`,
			ok:      true,
			want: `a2 = v2
`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			o := NewEditOperator(NewBodyGetFilter(tc.address))
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
