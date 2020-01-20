package hclwritex

import (
	"fmt"
	"io"
	"io/ioutil"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

// Sink is an interface which reads HCL and writes string
type Sink interface {
	// Sink reads HCL and writes string
	Sink(*hclwrite.File) (string, error)
}

// SinkHCL reads HCL from io.Reader, and writes arbitrary string to io.Writer.
// Note that a filename is used only for an error message.
// If an error occurs, Nothing is written to the output stream.
func SinkHCL(r io.Reader, w io.Writer, filename string, sink Sink) error {
	input, err := ioutil.ReadAll(r)
	if err != nil {
		return fmt.Errorf("failed to read input: %s", err)
	}

	f, err := safeParseConfig(input, filename, hcl.Pos{Line: 1, Column: 1})
	if err != nil {
		return err
	}

	out, err := sink.Sink(f)
	if err != nil {
		return err
	}

	if _, err := w.Write([]byte(out)); err != nil {
		return fmt.Errorf("failed to write output: %s", err)
	}

	return nil
}
