package editor

import (
	"fmt"
	"io"
	"io/ioutil"
)

// deriveOperator is an implementation of Operator for deriving any bytes from HCL.
type deriveOperator struct {
	source Source
	sink   Sink
}

var _ Operator = (*deriveOperator)(nil)

// NewDeriveOperator creates a new instance of operator for deriving any bytes from HCL.
func NewDeriveOperator(sink Sink) Operator {
	return &deriveOperator{
		source: NewParserSource(),
		sink:   sink,
	}
}

// Apply reads an input stream, applies a given sink for deriving, and writes an output stream.
// The input stream contains arbitrary bytes in HCL,
// and the output stream contains arbitrary bytes in non-HCL.
// Note that a filename is used only for an error message.
func (o *deriveOperator) Apply(r io.Reader, w io.Writer, filename string) error {
	input, err := ioutil.ReadAll(r)
	if err != nil {
		return fmt.Errorf("failed to read input: %s", err)
	}

	inFile, err := o.source.Source(input, filename)
	if err != nil {
		return err
	}

	out, err := o.sink.Sink(inFile)
	if err != nil {
		return err
	}

	if _, err := w.Write(out); err != nil {
		return fmt.Errorf("failed to write output: %s", err)
	}

	return nil
}

// DeriveHCL reads HCL from an input stream, applies a given sink,
// and writes arbitrary bytes to an output stream.
// This is intended to be used for the output is not HCL such as a listing operation.
func DeriveHCL(r io.Reader, w io.Writer, filename string, sink Sink) error {
	o := NewDeriveOperator(sink)
	return o.Apply(r, w, filename)
}
