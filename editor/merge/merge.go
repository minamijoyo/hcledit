// Package merge provides functionality to synchronize HCL blocks between two distinct configurations.
package merge

import (
	"slices"

	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
)

// findBlock is a helper to locate a block by its type and exact labels.
func findBlock(body *hclwrite.Body, blockType string, labels []string) *hclwrite.Block {
	for _, b := range body.Blocks() {
		if b.Type() == blockType && slices.Equal(b.Labels(), labels) {
			return b
		}
	}
	return nil
}

// contains checks if a string slice contains a specific value.
func contains(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

// Blocks synchronizes blocks of a specific type from source to target.
// It updates existing blocks and appends missing blocks.
// ignoreAttrs dictates which attributes should not be copied from the source.
// setAttrs dictates which attributes should be forcefully injected into newly added blocks.
func Blocks(source *hclwrite.File, target *hclwrite.File, blockType string, ignoreAttrs []string, setAttrs map[string]string) {
	sourceBody := source.Body()
	targetBody := target.Body()

	for _, srcBlock := range sourceBody.Blocks() {
		if srcBlock.Type() != blockType {
			continue
		}

		labels := srcBlock.Labels()
		tgtBlock := findBlock(targetBody, blockType, labels)

		if tgtBlock != nil {
			// Update existing block
			for attrName, attr := range srcBlock.Body().Attributes() {
				if !contains(ignoreAttrs, attrName) {
					tgtBlock.Body().SetAttributeRaw(attrName, attr.Expr().BuildTokens(nil))
				}
			}
		} else {
			// Add new block
			newBlock := hclwrite.NewBlock(blockType, labels)

			for attrName, attr := range srcBlock.Body().Attributes() {
				if !contains(ignoreAttrs, attrName) {
					newBlock.Body().SetAttributeRaw(attrName, attr.Expr().BuildTokens(nil))
				}
			}

			// Inject target-specific defaults for the new block
			for attrName, attrVal := range setAttrs {
				newBlock.Body().SetAttributeValue(attrName, cty.StringVal(attrVal))
			}

			targetBody.AppendNewline()
			targetBody.AppendBlock(newBlock)
		}
	}
}
