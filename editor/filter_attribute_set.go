package editor

import (
	"fmt"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

// AttributeSetFilter is a filter implementation for setting attribute.
type AttributeSetFilter struct {
	address string
	value   string
}

var _ Filter = (*AttributeSetFilter)(nil)

// NewAttributeSetFilter creates a new instance of AttributeSetFilter.
func NewAttributeSetFilter(address string, value string) *AttributeSetFilter {
	return &AttributeSetFilter{
		address: address,
		value:   value,
	}
}

// AttributeName returns the name of the attribute in this filter
func (f *AttributeSetFilter) AttributeName() string {
	a := strings.Split(f.address, ".")
	return a[len(a)-1]
}

// AttributeSubPaths returns a list of paths which together point to the end attribute in this filter
// This is a list paths because Attributes could contain bodies which could have another path.
// Paths in bodies are separated by '.' and bodies in attributes are marked by '='.
// So "block1.block2.attr=block3.attr2" points to an attribute named attr2 in the body called block3,
// which can be found in the body of attr in block2 in block1
// This method would in this case return a list containing "block1.block2.attr", and "block3.attr2"
func (f *AttributeSetFilter) AttributeSubPaths() []string {
	return strings.Split(f.address, "=")
}

// Errors to support pass a list of multiple errors as returned by hcl.Diagnostics
// Replace by errors.Join when we only support >go 1.20
type Errors []error

func ErrsFromError(err error) error {
	if err == nil {
		return err
	}
	return Errors{err}
}

func (errs Errors) Error() string {
	var sErrs []string
	for _, err := range errs {
		sErrs = append(sErrs, err.Error())
	}
	return strings.Join(sErrs, " -> ")
}

func pathCompare(a []string, b []string) bool {
	if len(a) != len(b) {
		return false
	} else {
		for i := range a {
			if strings.Compare(a[i], b[i]) != 0 {
				return false
			}
		}
	}
	return true
}

// bodySetSubAttribute allows for an attribute path with = signs.
// It recursively traverses through the attributes in attributes recursively until it has found the exact attribute
// it then sets the attribute and builds a new body (recursing back)...
func bodySetSubAttribute(body *hclwrite.Body, attrPath []string, tokens hclwrite.Tokens) error {
	switch len(attrPath) {
	case 0:
		return ErrsFromError(fmt.Errorf("invalid path"))
	case 1:
		body.SetAttributeRaw(attrPath[0], tokens)
	default:
		if attribute, _, err := findAttribute(body, attrPath[0]); attribute == nil {
			return ErrsFromError(err)
		} else {
			attributeTokens := attribute.BuildTokens(nil)
			attributeTokens = attributeTokens[3 : len(attributeTokens)-2]
			attributeBody := attributeTokens.Bytes()
			if pseudoFile, diags := hclwrite.ParseConfig(attributeBody, attrPath[0],
				hcl.Pos{Line: 0, Column: 0, Byte: 0}); diags.HasErrors() {
				return Errors(diags.Errs())
			} else {
				if errs := bodySetSubAttribute(pseudoFile.Body(), attrPath[1:], tokens); errs != nil {
					return errs
				}
				attrSubPath := strings.Split(attrPath[0], ".")
				attrName := attrSubPath[len(attrSubPath)-1]
				filter := NewAttributeSetFilter(attrName, fmt.Sprintf("{ %s }", pseudoFile.Bytes()))
				//fmt.Printf("from => \n%s", body.BuildTokens(nil).Bytes())
				blocks := body.Blocks()
				if len(blocks) > 0 {
					for _, block := range blocks {
						blockPath := append([]string{block.Type()}, block.Labels()...)
						if pathCompare(blockPath, attrSubPath[:len(attrSubPath)-1]) {
							filter.FilterBody(block.Body())
						}
					}
				} else {
					filter.FilterBody(body)
				}
				//fmt.Printf("to => \n%s", body.BuildTokens(nil).Bytes())
			}
		}
	}
	return nil
}

// FilterBody reads HCL and updates a value of matched an attribute at a given address.
func (f *AttributeSetFilter) Filter(inFile *hclwrite.File) (*hclwrite.File, error) {
	_, err := f.FilterBody(inFile.Body())
	return inFile, err
}

// FilterBody reads HCL and updates a value of matched an attribute at a given address.
func (f *AttributeSetFilter) FilterBody(inBody *hclwrite.Body) (*hclwrite.Body, error) {
	// To delegate expression parsing to the hclwrite parser,
	// We build a new expression and set back to the attribute by tokens.
	expr, err := buildExpression(f.AttributeName(), f.value)
	if err != nil {
		return nil, err
	}
	return inBody, bodySetSubAttribute(inBody, f.AttributeSubPaths(), expr.BuildTokens(nil))

}

// buildExpression returns a new expressions for a given name and value of attribute.
// At the time of writing this, there is no way to parse expression from string.
// So we generate a temporary config on memory and parse it, and extract a generated expression.
func buildExpression(name string, value string) (*hclwrite.Expression, error) {
	if strings.Contains(name, "=") {
		parts := strings.Split(name, "=")
		name = parts[len(parts)-1]
	}
	src := name + " = " + value
	f, err := safeParseConfig([]byte(src), "generated_by_buildExpression", hcl.Pos{Line: 1, Column: 1})
	if err != nil {
		return nil, fmt.Errorf("failed to build expression at the parse phase: %s", err)
	}

	attr := f.Body().GetAttribute(name)
	if attr == nil {
		return nil, fmt.Errorf("failed to build expression at the get phase. name = %s, value = %s", name, value)
	}

	return attr.Expr(), nil
}
