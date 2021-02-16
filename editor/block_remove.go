package editor

import (
	"io"

	"github.com/hashicorp/hcl/v2/hclwrite"
)

// RemoveBlock reads HCL from io.Reader, and removes a matched block,
// and writes the updated HCL to io.Writer.
// Note that a filename is used only for an error message.
// If an error occurs, Nothing is written to the output stream.
func RemoveBlock(r io.Reader, w io.Writer, filename string, address string) error {
	filter := &blockRemove{address: address}
	sink := &verticalFormater{}
	return EditHCL(r, w, filename, filter, sink)
}

// blockRemove is a filter implementation for block.
type blockRemove struct {
	address string
}

// Filter reads HCL and removes only matched blocks at a given address.
func (f *blockRemove) Filter(inFile *hclwrite.File) (*hclwrite.File, error) {
	typeName, labels, err := parseAddress(f.address)
	if err != nil {
		return nil, err
	}

	matched := findBlocks(inFile.Body(), typeName, labels)

	for _, b := range matched {
		inFile.Body().RemoveBlock(b)
	}

	return inFile, nil
}
