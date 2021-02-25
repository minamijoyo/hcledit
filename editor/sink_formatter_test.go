package editor

import (
	"bytes"
	"testing"
)

func TestFormat(t *testing.T) {
	cases := []struct {
		name string
		src  string
		ok   bool
		want string
	}{
		{
			name: "unformatted",
			src: `
  b1   {
  a1 = v1
	a2=v2
}
`,
			ok: true,
			want: `
b1 {
  a1 = v1
  a2 = v2
}
`,
		},
		{
			name: "formatted",
			src: `
b1 {
  a1 = v1
  a2 = v2
}
`,
			ok: true,
			want: `
b1 {
  a1 = v1
  a2 = v2
}
`,
		},
		{
			name: "syntax error",
			src: `
b1 {
  a1 = v1
`,
			ok:   false,
			want: "",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			inStream := bytes.NewBufferString(tc.src)
			outStream := new(bytes.Buffer)
			o := NewDeriveOperator(NewFormatterSink())
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
