package editor

import "github.com/hashicorp/hcl/v2/hclwrite"

// MultiFilter is a Filter implementation which applies multiple filters in sequence.
type MultiFilter struct {
	filters []Filter
}

var _ Filter = (*MultiFilter)(nil)

// NewMultiFilter creates a new instance of MultiFilter.
func NewMultiFilter(filters []Filter) Filter {
	return &MultiFilter{
		filters: filters,
	}
}

// Filter applies multiple filters in sequence.
func (f *MultiFilter) Filter(inFile *hclwrite.File) (*hclwrite.File, error) {
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
