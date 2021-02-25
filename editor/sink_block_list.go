package editor

import (
	"strings"

	"github.com/hashicorp/hcl/v2/hclwrite"
)

// BlockListSink is a sink implementation for getting a list of block addresses.
type BlockListSink struct {
}

var _ Sink = (*BlockListSink)(nil)

// NewBlockListSink creates a new instance of BlockListSink.
func NewBlockListSink() Sink {
	return &BlockListSink{}
}

// Sink reads HCL and writes a list of block addresses.
func (s *BlockListSink) Sink(inFile *hclwrite.File) ([]byte, error) {
	addrs := []string{}
	for _, b := range inFile.Body().Blocks() {
		addrs = append(addrs, toAddress(b))
	}

	out := strings.Join(addrs, "\n")
	if len(out) != 0 {
		// append a new line if output is not empty.
		out += "\n"
	}
	return []byte(out), nil
}

func toAddress(b *hclwrite.Block) string {
	addr := []string{}
	addr = append(addr, b.Type())
	addr = append(addr, (b.Labels())...)
	return strings.Join(addr, ".")
}
