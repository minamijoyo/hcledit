package editor

import (
	"bytes"
	"testing"
)

func TestOperatorDeriveApply(t *testing.T) {
	cases := []struct {
		name         string
		src          string
		address      string
		withComments bool
		ok           bool
		want         string
	}{
		{
			name: "match",
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
			name: "not found",
			src: `
a0 = v0
a1 = v1
`,
			address:      "a2",
			withComments: false,
			ok:           true,
			want:         "",
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
			o := NewDeriveOperator(NewAttributeGetSink(tc.address, tc.withComments))
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

func TestDeriveStream(t *testing.T) {
	cases := []struct {
		name         string
		src          string
		address      string
		withComments bool
		ok           bool
		want         string
	}{
		{
			name: "match",
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
			name: "not found",
			src: `
a0 = v0
a1 = v1
`,
			address:      "a2",
			withComments: false,
			ok:           true,
			want:         "",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			inStream := bytes.NewBufferString(tc.src)
			outStream := new(bytes.Buffer)
			sink := NewAttributeGetSink(tc.address, tc.withComments)
			err := DeriveStream(inStream, outStream, "test", sink)
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

func TestDeriveFile(t *testing.T) {
	cases := []struct {
		name         string
		src          string
		address      string
		withComments bool
		value        string
		ok           bool
		want         string
	}{
		{
			name: "match",
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
			name: "not found",
			src: `
a0 = v0
a1 = v1
`,
			address:      "a2",
			withComments: false,
			ok:           true,
			want:         "",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			path := setupTestFile(t, tc.src)
			outStream := new(bytes.Buffer)
			sink := NewAttributeGetSink(tc.address, tc.withComments)
			err := DeriveFile(path, outStream, sink)
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

			input := readTestFile(t, path)
			if input != tc.src {
				t.Fatalf("the input file should be read-only, but changed: \n%s", input)
			}
		})
	}
}

func TestDeriveFileNotFound(t *testing.T) {
	sink := NewAttributeGetSink("foo", false)
	outStream := new(bytes.Buffer)
	err := DeriveFile("not_found", outStream, sink)
	if err == nil {
		t.Error("expected to return an error, but no error")
	}
}
