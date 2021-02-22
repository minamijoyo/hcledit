package editor

import (
	"github.com/hashicorp/hcl/v2/hclwrite"
)

// formatter is a Sink implementation to format HCL.
type formatter struct {
}

// Sink reads HCL and writes formatted contents.
func (f *formatter) Sink(inFile *hclwrite.File) ([]byte, error) {
	raw := inFile.BuildTokens(nil).Bytes()
	out := hclwrite.Format(raw)
	return out, nil
}
