package editor

import (
	"github.com/hashicorp/hcl/v2/hclwrite"
)

// DefaultFormatter is a default Formatter implementation for formatting HCL.
type DefaultFormatter struct {
}

var _ Formatter = (*DefaultFormatter)(nil)

// NewDefaultFormatter creates a new instance of DefaultFormatter.
func NewDefaultFormatter() Formatter {
	return &DefaultFormatter{}
}

// Format reads HCL, formats tokens and writes formatted contents.
func (f *DefaultFormatter) Format(inFile *hclwrite.File) ([]byte, error) {
	raw := inFile.BuildTokens(nil).Bytes()
	out := hclwrite.Format(raw)
	return out, nil
}
