package diagnostic

import "iter"

type Severity string

const (
	SeverityDebug   Severity = "debug"
	SeverityNotice  Severity = "notice"
	SeverityWarning Severity = "warning"
	SeverityError   Severity = "error"
)

var allSeverities = []Severity{
	SeverityDebug,
	SeverityNotice,
	SeverityWarning,
	SeverityError,
}

// Severities returns a sequence of all severities.
func Severities() iter.Seq[Severity] {
	return func(yield func(Severity) bool) {
		for _, severity := range allSeverities {
			if !yield(severity) {
				return
			}
		}
	}
}
