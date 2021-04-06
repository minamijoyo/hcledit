package editor

import (
	"bytes"
	"testing"
)

func TestOperatorEditApply(t *testing.T) {
	cases := []struct {
		name    string
		src     string
		address string
		value   string
		ok      bool
		want    string
	}{
		{
			name: "match (formatted)",
			src: `
a0 = v0
a1 = v1
`,
			address: "a0",
			value:   "v2",
			ok:      true,
			want: `
a0 = v2
a1 = v1
`,
		},
		{
			name: "match (unformatted)",
			src: `
a0 = v0
a1= v1
`,
			address: "a0",
			value:   "v2",
			ok:      true,
			want: `
a0 = v2
a1 = v1
`,
		},
		{
			name: "not found (formatted)",
			src: `
a0 = v0
a1 = v1
`,
			address: "a3",
			value:   "v3",
			ok:      true,
			want: `
a0 = v0
a1 = v1
`,
		},
		{
			name: "not found (unformatted)", // skip format
			src: `
a0 = v0
a1= v1
`,
			address: "a3",
			value:   "v3",
			ok:      true,
			want: `
a0 = v0
a1= v1
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
			o := NewEditOperator(NewAttributeSetFilter(tc.address, tc.value))
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

func TestEditStream(t *testing.T) {
	cases := []struct {
		name    string
		src     string
		address string
		value   string
		ok      bool
		want    string
	}{
		{
			name: "match",
			src: `
a0 = v0
a1 = v1
`,
			address: "a0",
			value:   "v2",
			ok:      true,
			want: `
a0 = v2
a1 = v1
`,
		},
		{
			name: "not found",
			src: `
a0 = v0
a1 = v1
`,
			address: "a3",
			value:   "v3",
			ok:      true,
			want: `
a0 = v0
a1 = v1
`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			inStream := bytes.NewBufferString(tc.src)
			outStream := new(bytes.Buffer)
			filter := NewAttributeSetFilter(tc.address, tc.value)
			err := EditStream(inStream, outStream, "test", filter)
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

func TestUpdateFile(t *testing.T) {
	cases := []struct {
		name    string
		src     string
		address string
		value   string
		ok      bool
		want    string
	}{
		{
			name: "match",
			src: `
a0 = v0
a1 = v1
`,
			address: "a0",
			value:   "v2",
			ok:      true,
			want: `
a0 = v2
a1 = v1
`,
		},
		{
			name: "not found",
			src: `
a0 = v0
a1 = v1
`,
			address: "a3",
			value:   "v3",
			ok:      true,
			want: `
a0 = v0
a1 = v1
`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			path := setupTestFile(t, tc.src)
			filter := NewAttributeSetFilter(tc.address, tc.value)
			err := UpdateFile(path, filter)
			if tc.ok && err != nil {
				t.Fatalf("unexpected err = %s", err)
			}

			got := readTestFile(t, path)
			if !tc.ok && err == nil {
				t.Fatalf("expected to return an error, but no error, contents: \n%s", got)
			}

			if got != tc.want {
				t.Fatalf("got:\n%s\nwant:\n%s", got, tc.want)
			}
		})
	}
}

func TestUpdateFileNotFound(t *testing.T) {
	filter := NewAttributeSetFilter("foo", "bar")
	err := UpdateFile("not_found", filter)
	if err == nil {
		t.Error("expected to return an error, but no error")
	}
}

func TestReadFile(t *testing.T) {
	cases := []struct {
		name    string
		src     string
		address string
		value   string
		ok      bool
		want    string
	}{
		{
			name: "match",
			src: `
a0 = v0
a1 = v1
`,
			address: "a0",
			value:   "v2",
			ok:      true,
			want: `
a0 = v2
a1 = v1
`,
		},
		{
			name: "not found",
			src: `
a0 = v0
a1 = v1
`,
			address: "a3",
			value:   "v3",
			ok:      true,
			want: `
a0 = v0
a1 = v1
`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			path := setupTestFile(t, tc.src)
			outStream := new(bytes.Buffer)
			filter := NewAttributeSetFilter(tc.address, tc.value)
			err := ReadFile(path, outStream, filter)
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

func TestReadFileNotFound(t *testing.T) {
	filter := NewAttributeSetFilter("foo", "bar")
	outStream := new(bytes.Buffer)
	err := ReadFile("not_found", outStream, filter)
	if err == nil {
		t.Error("expected to return an error, but no error")
	}
}
