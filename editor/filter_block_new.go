package editor

import (
	"fmt"

	"github.com/hashicorp/hcl/v2/hclwrite"
)

// BlockNewFilter is a filter implementation for creating a new block
type BlockNewFilter struct {
	address string
	newline bool
}

// Filter reads HCL and creates a new block with the given type and labels.
func (f *BlockNewFilter) Filter(inFile *hclwrite.File) (*hclwrite.File, error) {
	typeName, labels, err := parseAddress(f.address)
	if err != nil {
		return nil, fmt.Errorf("failed to parse address: %s", err)
	}

	if f.newline {
		inFile.Body().AppendNewline()
	}
	inFile.Body().AppendNewBlock(typeName, labels)
	return inFile, nil
}

var _ Filter = (*BlockNewFilter)(nil)

func NewBlockNewFilter(address string, newline bool) Filter {
	return &BlockNewFilter{
		address: address,
		newline: newline,
	}
}
