package editor

import (
	"fmt"
	"strings"

	"github.com/hashicorp/hcl/v2/hclwrite"
)

// BlockGetFilter is a filter implementation for getting block.
type BlockGetFilter struct {
	address string
}

var _ Filter = (*BlockGetFilter)(nil)

// NewBlockGetFilter creates a new instance of BlockGetFilter.
func NewBlockGetFilter(address string) Filter {
	return &BlockGetFilter{
		address: address,
	}
}

// Filter reads HCL and writes only matched blocks at a given address.
func (f *BlockGetFilter) Filter(inFile *hclwrite.File) (*hclwrite.File, error) {
	typeName, labels, err := parseAddress(f.address)
	if err != nil {
		return nil, err
	}

	// find the top-level blocks first.
	matched := findBlocks(inFile.Body(), typeName, labels)
	if len(matched) == 0 {
		// If not found, then find nested blocks.
		// I'll reuse the findLongestMatchingBlocks to implement it as a compromise for now,
		// but it doesn't support the wildcard match. There is a bit inconsistency here.
		// To fix it, we will need to merge implementations of findBlocks and findLongestMatchingBlocks.
		matched, err = findLongestMatchingBlocks(inFile.Body(), f.address)
		if err != nil {
			return nil, err
		}
	}

	outFile := hclwrite.NewEmptyFile()
	for i, b := range matched {
		if i != 0 {
			// when adding a new block, insert a new line before the block.
			outFile.Body().AppendNewline()
		}
		outFile.Body().AppendBlock(b)

	}

	return outFile, nil
}

func parseAddress(address string) (string, []string, error) {
	if len(address) == 0 {
		return "", []string{}, fmt.Errorf("failed to parse address: %s", address)
	}

	a := strings.Split(address, ".")
	typeName := a[0]
	labels := []string{}
	if len(a) > 1 {
		labels = a[1:]
	}
	return typeName, labels, nil
}

// findBlocks returns matching blocks from the body that have the given name
// and labels or returns an empty list if there is currently no matching block.
// The labels can be wildcard (*), but numbers of label must be equal.
func findBlocks(b *hclwrite.Body, typeName string, labels []string) []*hclwrite.Block {
	var matched []*hclwrite.Block
	for _, block := range b.Blocks() {
		if typeName == block.Type() {
			labelNames := block.Labels()
			if len(labels) == 0 && len(labelNames) == 0 {
				matched = append(matched, block)
				continue
			}
			if matchLabels(labels, labelNames) {
				matched = append(matched, block)
			}
		}
	}

	return matched
}

// matchLabels returns true only if the matched and false otherwise.
// The labels can be wildcard (*), but numbers of label must be equal.
func matchLabels(lhs []string, rhs []string) bool {
	if len(lhs) != len(rhs) {
		return false
	}

	for i := range lhs {
		if !(lhs[i] == rhs[i] || lhs[i] == "*" || rhs[i] == "*") {
			return false
		}
	}

	return true
}
