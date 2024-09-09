package diagnostic_test

import (
	"io"
	"os"
	"runtime"

	"github.com/friendly-fhir/.automation/pkg/diagnostic"
	"golang.org/x/term"
)

func Example() {
	cwd, _ := os.Getwd()
	reporter := &diagnostic.ANSIReporter{
		W:        os.Stdout,
		MaxWidth: termwidth(os.Stdout),
		Root:     cwd,
	}
	msg := &diagnostic.Diagnostic{
		Severity: diagnostic.SeverityDebug,
		Title:    "This is a test error",
		Message:  "You made me do this. It didn't have to end this way.\nReally, we could do better.",
		File:     currentSource(),
		Start: diagnostic.Cursor{
			Line:   15,
			Column: 2,
		},
		End: diagnostic.Cursor{
			Line:   17,
			Column: 15,
		},
	}
	for severity := range diagnostic.Severities() {
		msg.Severity = severity
		reporter.Report(msg)
	}
}

func termwidth(w io.Writer) int {
	if f, ok := w.(*os.File); ok {
		width, _, err := term.GetSize(int(f.Fd()))
		if err == nil {
			return int(width)
		}
	}
	return 80
}

func currentSource() string {
	_, file, _, _ := runtime.Caller(1)
	return file
}
