package editor

import (
	"io"
)

// Format reads HCL from io.Reader, and writes formatted contents to io.Writer.
// Note that a filename is used only for an error message.
// If an error occurs, nothing is written to the output stream.
func Format(r io.Reader, w io.Writer, filename string) error {
	filter := &noop{}
	return FilterHCL(r, w, filename, filter)
}
