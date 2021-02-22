package editor

import "github.com/hashicorp/hcl/v2/hclwrite"

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
