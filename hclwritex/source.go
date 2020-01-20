package hclwritex

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"runtime/debug"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

// Source is an interface which reads string and writes HCL
type Source interface {
	// Sink reads string and writes HCL
	Source([]byte) (*hclwrite.File, error)
}

// ParseHCL is a helper function which reads io.Reader and returns *hclwrite.File
// Note that a filename is used only for an error message.
func ParseHCL(r io.Reader, filename string) (*hclwrite.File, error) {
	input, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("failed to read input: %s", err)
	}

	p := &parser{
		filename: filename,
	}

	return p.Source(input)
}

// parser is a Source implementation to parse HCL.
type parser struct {
	// filename is a metadata of input stream and used only for an error message.
	filename string
}

// Source parses HCL and returns *hclwrite.File
func (p *parser) Source(src []byte) (*hclwrite.File, error) {
	return safeParseConfig(src, p.filename, hcl.Pos{Line: 1, Column: 1})
}

// safeParseConfig parses config and recovers if panic occurs.
// The current hclwrite implementation is no perfect and will panic if
// unparseable input is given. We just treat it as a parse error so as not to
// surprise users.
func safeParseConfig(src []byte, filename string, start hcl.Pos) (f *hclwrite.File, e error) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("[DEBUG] failed to parse input: %s\nstacktrace: %s", filename, string(debug.Stack()))
			// Set a return value from panic recover
			e = fmt.Errorf(`failed to parse input: %s
panic: %s
This may be caused by a bug in the hclwrite parser`, filename, err)
		}
	}()

	f, diags := hclwrite.ParseConfig(src, filename, start)

	if diags.HasErrors() {
		return nil, fmt.Errorf("failed to parse input: %s", diags)
	}

	return f, nil
}
