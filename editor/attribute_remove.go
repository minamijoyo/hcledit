package editor

import (
	"io"
	"strings"

	"github.com/hashicorp/hcl/v2/hclwrite"
)

// RemoveAttribute reads HCL from io.Reader, and removes a matched attribute,
// and writes the updated HCL to io.Writer.
// Note that a filename is used only for an error message.
// If an error occurs, Nothing is written to the output stream.
func RemoveAttribute(r io.Reader, w io.Writer, filename string, address string) error {
	filter := &attributeRemove{address: address}
	return EditHCL(r, w, filename, filter)
}

// attributeRemove is a filter implementation for attribute.
type attributeRemove struct {
	address string
}

// Filter reads HCL and remove a matched attribute at a given address.
func (f *attributeRemove) Filter(inFile *hclwrite.File) (*hclwrite.File, error) {
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
