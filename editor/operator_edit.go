package editor

import (
	"fmt"
	"io"
	"io/ioutil"
)

// EditOperator is an implementation of Operator for editing HCL.
type EditOperator struct {
	source Source
	filter Filter
	sink   Sink
}

var _ Operator = (*EditOperator)(nil)

// NewEditOperator creates a new instance of operator for editing HCL.
// If you want to apply multiple filters, use the MultiFilter to compose them.
func NewEditOperator(filter Filter) Operator {
	return &EditOperator{
		source: NewParserSource(),
		filter: filter,
		sink:   NewFormatterSink(),
	}
}

// Apply reads an input stream, applies some filters and formatter, and writes an output stream.
// The input and output streams contain arbitrary bytes in HCL.
// Note that a filename is used only for an error message.
func (o *EditOperator) Apply(r io.Reader, w io.Writer, filename string) error {
	input, err := ioutil.ReadAll(r)
	if err != nil {
		return fmt.Errorf("failed to read input: %s", err)
	}

	inFile, err := o.source.Source(input, filename)
	if err != nil {
		return err
	}

	tmpFile, err := o.filter.Filter(inFile)
	if err != nil {
		return err
	}

	out, err := o.sink.Sink(tmpFile)
	if err != nil {
		return err
	}

	if _, err := w.Write(out); err != nil {
		return fmt.Errorf("failed to write output: %s", err)
	}

	return nil
}
