package editor

import (
	"fmt"
	"io"
	"os"
)

// DeriveOperator is an implementation of Operator for deriving any bytes from HCL.
type DeriveOperator struct {
	source Source
	sink   Sink
}

var _ Operator = (*DeriveOperator)(nil)

// NewDeriveOperator creates a new instance of operator for deriving any bytes from HCL.
func NewDeriveOperator(sink Sink) Operator {
	return &DeriveOperator{
		source: NewParserSource(),
		sink:   sink,
	}
}

// Apply reads an input bytes, applies a given sink for deriving, and writes output.
// The input contains arbitrary bytes in HCL,
// and the output contains arbitrary bytes in non-HCL.
// Note that a filename is used only for an error message.
func (o *DeriveOperator) Apply(input []byte, filename string) ([]byte, error) {
	inFile, err := o.source.Source(input, filename)
	if err != nil {
		return nil, err
	}

	return o.sink.Sink(inFile)
}

// DeriveStream is a helper method which builds a DeriveOperator from a given
// sink and applies it to stream.
// Note that a filename is used only for an error message.
func DeriveStream(r io.Reader, w io.Writer, filename string, sink Sink) error {
	input, err := io.ReadAll(r)
	if err != nil {
		return fmt.Errorf("failed to read input: %s", err)
	}

	o := NewDeriveOperator(sink)
	output, err := o.Apply(input, filename)
	if err != nil {
		return err
	}

	if _, err := w.Write(output); err != nil {
		return fmt.Errorf("failed to write output: %s", err)
	}

	return nil
}

// DeriveFile is a helper method which builds an DeriveOperator from a given
// sink and applies it to a single file.
// The outputs are written to stream.
func DeriveFile(filename string, w io.Writer, sink Sink) error {
	input, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to open file: %s", err)
	}

	o := NewDeriveOperator(sink)
	output, err := o.Apply(input, filename)
	if err != nil {
		return err
	}

	if _, err := w.Write(output); err != nil {
		return fmt.Errorf("failed to write output: %s", err)
	}

	return nil
}
