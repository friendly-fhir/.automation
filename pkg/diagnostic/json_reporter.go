package diagnostic

import (
	"encoding/json"
	"fmt"
	"io"
)

// JSONReporter is a reporter that reports diagnostics in JSON format.
type JSONReporter struct {
	W      io.Writer
	Indent bool
}

// Report reports a diagnostic.
func (r *JSONReporter) Report(diagnostic *Diagnostic) {
	var bytes []byte
	var err error
	if r.Indent {
		bytes, err = json.Marshal(diagnostic)
	} else {
		bytes, err = json.MarshalIndent(diagnostic, "", "  ")
	}

	if err == nil {
		_, _ = fmt.Fprintf(r.W, "%s", string(bytes))
	} else {
		// If we encounter an error, serialize the error
		_, _ = fmt.Fprintf(r.W, "%q", err)
	}
}

func (r *JSONReporter) reporter() {}

var _ reporter = (*JSONReporter)(nil)
