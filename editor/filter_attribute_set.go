package editor

import (
	"fmt"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

// AttributeSetFilter is a filter implementation for setting attribute.
type AttributeSetFilter struct {
	address string
	value   string
}

var _ Filter = (*AttributeSetFilter)(nil)

// NewAttributeSetFilter creates a new instance of AttributeSetFilter.
func NewAttributeSetFilter(address string, value string) Filter {
	return &AttributeSetFilter{
		address: address,
		value:   value,
	}
}

// Filter reads HCL and updates a value of matched an attribute at a given address.
func (f *AttributeSetFilter) Filter(inFile *hclwrite.File) (*hclwrite.File, error) {
	attr, body, err := findAttribute(inFile.Body(), f.address)
	if err != nil {
		return nil, err
	}

	if attr != nil {
		a := strings.Split(f.address, ".")
		attrName := a[len(a)-1]

		// To delegate expression parsing to the hclwrite parser,
		// We build a new expression and set back to the attribute by tokens.
		expr, err := buildExpression(attrName, f.value)
		if err != nil {
			return nil, err
		}
		body.SetAttributeRaw(attrName, expr.BuildTokens(nil))
	}

	return inFile, nil
}

// buildExpression returns a new expressions for a given name and value of attribute.
// At the time of wrting this, there is no way to parse expression from string.
// So we generate a temporarily config on memory and parse it, and extract a generated expression.
func buildExpression(name string, value string) (*hclwrite.Expression, error) {
	src := name + " = " + value
	f, err := safeParseConfig([]byte(src), "generated_by_buildExpression", hcl.Pos{Line: 1, Column: 1})
	if err != nil {
		return nil, fmt.Errorf("failed to build expression at the parse phase: %s", err)
	}

	attr := f.Body().GetAttribute(name)
	if attr == nil {
		return nil, fmt.Errorf("failed to build expression at the get phase. name = %s, value = %s", name, value)
	}

	return attr.Expr(), nil
}
