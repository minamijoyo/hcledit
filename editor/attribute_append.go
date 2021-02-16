package editor

import (
	"fmt"
	"io"
	"strings"

	"github.com/hashicorp/hcl/v2/hclwrite"
)

// AppendAttribute reads HCL from io.Reader, and appends a new attribute to a
// given address, and writes the updated HCL to io.Writer.
// If a matched block not found, nothing happens.
// If the given attribute already exists, it returns an error.
// If a newline flag is true, it also appends a newline before the new attribute.
// Note that a filename is used only for an error message.
// If an error occurs, Nothing is written to the output stream.
func AppendAttribute(r io.Reader, w io.Writer, filename string, address string, value string, newline bool) error {
	filter := &attributeAppend{
		address: address,
		value:   value,
		newline: newline,
	}
	sink := &formater{}
	return EditHCL(r, w, filename, filter, sink)
}

// attributeAppend is a filter implementation for attribute.
type attributeAppend struct {
	address string
	value   string
	newline bool
}

// Filter reads HCL and appends a new attribute to a given address.
func (f *attributeAppend) Filter(inFile *hclwrite.File) (*hclwrite.File, error) {
	attrName := f.address
	body := inFile.Body()

	a := strings.Split(f.address, ".")
	if len(a) > 1 {
		// if address contains dots, the last element is an attribute name,
		// and the rest is the address of the block.
		attrName = a[len(a)-1]
		blockAddr := strings.Join(a[:len(a)-1], ".")
		blocks, err := findLongestMatchingBlocks(body, blockAddr)
		if err != nil {
			return nil, err
		}

		if len(blocks) == 0 {
			// not found
			return inFile, nil
		}

		// Use first matching one.
		body = blocks[0].Body()
		if body.GetAttribute(attrName) != nil {
			return nil, fmt.Errorf("attribute already exists: %s", f.address)
		}
	}

	// To delegate expression parsing to the hclwrite parser,
	// We build a new expression and set back to the attribute by tokens.
	expr, err := buildExpression(attrName, f.value)
	if err != nil {
		return nil, err
	}

	if f.newline {
		body.AppendNewline()
	}
	body.SetAttributeRaw(attrName, expr.BuildTokens(nil))

	return inFile, nil
}
