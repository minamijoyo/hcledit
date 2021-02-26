package editor

import (
	"github.com/hashicorp/hcl/v2/hclwrite"
)

// BlockRenameFilter is a filter implementation for renaming block.
type BlockRenameFilter struct {
	from string
	to   string
}

var _ Filter = (*BlockRenameFilter)(nil)

// NewBlockRenameFilter creates a new instance of BlockRenameFilter.
func NewBlockRenameFilter(from string, to string) Filter {
	return &BlockRenameFilter{
		from: from,
		to:   to,
	}
}

// Filter reads HCL and renames matched blocks at a given address.
// The blocks which do not match the from address are output as is.
// Rename means setting the block type and labels corresponding to the new
// address.  changing the block type does not make sense on an application
// context, but filters can chain to others and the later filter may edit its
// attributes. So we allow this filter to any block type and labels.
func (f *BlockRenameFilter) Filter(inFile *hclwrite.File) (*hclwrite.File, error) {
	fromTypeName, fromLabels, err := parseAddress(f.from)
	if err != nil {
		return nil, err
	}

	toTypeName, toLabels, err := parseAddress(f.to)
	if err != nil {
		return nil, err
	}

	matched := findBlocks(inFile.Body(), fromTypeName, fromLabels)

	for _, b := range matched {
		b.SetType(toTypeName)
		b.SetLabels(toLabels)
	}

	return inFile, nil
}
