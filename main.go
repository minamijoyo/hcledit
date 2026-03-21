package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/hashicorp/logutils"
	"github.com/minamijoyo/hcledit/cmd"
)

func main() {
	log.SetOutput(logOutput())
	// Sanitize os.Args to prevent CRLF log injection
	sanitizedArgs := make([]string, len(os.Args))
	for i, arg := range os.Args {
		clean := strings.ReplaceAll(arg, "\n", " ")
		clean = strings.ReplaceAll(clean, "\r", "")
		sanitizedArgs[i] = clean
	}
	log.Printf("[INFO] CLI args: %q", sanitizedArgs)
	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func logOutput() io.Writer {
	levels := []logutils.LogLevel{"TRACE", "DEBUG", "INFO", "WARN", "ERROR"}
	minLevel := os.Getenv("HCLEDIT_LOG")

	// default log writer is null device.
	writer := io.Discard
	if minLevel != "" {
		writer = os.Stderr
	}

	filter := &logutils.LevelFilter{
		Levels:   levels,
		MinLevel: logutils.LogLevel(minLevel),
		Writer:   writer,
	}

	return filter
}
