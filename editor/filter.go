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

// multiop is a Filter implementation which applies multiple filters in sequence.
type multiop struct {
	filters []Filter
}

// Filter applies multiple filters in sequence.
func (f *multiop) Filter(inFile *hclwrite.File) (*hclwrite.File, error) {
	current := inFile
	for _, f := range f.filters {
		next, err := f.Filter(current)
		if err != nil {
			return nil, err
		}
		current = next
	}
	return current, nil
}
