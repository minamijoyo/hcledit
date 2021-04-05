package editor

import (
	"bytes"
	"testing"
)

func TestClientEdit(t *testing.T) {
	stdin := `
a0 = v0
a1 = v1
`
	inFile := `
a0 = v0
a2 = v2
`
	filter := NewAttributeSetFilter("a0", "v3")

	cases := []struct {
		name     string
		stdin    string
		inFile   string
		filename string
		update   bool
		ok       bool
		stdout   string
		stderr   string
		outFile  string
	}{
		{
			name:     "stdin",
			stdin:    stdin,
			inFile:   inFile,
			filename: "-",
			update:   false,
			ok:       true,
			stdout: `
a0 = v3
a1 = v1
`,
			stderr:  "",
			outFile: inFile,
		},
		{
			name:     "read file",
			stdin:    "",
			inFile:   inFile,
			filename: "test",
			update:   false,
			ok:       true,
			stdout: `
a0 = v3
a2 = v2
`,
			stderr:  "",
			outFile: inFile,
		},
		{
			name:     "update file",
			stdin:    "",
			inFile:   inFile,
			filename: "test",
			update:   true,
			ok:       true,
			stdout:   "",
			stderr:   "",
			outFile: `
a0 = v3
a2 = v2
`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			path := setupTestFile(t, tc.inFile)
			inStream := bytes.NewBufferString(tc.stdin)
			outStream := new(bytes.Buffer)
			errStream := new(bytes.Buffer)
			o := &Option{
				InStream:  inStream,
				OutStream: outStream,
				ErrStream: errStream,
			}
			c := NewClient(o)
			filename := tc.filename
			if filename != "-" {
				// A test file is generated at runtime. We don't use tc.filename, and use the generated file.
				filename = path
			}
			err := c.Edit(filename, tc.update, filter)
			if tc.ok && err != nil {
				t.Errorf("unexpected err = %s", err)
			}

			gotStdout := outStream.String()
			if !tc.ok && err == nil {
				t.Errorf("expected to return an error, but no error, outStream: \n%s", gotStdout)
			}

			if gotStdout != tc.stdout {
				t.Errorf("unexpected stdout. got:\n%s\nwant:\n%s", gotStdout, tc.stdout)
			}

			gotStderr := errStream.String()
			if gotStderr != tc.stderr {
				t.Errorf("unexpected stderr. got:\n%s\nwant:\n%s", gotStderr, tc.stderr)
			}

			gotOutFile := readTestFile(t, path)
			if gotOutFile != tc.outFile {
				t.Errorf("unexpected outFile. got:\n%s\nwant:\n%s", gotOutFile, tc.outFile)
			}
		})
	}
}

func TestClientDerive(t *testing.T) {
	stdin := `
a0 = v0
a1 = v1
`
	inFile := `
a0 = v3
a2 = v2
`
	sink := NewAttributeGetSink("a0")

	cases := []struct {
		name     string
		stdin    string
		inFile   string
		filename string
		ok       bool
		stdout   string
		stderr   string
		outFile  string
	}{
		{
			name:     "stdin",
			stdin:    stdin,
			inFile:   inFile,
			filename: "-",
			ok:       true,
			stdout:   "v0\n",
			stderr:   "",
			outFile:  inFile,
		},
		{
			name:     "read file",
			stdin:    "",
			inFile:   inFile,
			filename: "test",
			ok:       true,
			stdout:   "v3\n",
			stderr:   "",
			outFile:  inFile,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			path := setupTestFile(t, tc.inFile)
			inStream := bytes.NewBufferString(tc.stdin)
			outStream := new(bytes.Buffer)
			errStream := new(bytes.Buffer)
			o := &Option{
				InStream:  inStream,
				OutStream: outStream,
				ErrStream: errStream,
			}
			c := NewClient(o)
			filename := tc.filename
			if filename != "-" {
				// A test file is generated at runtime. We don't use tc.filename, and use the generated file.
				filename = path
			}
			err := c.Derive(filename, sink)
			if tc.ok && err != nil {
				t.Errorf("unexpected err = %s", err)
			}

			gotStdout := outStream.String()
			if !tc.ok && err == nil {
				t.Errorf("expected to return an error, but no error, outStream: \n%s", gotStdout)
			}

			if gotStdout != tc.stdout {
				t.Errorf("unexpected stdout. got:\n%s\nwant:\n%s", gotStdout, tc.stdout)
			}

			gotStderr := errStream.String()
			if gotStderr != tc.stderr {
				t.Errorf("unexpected stderr. got:\n%s\nwant:\n%s", gotStderr, tc.stderr)
			}

			gotOutFile := readTestFile(t, path)
			if gotOutFile != tc.outFile {
				t.Errorf("unexpected outFile. got:\n%s\nwant:\n%s", gotOutFile, tc.outFile)
			}
		})
	}
}
