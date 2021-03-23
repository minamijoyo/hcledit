package editor

import (
	"github.com/hashicorp/hcl/v2/hclwrite"
)

// Operator is an interface which abstracts stream operations.
// The hcledit provides not only operations for editing HCL,
// but also for deriving such as listing.
// They need similar but different implementations.
type Operator interface {
	// Apply reads input bytes, apply some operations, and writes outputs.
	// The input and output contain arbitrary bytes (maybe HCL or not).
	// Note that a filename is used only for an error message.
	Apply(input []byte, filename string) ([]byte, error)
}

// Source is an interface which reads string and writes HCL
type Source interface {
	// Source parses HCL and returns *hclwrite.File
	// filename is a metadata of input stream and used only for an error message.
	Source(src []byte, filename string) (*hclwrite.File, error)
}

// Filter is an interface which reads HCL and writes HCL
type Filter interface {
	// Filter reads HCL and writes HCL
	Filter(*hclwrite.File) (*hclwrite.File, error)
}

// Sink is an interface which reads HCL and writes bytes.
type Sink interface {
	// Sink reads HCL and writes bytes.
	Sink(*hclwrite.File) ([]byte, error)
}

// Formatter is an interface which reads HCL, formats tokens and writes bytes.
// Formatter has a signature similar to Sink, but they have different features,
// so we distinguish them with types.
type Formatter interface {
	// Format reads HCL, formats tokens and writes bytes.
	Format(*hclwrite.File) ([]byte, error)
}
