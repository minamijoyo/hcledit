package hclwritex

import (
	"fmt"
	"io"
	"strings"

	"github.com/hashicorp/hcl/v2/hclwrite"
)

// GetBlock reads HCL from io.Reader, and writes matched blocks to io.Writer.
// Note that a filename is used only for an error message.
// If an error occurs, Nothing is written to the output stream.
func GetBlock(r io.Reader, w io.Writer, filename string, address string) error {
	f := &BlockFilter{
		address: address,
	}

	return FilterHCL(r, w, filename, f)
}

// BlockFilter is a filter implementation for block.
type BlockFilter struct {
	address string
}

// Process gets blocks at a given address.
func (f *BlockFilter) Process(inFile *hclwrite.File) (*hclwrite.File, error) {
	typeName, labels, err := parseAddress(f.address)
	if err != nil {
		return nil, err
	}

	b := inFile.Body().FirstMatchingBlock(typeName, labels)

	outFile := hclwrite.NewEmptyFile()
	if b != nil {
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
