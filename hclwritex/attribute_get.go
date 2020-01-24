package hclwritex

import (
	"fmt"
	"io"
	"strings"

	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

// GetAttribute reads HCL from io.Reader, and writes a value to matched
// attribute to io.Writer.
// Note that a filename is used only for an error message.
// If an error occurs, Nothing is written to the output stream.
func GetAttribute(r io.Reader, w io.Writer, filename string, address string) error {
	e := &Editor{
		source: &parser{filename: filename},
		filters: []Filter{
			&attributeGet{address: address},
		},
		sink: &attributeGet{address: address},
	}

	return e.Apply(r, w)
}

// attributeGet is a filter and sink implementation for attribute.
type attributeGet struct {
	address string
}

// Filter reads HCL and writes only matched an attribute at a given address.
func (f *attributeGet) Filter(inFile *hclwrite.File) (*hclwrite.File, error) {
	// Check only top level attribute for now.
	attrName := f.address
	attr := inFile.Body().GetAttribute(attrName)

	outFile := hclwrite.NewEmptyFile()
	if attr != nil {
		outFile.Body().SetAttributeRaw(attrName, attr.BuildTokens(nil))
	}

	return outFile, nil
}

// Sink reads HCL and writes value of attribute.
func (f *attributeGet) Sink(inFile *hclwrite.File) ([]byte, error) {
	attrName := f.address
	attr := inFile.Body().GetAttribute(attrName)
	if attr == nil {
		return []byte{}, nil
	}

	// treat expr as a string without interpreting its meaning.
	out, err := getAttributeValueAsString(attr)
	if err != nil {
		return []byte{}, err
	}

	return []byte(out), nil
}

// getAttributeValueAsString returns a value of Attribute as string.
// There is no way to get value as string directly,
// so we parses tokens of Attribute and build string representation.
func getAttributeValueAsString(attr *hclwrite.Attribute) (string, error) {
	// find TokenEqual
	expr := attr.Expr()
	exprTokens := expr.BuildTokens(nil)
	i := 0
	for exprTokens[i].Type != hclsyntax.TokenEqual {
		i++
	}

	if i == len(exprTokens) {
		return "", fmt.Errorf("failed to find TokenEqual: %#v", attr)
	}

	// append tokens until find TokenComment
	var valueTokens hclwrite.Tokens
	for _, t := range exprTokens[(i + 1):] {
		if t.Type == hclsyntax.TokenComment {
			break
		}
		valueTokens = append(valueTokens, t)
	}

	// TokenIdent records SpaceBefore, but we should ignore it here.
	value := strings.TrimSpace(string(valueTokens.Bytes()))

	// Note that the value may be quoted.
	// Most of the case, we need unquoted string.
	// We should strictly check TokenOQuote / TokenCQuote and unquote TokenQuotedLit.
	// TokenQuotedLit may contain escape sequences.
	// To implement it exactly, it becomes more complicated logic.
	// Trim double quotes in string for now.
	unquoted := strings.Trim(value, `"`)

	return unquoted, nil
}
