package diagnostic

import "io"

// GitHubReporter is a reporter that reports diagnostics to GitHub.
type GitHubReporter struct {
	W io.Writer
}

// Report reports a diagnostic.
func (r *GitHubReporter) Report(diagnostic *Diagnostic) {

}

func (r *GitHubReporter) reporter() {}

var _ reporter = (*GitHubReporter)(nil)
