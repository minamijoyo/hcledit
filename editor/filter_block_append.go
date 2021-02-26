package editor

import (
	"fmt"

	"github.com/hashicorp/hcl/v2/hclwrite"
)

// BlockAppendFilter is a filter implementation for appending block.
type BlockAppendFilter struct {
	parent  string
	child   string
	newline bool
}

var _ Filter = (*BlockAppendFilter)(nil)

// NewBlockAppendFilter creates a new instance of BlockAppendFilter.
func NewBlockAppendFilter(parent string, child string, newline bool) Filter {
	return &BlockAppendFilter{
		parent:  parent,
		child:   child,
		newline: newline,
	}
}

// Filter reads HCL and appends only matched blocks at a given address.
// The child address is relative to parent one.
// If a newline flag is true, it also appends a newline before the new block.
func (f *BlockAppendFilter) Filter(inFile *hclwrite.File) (*hclwrite.File, error) {
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
