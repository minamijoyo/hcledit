package editor

import (
	"github.com/hashicorp/hcl/v2/hclwrite"
)

// FormatterSink is a Sink implementation for formatting HCL.
type FormatterSink struct {
}

var _ Sink = (*FormatterSink)(nil)

// NewFormatterSink creates a new instance of FormatterSink.
func NewFormatterSink() Sink {
	return &FormatterSink{}
}

// Sink reads HCL and writes formatted contents.
func (s *FormatterSink) Sink(inFile *hclwrite.File) ([]byte, error) {
	raw := inFile.BuildTokens(nil).Bytes()
	out := hclwrite.Format(raw)
	return out, nil
}
