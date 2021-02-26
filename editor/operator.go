package editor

import (
	"io"

	"github.com/hashicorp/hcl/v2/hclwrite"
)

// Operator is an interface which abstracts stream operations.
// The hcledit provides not only operations for editing HCL,
// but also for deriving such as listing.
// They need similar but different implementations.
type Operator interface {
	// Apply reads an input stream, apply some operations, and writes an output stream.
	// The input and output streams contain arbitrary bytes (maybe HCL or not).
	// Note that a filename is used only for an error message.
	Apply(r io.Reader, w io.Writer, filename string) error
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
