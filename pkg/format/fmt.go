package format

import "github.com/friendly-fhir/.automation/pkg/ansi"

var (
	// Link is the ANSI escape sequence for a hyperlink.
	Link = ansi.Format{ansi.Underline, ansi.FGWhite}
)
