package editor

import (
	"fmt"
	"io"

	"github.com/hashicorp/hcl/v2/hclwrite"
)

// AppendBlock reads HCL from io.Reader, and appends a new child block to
// matched blocks at a given parent block address, and writes the updated HCL
// to io.Writer.
// The child address is relative to parent one.
// If a newline flag is true, it also appends a newline before the new block.
// Note that a filename is used only for an error message.
// If an error occurs, Nothing is written to the output stream.
func AppendBlock(r io.Reader, w io.Writer, filename string, parent string, child string, newline bool) error {
	filter := &blockAppend{
		parent:  parent,
		child:   child,
		newline: newline,
	}
	sink := &formater{}
	return EditHCL(r, w, filename, filter, sink)
}

// blockAppend is a filter implementation for block.
type blockAppend struct {
	parent  string
	child   string
	newline bool
}

// Filter reads HCL and appends only matched blocks at a given address.
func (f *blockAppend) Filter(inFile *hclwrite.File) (*hclwrite.File, error) {
	pTypeName, pLabels, err := parseAddress(f.parent)
	if err != nil {
		return nil, fmt.Errorf("failed to parse parent address: %s", err)
	}

	cTypeName, cLabels, err := parseAddress(f.child)
	if err != nil {
		return nil, fmt.Errorf("failed to parse child address: %s", err)
	}

	matched := findBlocks(inFile.Body(), pTypeName, pLabels)

	for _, b := range matched {
		if f.newline {
			b.Body().AppendNewline()
		}
		b.Body().AppendNewBlock(cTypeName, cLabels)
	}

	return inFile, nil
}
