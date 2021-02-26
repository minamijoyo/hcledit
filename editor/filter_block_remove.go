package editor

import (
	"github.com/hashicorp/hcl/v2/hclwrite"
)

// BlockRemoveFilter is a filter implementation for removing block.
type BlockRemoveFilter struct {
	address string
}

var _ Filter = (*BlockRemoveFilter)(nil)

// NewBlockRemoveFilter creates a new instance of BlockRemoveFilter.
func NewBlockRemoveFilter(address string) Filter {
	return &BlockRemoveFilter{
		address: address,
	}
}

// Filter reads HCL and removes only matched blocks at a given address.
func (f *BlockRemoveFilter) Filter(inFile *hclwrite.File) (*hclwrite.File, error) {
	m := NewMultiFilter([]Filter{
		&unformattedBlockRemoveFilter{address: f.address},
		&verticalFormatterFilter{},
	})
	return m.Filter(inFile)
}

// unformattedBlockRemoveFilter is a filter implementation for removing block without formatting.
type unformattedBlockRemoveFilter struct {
	address string
}

var _ Filter = (*unformattedBlockRemoveFilter)(nil)

// Filter reads HCL and removes only matched blocks at a given address without formatting.
func (f *unformattedBlockRemoveFilter) Filter(inFile *hclwrite.File) (*hclwrite.File, error) {
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
