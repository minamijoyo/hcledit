package editor

import (
	"fmt"

	"github.com/hashicorp/hcl/v2/hclwrite"
)

// AttributeRenameFilter is a filter implementation for renaming attribute.
type AttributeRenameFilter struct {
	from string
	to   string
}

var _ Filter = (*AttributeRenameFilter)(nil)

// NewAttributeRenameFilter creates a new instance of AttributeRenameFilter.
func NewAttributeRenameFilter(from string, to string) Filter {
	return &AttributeRenameFilter{
		from: from,
		to:   to,
	}
}

// Filter reads HCL and renames matched an attribute at a given address.
// The current implementation does not allow moving an attribute across blocks,
// but it accepts addresses as arguments, which allows for future extensions.
func (f *AttributeRenameFilter) Filter(inFile *hclwrite.File) (*hclwrite.File, error) {
	fromAttr, fromBody, err := findAttribute(inFile.Body(), f.from)
	if err != nil {
		return nil, err
	}

	if fromAttr != nil {
		fromBlockAddress, fromAttributeName, err := parseAttributeAddress(f.from)
		if err != nil {
			return nil, err
		}
		toBlockAddress, toAttributeName, err := parseAttributeAddress(f.to)
		if err != nil {
			return nil, err
		}

		if fromBlockAddress == toBlockAddress {
			// The Body.RenameAttribute() returns false if fromName does not exist or
			// toName already exists. However, here, we want to return an error only
			// if toName already exists, so we check it ourselves.
			toAttr := fromBody.GetAttribute(toAttributeName)
			if toAttr != nil {
				return nil, fmt.Errorf("attribute already exists: %s", f.to)
			}

			_ = fromBody.RenameAttribute(fromAttributeName, toAttributeName)
		} else {
			return nil, fmt.Errorf("moving an attribute across blocks has not been implemented yet: %s -> %s", f.from, f.to)
		}
	}

	return inFile, nil
}
