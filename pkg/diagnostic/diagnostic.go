package diagnostic

import (
	"iter"
)

// Cursor represents a position in a file.
type Cursor struct {
	Line   int `json:"line"`
	Column int `json:"column"`
}

// Diagnostic represents a diagnostic message that is attached to a specific
// file.
type Diagnostic struct {
	Severity Severity `json:"severity"`
	Title    string   `json:"title"`
	Message  string   `json:"message"`

	File  string `json:"file"`
	Start Cursor `json:"start"`
	End   Cursor `json:"end"`

	Suggestion string `json:"suggestion"`
}

// Checker is an interface that defines a type that can check a file for
type Checker interface {
	// Check checks a file for diagnostics.
	Check(file string) iter.Seq[*Diagnostic]
}

// Checkers is a collection of checkers.
type Checkers []Checker

// Check checks a file for diagnostics using all checkers in the collection.
func (c Checkers) Check(file string) iter.Seq[*Diagnostic] {
	return func(yield func(*Diagnostic) bool) {
		for _, checker := range c {
			for diagnostic := range checker.Check(file) {
				if !yield(diagnostic) {
					return
				}
			}
		}
	}
}
