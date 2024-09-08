package diagnostic

type reporter interface {
	Report(*Diagnostic)
	reporter()
}

// ReporterFunc is a function that reports diagnostics.
type ReporterFunc func(*Diagnostic)

// Report reports a diagnostic.
func (f ReporterFunc) Report(diagnostic *Diagnostic) {
	f(diagnostic)
}

func (f ReporterFunc) reporter() {}

var _ reporter = (*ReporterFunc)(nil)

// Reporter is an interface that defines a type that can report diagnostics.
type Reporter struct {
	reporter     reporter
	includeDebug bool

	debugs   int
	notices  int
	warnings int
	errors   int
}

func (r *Reporter) Debugs() int {
	return r.debugs
}

func (r *Reporter) Notices() int {
	return r.notices
}

func (r *Reporter) Warnings() int {
	return r.warnings
}

func (r *Reporter) Errors() int {
	return r.errors
}
