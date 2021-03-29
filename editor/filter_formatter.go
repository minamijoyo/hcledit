package editor

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

// FormatterFilter is a Filter implementation which applies the default formatter as a Filter.
type FormatterFilter struct {
}

var _ Filter = (*FormatterFilter)(nil)

// NewFormatterFilter creates a new instance of FormatterFilter.
func NewFormatterFilter() Filter {
	return &FormatterFilter{}
}

// Filter applies the default formatter as a Filter.
func (f *FormatterFilter) Filter(inFile *hclwrite.File) (*hclwrite.File, error) {
	formatter := NewDefaultFormatter()
	tmp, err := formatter.Format(inFile)
	if err != nil {
		return nil, err
	}

	// The hclwrite package doesn't provide a token-base interface, so we need to
	// parse it again. It's obviously inefficient, but the only way to match the
	// type signature.
	outFile, err := safeParseConfig(tmp, "generated_by_FormatterFilter", hcl.Pos{Line: 1, Column: 1})
	if err != nil {
		// should never happen.
		return nil, fmt.Errorf("failed to parse formatted bytes: %s", err)
	}

	return outFile, nil
}
