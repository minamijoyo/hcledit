package editor

import (
	"bytes"
	"fmt"
	"io"
	"os"
)

// EditOperator is an implementation of Operator for editing HCL.
type EditOperator struct {
	source    Source
	filter    Filter
	formatter Formatter
}

var _ Operator = (*EditOperator)(nil)

// NewEditOperator creates a new instance of operator for editing HCL.
// If you want to apply multiple filters, use the MultiFilter to compose them.
func NewEditOperator(filter Filter) Operator {
	return &EditOperator{
		source:    NewParserSource(),
		filter:    filter,
		formatter: NewDefaultFormatter(),
	}
}

// Apply reads input bytes, applies some filters and formatter, and writes output.
// The input and output contain arbitrary bytes in HCL.
// Note that a filename is used only for an error message.
func (o *EditOperator) Apply(input []byte, filename string) ([]byte, error) {
	inFile, err := o.source.Source(input, filename)
	if err != nil {
		return nil, err
	}

	tmpFile, err := o.filter.Filter(inFile)
	if err != nil {
		return nil, err
	}

	output := tmpFile.BuildTokens(nil).Bytes()
	// Skip the formatter if the filter didn't change contents to suppress meaningless diff
	if bytes.Equal(input, output) {
		return output, nil
	}

	return o.formatter.Format(tmpFile)
}

// EditStream is a helper method which builds an EditorOperator from a given
// filter and applies it to stream.
// Note that a filename is used only for an error message.
func EditStream(r io.Reader, w io.Writer, filename string, filter Filter) error {
	input, err := io.ReadAll(r)
	if err != nil {
		return fmt.Errorf("failed to read input: %s", err)
	}

	o := NewEditOperator(filter)
	output, err := o.Apply(input, filename)
	if err != nil {
		return err
	}

	if _, err := w.Write(output); err != nil {
		return fmt.Errorf("failed to write output: %s", err)
	}

	return nil
}

// UpdateFile is a helper method which builds an EditorOperator from a given
// filter and applies it to a single file.
// The outputs are written to the input file in-place.
func UpdateFile(filename string, filter Filter) error {
	input, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to open file: %s", err)
	}

	o := NewEditOperator(filter)
	output, err := o.Apply(input, filename)
	if err != nil {
		return err
	}

	// skip updating the timestamp of file if its contents has no change.
	if bytes.Equal(input, output) {
		return nil
	}

	// Write contents back to source file if changed.
	if err = os.WriteFile(filename, output, os.ModePerm); err != nil {
		return fmt.Errorf("failed to write file: %s", err)
	}

	return nil
}

// ReadFile is a helper method which builds an EditorOperator from a given
// filter and applies it to a single file.
// The outputs are written to stream.
func ReadFile(filename string, w io.Writer, filter Filter) error {
	input, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to open file: %s", err)
	}

	o := NewEditOperator(filter)
	output, err := o.Apply(input, filename)
	if err != nil {
		return err
	}

	if _, err := w.Write(output); err != nil {
		return fmt.Errorf("failed to write output: %s", err)
	}

	return nil
}
