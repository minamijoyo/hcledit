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

// Filter is an interface which edits *hclwrite.File
type Filter interface {
	// Process filters input HCL to output one.
	Process(*hclwrite.File) (*hclwrite.File, error)
}

// FilterHCL reads HCL from io.Reader, and writes filtered contents to io.Writer.
// Note that a filename is used only for an error message.
// If an error occurs, Nothing is written to the output stream.
func FilterHCL(r io.Reader, w io.Writer, filename string, filter Filter) error {
	input, err := ioutil.ReadAll(r)
	if err != nil {
		return fmt.Errorf("failed to read input: %s", err)
	}

	f, err := safeParseConfig(input, filename, hcl.Pos{Line: 1, Column: 1})
	if err != nil {
		return err
	}

	filtered, err := filter.Process(f)
	if err != nil {
		return err
	}

	raw := filtered.BuildTokens(nil).Bytes()
	formatted := hclwrite.Format(raw)

	if _, err := w.Write(formatted); err != nil {
		return fmt.Errorf("failed to write output: %s", err)
	}

	return nil
}

// safeParseConfig parses config and recovers if panic occurs.
// The current hclwrite implementation is no perfect and will panic if
// unparseable input is given. We just treat it as a parse error so as not to
// surprise users of tfupdate.
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
