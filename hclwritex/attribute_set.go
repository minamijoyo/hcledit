package hclwritex

import (
	"io"
	"strings"

	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
)

// SetAttribute reads HCL from io.Reader, and updates a value of matched
// attribute, and writes the updated HCL to io.Writer.
// Note that a filename is used only for an error message.
// If an error occurs, Nothing is written to the output stream.
func SetAttribute(r io.Reader, w io.Writer, filename string, address string, value string) error {
	e := &Editor{
		source: &parser{filename: filename},
		filters: []Filter{
			&attributeSet{address: address, value: value},
		},
		sink: &formater{},
	}

	return e.Apply(r, w)
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
		// Set a new value as string literal for now.
		// We should be able to handle a variable reference.
		body.SetAttributeValue(attrName, cty.StringVal(f.value))
	}

	return inFile, nil
}
