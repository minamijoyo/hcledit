package editor

import (
	"fmt"
	"io"
	"io/ioutil"
)

// Editor is an interface which assembles a pipeline to edit HCL.
type Editor interface {
	// Edit reads an input stream, edits the contents, and writes an output stream.
	// The input and output streams contain arbitrary bytes (maybe HCL or not).
	// Note that a filename is used only for an error message.
	Edit(r io.Reader, w io.Writer, filename string) error
}

// editor is an implementation of Editor.
type editor struct {
	source Source
	filter Filter
	sink   Sink
}

// NewFilterEditor creates a new instance of editor with a given filter.
func NewFilterEditor(filename string, filter Filter) Editor {
	return &editor{
		source: &parser{},
		filter: filter,
		sink:   &formatter{},
	}
}

// NewSinkEditor creates a new instance of editor with a given sink.
func NewSinkEditor(sink Sink) Editor {
	return &editor{
		source: &parser{},
		filter: &noop{},
		sink:   sink,
	}
}

// Edit reads an input stream, applies some filters, and writes an output stream.
// The input and output streams contain arbitrary bytes (maybe HCL or not).
// Note that a filename is used only for an error message.
func (e *editor) Edit(r io.Reader, w io.Writer, filename string) error {
	input, err := ioutil.ReadAll(r)
	if err != nil {
		return fmt.Errorf("failed to read input: %s", err)
	}

	inFile, err := e.source.Source(input, filename)
	if err != nil {
		return err
	}

	tmpFile, err := e.filter.Filter(inFile)
	if err != nil {
		return err
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
	e := NewFilterEditor(filename, filter)
	return e.Edit(r, w, filename)
}

// SinkHCL reads HCL from an input stream, applies a sink,
// and writes arbitrary bytes to an output stream.
// This is intended to be used for the output is not HCL such as a "list" operation.
func SinkHCL(r io.Reader, w io.Writer, filename string, sink Sink) error {
	e := NewSinkEditor(sink)
	return e.Edit(r, w, filename)
}
