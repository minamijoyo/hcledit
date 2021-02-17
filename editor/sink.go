package editor

import (
	"github.com/hashicorp/hcl/v2/hclwrite"
)

// Sink is an interface which reads HCL and writes bytes.
type Sink interface {
	// Sink reads HCL and writes bytes.
	Sink(*hclwrite.File) ([]byte, error)
}

// formatter is a Sink implementation to format HCL.
type formatter struct {
}

// Sink reads HCL and writes formatted contents.
func (f *formatter) Sink(inFile *hclwrite.File) ([]byte, error) {
	raw := inFile.BuildTokens(nil).Bytes()
	out := hclwrite.Format(raw)
	return out, nil
}
