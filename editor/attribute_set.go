package editor

import (
	"fmt"
	"io"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

// SetAttribute reads HCL from io.Reader, and updates a value of matched
// attribute, and writes the updated HCL to io.Writer.
// Note that a filename is used only for an error message.
// If an error occurs, Nothing is written to the output stream.
func SetAttribute(r io.Reader, w io.Writer, filename string, address string, value string) error {
	filter := &attributeSet{address: address, value: value}
	sink := &formater{}
	return EditHCL(r, w, filename, filter, sink)
}

// attributeSet is a filter implementation for attribute.
type attributeSet struct {
	address string
	value   string
}

// Filter reads HCL and updates a value of matched an attribute at a given address.
func (f *attributeSet) Filter(inFile *hclwrite.File) (*hclwrite.File, error) {
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
