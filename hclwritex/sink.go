package hclwritex

import (
	"fmt"
	"io"

	"github.com/hashicorp/hcl/v2/hclwrite"
)

// Sink is an interface which reads HCL and writes bytes.
type Sink interface {
	// Sink reads HCL and writes bytes.
	Sink(*hclwrite.File) ([]byte, error)
}

// WriteFormattedHCL is a helper funciton which writes formatted HCL to io.Writer.
func WriteFormattedHCL(inFile *hclwrite.File, w io.Writer) error {
	f := &formater{}
	out, err := f.Sink(inFile)
	if err != nil {
		return err
	}

	return writeRawBytes(out, w)
}

// formater is a Sink implementation to format HCL.
type formater struct {
}

// Sink reads HCL and writes formatted contents.
func (f *formater) Sink(inFile *hclwrite.File) ([]byte, error) {
	raw := inFile.BuildTokens(nil).Bytes()
	out := hclwrite.Format(raw)
	return out, nil
}

// writeRawBytes is a helper funciton which writes raw bytes to io.Writer.
func writeRawBytes(out []byte, w io.Writer) error {
	if _, err := w.Write(out); err != nil {
		return fmt.Errorf("failed to write output: %s", err)
	}

	return nil
}
