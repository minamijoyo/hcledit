package editor

import "io"

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
