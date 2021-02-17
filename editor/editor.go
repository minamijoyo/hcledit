package editor

import (
	"fmt"
	"io"
	"io/ioutil"
)

// Editor assembles a pipeline to edit HCL.
type Editor struct {
	source  Source
	filters []Filter
	sink    Sink
}

// Apply reads an input stream, applies some filters, and writes an output stream.
// The input and output streams contain arbitrary bytes (maybe HCL or not).
func (e *Editor) Apply(r io.Reader, w io.Writer) error {
	input, err := ioutil.ReadAll(r)
	if err != nil {
		return fmt.Errorf("failed to read input: %s", err)
	}

	inFile, err := e.source.Source(input)
	if err != nil {
		return err
	}

	tmpFile := inFile
	for _, filter := range e.filters {
		tmpFile, err = filter.Filter(tmpFile)
		if err != nil {
			return err
		}
	}

	out, err := e.sink.Sink(tmpFile)
	if err != nil {
		return err
	}

	if _, err := w.Write(out); err != nil {
		return fmt.Errorf("failed to write output: %s", err)
	}

	return nil
}

// FilterHCL reads HCL from an input stream, applies a filter,
// and writes HCL to an output stream.
func FilterHCL(r io.Reader, w io.Writer, filename string, filter Filter) error {
	e := &Editor{
		source:  &parser{filename: filename},
		filters: []Filter{filter},
		sink:    &formatter{},
	}

	return e.Apply(r, w)
}

// SinkHCL reads HCL from an input stream, applies a sink,
// and writes arbitrary bytes to an output stream.
// This is intended to be used for the output is not HCL such as a "list" operation.
func SinkHCL(r io.Reader, w io.Writer, filename string, sink Sink) error {
	filter := &noop{}
	e := &Editor{
		source:  &parser{filename: filename},
		filters: []Filter{filter},
		sink:    sink,
	}

	return e.Apply(r, w)
}
