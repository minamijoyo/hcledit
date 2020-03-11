package editor

import (
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

// Sink is an interface which reads HCL and writes bytes.
type Sink interface {
	// Sink reads HCL and writes bytes.
	Sink(*hclwrite.File) ([]byte, error)
}

// formater is a Sink implementation to format HCL.
type formater struct {
}

// Sink reads HCL and writes formatted contents.
func (f *formater) Sink(inFile *hclwrite.File) ([]byte, error) {
	raw := inFile.BuildTokens(nil).Bytes()
	out := hclwrite.Format(raw)
	return out, nil
}

// verticalFormater is a Sink implementation to format HCL.
// At time of writing, the default hcl formatter does not support vertical
// formatting. However, it's useful in some cases such as removing a block
// because leading and trailing newline tokens don't belong to a block, so
// deleting a block leaves extra newline tokens.
// This is not included in the original hcl implementation, so we should not be
// the default behavior of the formater not to break existing fomatted hcl configurations.
// Opt-in only where you neeed this feature.
// Note that verticalFormatter formats not only in vertical but also horizontal
// because we cannot use multiple Sink implementations at once.
type verticalFormater struct {
}

// Sink reads HCL and writes formatted contents in vertical and horizontal.
func (f *verticalFormater) Sink(inFile *hclwrite.File) ([]byte, error) {
	tokens := inFile.BuildTokens(nil)

	vertical := VerticalFormat(tokens)

	// default horizontal format
	out := hclwrite.Format(vertical.Bytes())
	return out, nil
}

// VerticalFormat formats token in vertical.
func VerticalFormat(tokens hclwrite.Tokens) hclwrite.Tokens {
	trimmed := trimNewLine(tokens)
	removed := removeDuplicatedNewLine(trimmed)
	return removed
}

// trimNewLine trimsleading and trailing newlines from tokens
func trimNewLine(tokens hclwrite.Tokens) hclwrite.Tokens {
	begin := 0
	for ; begin < len(tokens); begin++ {
		if tokens[begin].Type != hclsyntax.TokenNewline {
			break
		}
	}

	end := len(tokens)
	var eof *hclwrite.Token
	for ; end > begin; end-- {
		if tokens[end-1].Type == hclsyntax.TokenEOF {
			// skip EOF
			eof = tokens[end-1]
			continue
		}
		if tokens[end-1].Type != hclsyntax.TokenNewline {
			break
		}
	}

	ret := append(tokens[begin:end], eof)
	return ret
}

// removeDuplicatedNewLine removes duplicated newlines
func removeDuplicatedNewLine(tokens hclwrite.Tokens) hclwrite.Tokens {
	var removed hclwrite.Tokens
	before := false

	for _, token := range tokens {
		if token.Type != hclsyntax.TokenNewline {
			removed = append(removed, token)
			before = false
		} else if !before {
			removed = append(removed, token)
			before = true
		}
	}

	return removed
}
