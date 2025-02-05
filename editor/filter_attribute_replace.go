package editor

import (
	"fmt"

	"github.com/hashicorp/hcl/v2/hclwrite"
)

// AttributeReplaceFilter is a filter implementation for replacing attribute.
type AttributeReplaceFilter struct {
	address string
	name    string
	value   string
}

var _ Filter = (*AttributeReplaceFilter)(nil)

// NewAttributeReplaceFilter creates a new instance of AttributeReplaceFilter.
func NewAttributeReplaceFilter(address string, name string, value string) Filter {
	return &AttributeReplaceFilter{
		address: address,
		name:    name,
		value:   value,
	}
}

// Filter reads HCL and replaces both the name and value of matched an
// attribute at a given address.
func (f *AttributeReplaceFilter) Filter(inFile *hclwrite.File) (*hclwrite.File, error) {
	attr, body, err := findAttribute(inFile.Body(), f.address)
	if err != nil {
		return nil, err
	}

	if attr != nil {
		_, fromAttributeName, err := parseAttributeAddress(f.address)
		if err != nil {
			return nil, err
		}
		toAttributeName := f.name

		// The Body.RenameAttribute() returns false if fromName does not exist or
		// toName already exists. However, here, we want to return an error only
		// if toName already exists, so we check it ourselves.
		toAttr := body.GetAttribute(toAttributeName)
		if toAttr != nil {
			return nil, fmt.Errorf("attribute already exists: %s", toAttributeName)
		}

		_ = body.RenameAttribute(fromAttributeName, toAttributeName)

		// To delegate expression parsing to the hclwrite parser,
		// We build a new expression and set back to the attribute by tokens.
		expr, err := buildExpression(toAttributeName, f.value)
		if err != nil {
			return nil, err
		}
		body.SetAttributeRaw(toAttributeName, expr.BuildTokens(nil))
	}

	return inFile, nil
}
