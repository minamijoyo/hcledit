package command

import (
	"io"

	"github.com/mitchellh/cli"
	"github.com/spf13/afero"
)

// Meta are the meta-options that are available on all or most commands.
type Meta struct {
	// UI is a user interface representing input and output.
	UI cli.Ui

	// input is an input device.
	// This is a normally stdin, but can be mocked for testing.
	Input io.Reader

	// Fs is an afero filesystem.
	Fs afero.Fs
}
