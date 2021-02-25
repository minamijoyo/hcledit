package editor

import (
	"fmt"
	"io"
	"io/ioutil"
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

// Apply reads an input stream, applies a given sink for deriving, and writes an output stream.
// The input stream contains arbitrary bytes in HCL,
// and the output stream contains arbitrary bytes in non-HCL.
// Note that a filename is used only for an error message.
func (o *DeriveOperator) Apply(r io.Reader, w io.Writer, filename string) error {
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
