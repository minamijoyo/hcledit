package editor

import (
	"strings"

	"github.com/hashicorp/hcl/v2/hclwrite"
)

// AttributeRemoveFilter is a filter implementation for removing attribute.
type AttributeRemoveFilter struct {
	address string
}

var _ Filter = (*AttributeRemoveFilter)(nil)

// NewAttributeRemoveFilter creates a new instance of AttributeRemoveFilter.
func NewAttributeRemoveFilter(address string) Filter {
	return &AttributeRemoveFilter{
		address: address,
	}
}

// Filter reads HCL and remove a matched attribute at a given address.
func (f *AttributeRemoveFilter) Filter(inFile *hclwrite.File) (*hclwrite.File, error) {
	attr, body, err := findAttribute(inFile.Body(), f.address)
	if err != nil {
		return nil, err
	}

	if attr != nil {
		a := strings.Split(f.address, ".")
		attrName := a[len(a)-1]
		body.RemoveAttribute(attrName)
	}

	return inFile, nil
}
