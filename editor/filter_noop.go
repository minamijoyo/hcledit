package editor

import "github.com/hashicorp/hcl/v2/hclwrite"

// noop is a Filter inmplementation which does nothing.
type noop struct {
}

// Filter does nothing, just passes a given input to output as it is.
func (f *noop) Filter(inFile *hclwrite.File) (*hclwrite.File, error) {
	return inFile, nil
}
