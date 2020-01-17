package hclwritex

import (
	"fmt"
	"io"
	"io/ioutil"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

// Filter is an interface which reads HCL and writes HCL
type Filter interface {
	// Filter reads HCL and writes HCL
	Filter(*hclwrite.File) (*hclwrite.File, error)
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

	filtered, err := filter.Filter(f)
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
