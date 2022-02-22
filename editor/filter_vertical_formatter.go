package editor

import (
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

// verticalFormatterFilter is a Filter implementation to format HCL.
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
type verticalFormatterFilter struct {
}

var _ Filter = (*verticalFormatterFilter)(nil)

// Filter reads HCL and writes formatted contents in vertical.
func (f *verticalFormatterFilter) Filter(inFile *hclwrite.File) (*hclwrite.File, error) {
	tokens := inFile.BuildTokens(nil)
	vertical := VerticalFormat(tokens)

	outFile := hclwrite.NewEmptyFile()
	outFile.Body().AppendUnstructuredTokens(vertical)

	return outFile, nil
}

// VerticalFormat formats token in vertical.
func VerticalFormat(tokens hclwrite.Tokens) hclwrite.Tokens {
	trimmedLeading := trimLeadingNewLine(tokens)
	removed := removeRedundantNewLine(trimmedLeading)
	trimmedTrailing := trimTrailingDuplicatedNewLine(removed)
	return trimmedTrailing
}

// trimLeadingNewLine trims leading newlines from tokens.
func trimLeadingNewLine(tokens hclwrite.Tokens) hclwrite.Tokens {
	begin := 0
	for ; begin < len(tokens); begin++ {
		if tokens[begin].Type != hclsyntax.TokenNewline {
			break
		}
	}

	return tokens[begin:]
}

// trimTrailingDuplicatedNewLine trims trailing newlines from tokens.
// We should not trim the last newlines because the last one means the end of
// line.
func trimTrailingDuplicatedNewLine(tokens hclwrite.Tokens) hclwrite.Tokens {
	end := len(tokens)
	var eof *hclwrite.Token
	for ; end > 1; end-- {
		if tokens[end-1].Type == hclsyntax.TokenEOF {
			// skip EOF
			eof = tokens[end-1]
			continue
		}
		if tokens[end-1].Type == hclsyntax.TokenNewline &&
			tokens[end-2].Type != hclsyntax.TokenNewline {
			break
		}
	}

	ret := tokens[:end]
	if eof != nil {
		// restore EOF
		ret = append(ret, eof)
	}
	return ret
}

// removeRedundantNewLine removes Redundant newlines.
// Two consecutive blank lines should be removed.
// In other words, if there are three consecutive TokenNewline tokens,
// the third and subsequent TokenNewline tokens are removed.
func removeRedundantNewLine(tokens hclwrite.Tokens) hclwrite.Tokens {
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
