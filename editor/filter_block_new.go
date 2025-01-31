package editor

import "github.com/hashicorp/hcl/v2/hclwrite"

// BlockNewFilter is a filter implementation for creating a new block
type BlockNewFilter struct {
	blockType string
	labels    []string
}

// Filter reads HCL and creates a new block with the given type and labels.
func (b *BlockNewFilter) Filter(inFile *hclwrite.File) (*hclwrite.File, error) {
	inFile.Body().AppendNewBlock(b.blockType, b.labels)
	return inFile, nil
}

var _ Filter = (*BlockNewFilter)(nil)

func NewBlockNewFilter(blockType string, labels []string) Filter {
	return &BlockNewFilter{
		blockType: blockType,
		labels:    labels,
	}
}
