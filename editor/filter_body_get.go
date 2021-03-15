package editor

import (
	"github.com/hashicorp/hcl/v2/hclwrite"
)

// BodyGetFilter is a filter implementation for getting body of first matched block.
type BodyGetFilter struct {
	address string
}

var _ Filter = (*BodyGetFilter)(nil)

// NewBodyGetFilter creates a new instance of BodyGetFilter.
func NewBodyGetFilter(address string) Filter {
	return &BodyGetFilter{
		address: address,
	}
}

// Filter reads HCL and writes body of first matched block at a given address.
func (f *BodyGetFilter) Filter(inFile *hclwrite.File) (*hclwrite.File, error) {
	m := NewMultiFilter([]Filter{
		NewBlockGetFilter(f.address),
		&firstBodyFilter{},
		// body contains a leading NewLine token, it's natural to trim it.
		&verticalFormatterFilter{},
	})
	return m.Filter(inFile)
}

// firstBodyFilter is a filter implementation for getting body of first block.
type firstBodyFilter struct {
}

var _ Filter = (*firstBodyFilter)(nil)

// Filter reads HCL and writes body of first block.
func (f *firstBodyFilter) Filter(inFile *hclwrite.File) (*hclwrite.File, error) {
	outFile := hclwrite.NewEmptyFile()

	matched := inFile.Body().Blocks()
	if len(matched) > 0 {
		// The current implementation doesn't support index in address format.
		// Merging body of contents for multiple blocks doesn't make sense,
		// so we take the first matched block.
		body := matched[0].Body()
		tokens := body.BuildTokens(nil)
		outFile.Body().AppendUnstructuredTokens(tokens)
	}

	return outFile, nil
}
