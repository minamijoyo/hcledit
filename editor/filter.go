package editor

import (
	"github.com/hashicorp/hcl/v2/hclwrite"
)

// Filter is an interface which reads HCL and writes HCL
type Filter interface {
	// Filter reads HCL and writes HCL
	Filter(*hclwrite.File) (*hclwrite.File, error)
}

// noop is a Filter inmplementation which does nothing.
type noop struct {
}

// Filter does nothing, just passes a given input to output as it is.
func (f *noop) Filter(inFile *hclwrite.File) (*hclwrite.File, error) {
	return inFile, nil
}
