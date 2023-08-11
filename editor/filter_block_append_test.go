package editor

import (
	"testing"
)

func TestBlockAppendFilter(t *testing.T) {
	cases := []struct {
		name    string
		src     string
		parent  string
		child   string
		newline bool
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
			parent:  "b1",
			child:   "b11",
			newline: false,
			ok:      true,
			want: `
a0 = v0
b1 {
  a2 = v2
  b11 {
  }
}

b2 l1 {
}
`,
		},
		{
			name: "no match",
			src: `
a0 = v0
b1 {
  a2 = v2
}

b2 l1 {
}
`,
			parent:  "not_found",
			child:   "b11",
			newline: false,
			ok:      true,
			want: `
a0 = v0
b1 {
  a2 = v2
}

b2 l1 {
}
`,
		},
		{
			name:    "empty",
			parent:  "",
			child:   "b11",
			newline: false,
			ok:      false,
			want:    "",
		},
		{
			name: "with label",
			src: `
a0 = v0
b1 {
  a2 = v2
}

b1 l1 {
}
`,
			parent:  "b1.l1",
			child:   "b11.l11.l12",
			newline: false,
			ok:      true,
			want: `
a0 = v0
b1 {
  a2 = v2
}

b1 l1 {
  b11 "l11" "l12" {
  }
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
			parent:  "b1.l1",
			child:   "b11.l11.l12",
			newline: false,
			ok:      true,
			want: `
b1 {
}

b1 l1 {
  b11 "l11" "l12" {
  }
}

b1 l1 {
  b11 "l11" "l12" {
  }
}
`,
		},
		{
			name: "append newline",
			src: `
b1 {
  a1 = v1
}
`,
			parent:  "b1",
			child:   "b11",
			newline: true,
			ok:      true,
			want: `
b1 {
  a1 = v1

  b11 {
  }
}
`,
		},
		{
			name: "escaped address",
			src: `
a0 = v0
b1 {
  a2 = v2
}

b1 "l.1" {
}
`,
			parent:  `b1.l\.1`,
			child:   `b11.l\.1.l12`,
			newline: false,
			ok:      true,
			want: `
a0 = v0
b1 {
  a2 = v2
}

b1 "l.1" {
  b11 "l.1" "l12" {
  }
}
`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			o := NewEditOperator(NewBlockAppendFilter(tc.parent, tc.child, tc.newline))
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
