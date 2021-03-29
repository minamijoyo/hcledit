package editor

import (
	"io"
)

// Client is an interface for the entrypoint of editor package
type Client interface {
	// Edit reads a HCL file and appies a given filter.
	// If filename is `-`, reads the input from stdin.
	// If update is true, the outputs is written to the input file, else to stdout.
	Edit(filename string, update bool, filter Filter) error
	// Derive reads a HCL file and appies a given sink.
	// If filename is `-`, reads the input from stdin.
	// The outputs is always written to stdout.
	Derive(filename string, sink Sink) error
}

// Option is a set of options for Client.
type Option struct {
	// InStream is the stdin stream.
	InStream io.Reader
	// OutStream is the stdout stream.
	OutStream io.Writer
	// ErrStream is the stderr stream.
	ErrStream io.Writer
}

// client implements the Client interface.
type client struct {
	o *Option
}

var _ Client = (*client)(nil)

// NewClient creates a new instance of Client.
func NewClient(o *Option) Client {
	return &client{
		o: o,
	}
}

// Edit reads a HCL file and appies a given filter.
// If filename is `-`, reads the input from stdin.
// If update is true, the outputs is written to the input file, else to stdout.
func (c *client) Edit(filename string, update bool, filter Filter) error {
	if filename == "-" {
		return EditStream(c.o.InStream, c.o.OutStream, filename, filter)
	}

	if update {
		return UpdateFile(filename, filter)
	}

	return ReadFile(filename, c.o.OutStream, filter)
}

// Derive reads a HCL file and appies a given sink.
// If filename is `-`, reads the input from stdin.
// The outputs is always written to stdout.
func (c *client) Derive(filename string, sink Sink) error {
	if filename == "-" {
		return DeriveStream(c.o.InStream, c.o.OutStream, filename, sink)
	}

	return DeriveFile(filename, c.o.OutStream, sink)
}
