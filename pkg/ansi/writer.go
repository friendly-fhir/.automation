package ansi

import (
	"io"
	"regexp"
)

// NoFormat returns a writer that strips ANSI escape codes from the output.
func NoFormat(w io.Writer) io.Writer {
	return &noFormatWriter{W: w}
}

// noFormatWriter is a writer that strips ANSI escape codes from the output.
type noFormatWriter struct {
	W io.Writer
}

// Write writes p to the underlying writer, stripping ANSI escape codes.
func (w *noFormatWriter) Write(p []byte) (n int, err error) {
	return w.W.Write(stripFormat.ReplaceAll(p, nil))
}

// stripFormat is a regular expression that matches ANSI escape codes.
var stripFormat = regexp.MustCompile(`\033\[[0-9;]*m`)
