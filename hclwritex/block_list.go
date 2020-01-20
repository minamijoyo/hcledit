package hclwritex

import (
	"io"
	"strings"

	"github.com/hashicorp/hcl/v2/hclwrite"
)

// ListBlock reads HCL from io.Reader, and writes a list of block addresses to io.Writer.
// Note that a filename is used only for an error message.
// If an error occurs, Nothing is written to the output stream.
func ListBlock(r io.Reader, w io.Writer, filename string) error {
	f := &BlockList{}

	return SinkHCL(r, w, filename, f)
}

// BlockList is a Sink implementation to get a list of block address..
type BlockList struct {
}

// Sink reads HCL and writes a list of block addresses.
func (f *BlockList) Sink(inFile *hclwrite.File) (string, error) {
	addrs := []string{}
	for _, b := range inFile.Body().Blocks() {
		addrs = append(addrs, toAddress(b))
	}

	out := strings.Join(addrs, "\n")
	if len(out) != 0 {
		// append a new line if output is not empty.
		out += "\n"
	}
	return out, nil
}

func toAddress(b *hclwrite.Block) string {
	addr := []string{}
	addr = append(addr, b.Type())
	addr = append(addr, (b.Labels())...)
	return strings.Join(addr, ".")
}
