package editor

import (
	"bytes"
	"testing"
)

func TestBlockList(t *testing.T) {
	cases := []struct {
		name string
		src  string
		ok   bool
		want string
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
			ok: true,
			want: `b1
b2.l1
`,
		},
		{
			name: "empty",
			src:  "",
			ok:   true,
			want: "",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			inStream := bytes.NewBufferString(tc.src)
			outStream := new(bytes.Buffer)
			o := NewDeriveOperator(NewBlockListSink())
			err := o.Apply(inStream, outStream, "test")
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
