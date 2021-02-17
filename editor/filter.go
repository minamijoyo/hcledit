package editor

import (
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

// Filter is an interface which reads HCL and writes HCL
type Filter interface {
	// Filter reads HCL and writes HCL
	Filter(*hclwrite.File) (*hclwrite.File, error)
}

// noop is a Filter inmplementation which does nothing.
type noop struct {
}

// Filter does nothing, just passes a given input to output as it is.
func (f *noop) Filter(inFile *hclwrite.File) (*hclwrite.File, error) {
	return inFile, nil
}

// multiop is a Filter implementation which applies multiple filters in sequence.
type multiop struct {
	filters []Filter
}

// Filter applies multiple filters in sequence.
func (f *multiop) Filter(inFile *hclwrite.File) (*hclwrite.File, error) {
	current := inFile
	for _, f := range f.filters {
		next, err := f.Filter(current)
		if err != nil {
			return nil, err
		}
		current = next
	}
	return current, nil
}

// verticalFormatter is a Filter implementation to format HCL.
// At time of writing, the default hcl formatter does not support vertical
// formatting. However, it's useful in some cases such as removing a block
// because leading and trailing newline tokens don't belong to a block, so
// deleting a block leaves extra newline tokens.
// This is not included in the original hcl implementation, so we should not be
// the default behavior of the formatter not to break existing fomatted hcl configurations.
// Opt-in only where you neeed this feature.
// Note that verticalFormatter formats only in vertical, and not in horizontal.
// This was originally implemented as a Sink, but I found it's better as a Filter,
// because using only default formatter as a Sink is more simple and consistent.
type verticalFormatter struct {
}

// Filter reads HCL and writes formatted contents in vertical.
func (f *verticalFormatter) Filter(inFile *hclwrite.File) (*hclwrite.File, error) {
	tokens := inFile.BuildTokens(nil)
	vertical := VerticalFormat(tokens)

	outFile := hclwrite.NewEmptyFile()
	outFile.Body().AppendUnstructuredTokens(vertical)

	return outFile, nil
}

// VerticalFormat formats token in vertical.
func VerticalFormat(tokens hclwrite.Tokens) hclwrite.Tokens {
	trimmed := trimLeadingNewLine(tokens)
	removed := removeDuplicatedNewLine(trimmed)
	return removed
}

// trimLeadingNewLine trims leading newlines from tokens.
// We don't need trim trailing newlines because the last newline should be
// kept and others are removed removeDuplicatedNewLine.
func trimLeadingNewLine(tokens hclwrite.Tokens) hclwrite.Tokens {
	begin := 0
	for ; begin < len(tokens); begin++ {
		if tokens[begin].Type != hclsyntax.TokenNewline {
			break
		}
	}

	return tokens[begin:len(tokens)]
}

// removeDuplicatedNewLine removes duplicated newlines
// Two consecutive blank lines should be removed.
// In other words, if there are three consecutive TokenNewline tokens,
// the third and subsequent TokenNewline tokens are removed.
func removeDuplicatedNewLine(tokens hclwrite.Tokens) hclwrite.Tokens {
	var removed hclwrite.Tokens
	beforeBefore := false
	before := false

	for _, token := range tokens {
		if token.Type != hclsyntax.TokenNewline {
			removed = append(removed, token)
			// reset
			beforeBefore = false
			before = false
			continue
		}
		// TokenNewLine
		if before && beforeBefore {
			// skip duplicated newlines
			continue
		}
		removed = append(removed, token)
		beforeBefore = before
		before = true
	}

	return removed
}
