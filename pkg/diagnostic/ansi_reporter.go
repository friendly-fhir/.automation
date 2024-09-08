package diagnostic

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/friendly-fhir/.automation/pkg/ansi"
	"github.com/friendly-fhir/.automation/pkg/format/wrap"
)

var (
	fileFormat  = ansi.Format{ansi.FGDefault, ansi.Underline}
	titleFormat = ansi.Format{ansi.FGBrightWhite, ansi.Bold}
)

// ANSIReporter is a reporter that reports diagnostics to the terminal.
type ANSIReporter struct {
	// W is the writer to write the output to.
	W io.Writer

	// MaxWidth is the maximum width of the output. If MaxWidth is 0, the output
	// is not wrapped.
	MaxWidth int

	// Root is the root directory for files to be relative to. If Root is empty,
	// the path is always absolute.
	Root string
}

// Report reports a diagnostic.
func (r *ANSIReporter) Report(diagnostic *Diagnostic) {
	var sb strings.Builder

	color := r.writeSeverity(&sb, diagnostic)

	hasTitle := true
	if diagnostic.Title != "" {
		_, _ = fmt.Fprintf(&sb, "%v%v%v\n", titleFormat, diagnostic.Title, ansi.Reset)
	} else {
		_, _ = sb.WriteRune('\n')
		hasTitle = false
	}

	digits := maxDigits(diagnostic.Start.Line, diagnostic.End.Line)
	changed := linesChanged(diagnostic)
	spaces := strings.Repeat(" ", digits)

	if diagnostic.File != "" {
		_, _ = fmt.Fprintf(&sb, "%s%v---> ", spaces, color)
		r.writeFileLocation(&sb, diagnostic)
		_, _ = fmt.Fprintf(&sb, "%v\n", ansi.Reset)

		lines := readLines(diagnostic)
		if len(lines) != 0 {
			_, _ = fmt.Fprintf(&sb, "%s %v|%v\n", spaces, color, ansi.Reset)
		}
		for i, line := range lines {
			_, _ = fmt.Fprintf(&sb, "%v%*d %v| %v%s\n", ansi.FGBrightBlack, digits, line.Number, color, ansi.FGBrightBlack, line.Text)
			r.writeUnderline(&sb, i != 0, spaces, color, &line)

			if i == 0 && changed > 2 {
				_, _ = fmt.Fprintf(&sb, "%s%v...%v\n", spaces, color, ansi.FGBrightBlack)
			}
		}
	}

	lines := r.resizeLines(digits+3, diagnostic.Message)
	if hasTitle && len(lines) != 0 {
		_, _ = fmt.Fprintf(&sb, "%s %v|%v\n", spaces, color, ansi.Reset)
	}
	for _, line := range lines {
		_, _ = fmt.Fprintf(&sb, "%s %v| %v%s\n", spaces, color, ansi.Reset, line)
	}
	_, _ = r.W.Write([]byte(sb.String()))
}

func (r *ANSIReporter) writeUnderline(w io.Writer, continuation bool, spaces string, color ansi.Format, line *line) {
	if line.StartCol == 0 || line.EndCol == 0 {
		return
	}

	columns := line.EndCol - line.StartCol + 1
	prefix := linePrefix(line.Text)
	position := line.StartCol - 1
	position = max(position-len(prefix), 0) // Handle the case of tabs so that spacing remains correct
	prefix += strings.Repeat(" ", position)

	var underline string
	if continuation {
		underline = strings.Repeat("~", columns-len(prefix))
	} else if columns == 0 {
		underline = "^"
	} else {
		underline = "^" + strings.Repeat("~", columns-1)
	}
	_, _ = fmt.Fprintf(w, "%s %v| %s%v%s\n", spaces, color, prefix, ansi.FGGreen, underline)
}

func (r *ANSIReporter) maxWidth() int {
	if r.MaxWidth == 0 {
		return 80
	}
	return r.MaxWidth
}

func (r *ANSIReporter) writeSeverity(w io.Writer, diagnostic *Diagnostic) ansi.Format {
	var color ansi.Format
	switch diagnostic.Severity {
	case SeverityDebug:
		color = append(color, ansi.FGMagenta)
		_, _ = fmt.Fprintf(w, "%vdebug:%v ", color, ansi.Reset)
	case SeverityNotice:
		color = append(color, ansi.FGBrightBlue)
		_, _ = fmt.Fprintf(w, "%vnotice:%v ", color, ansi.Reset)
	case SeverityWarning:
		color = append(color, ansi.FGYellow)
		_, _ = fmt.Fprintf(w, "%vwarning:%v ", color, ansi.Reset)
	case SeverityError:
		color = append(color, ansi.FGRed)
		_, _ = fmt.Fprintf(w, "%verror:%v ", color, ansi.Reset)
	}
	return color
}

func (r *ANSIReporter) resizeLines(offset int, message string) []string {
	wrapper := wrap.Wrapper{
		MaxWidth: min(max(r.maxWidth()-offset, 30), 120),
	}
	return wrapper.Lines(strings.Split(message, "\n")...)
}

func (r *ANSIReporter) relativePath(path string) string {
	if r.Root != "" {
		rel, err := filepath.Rel(r.Root, path)
		if err == nil {
			return rel
		}
	}
	return path
}

func (r *ANSIReporter) absPath(path string) string {
	if r.Root != "" {
		abs, err := filepath.Abs(path)
		if err == nil {
			return abs
		}
	}
	return path
}

func (r *ANSIReporter) writeFileLocation(w io.Writer, diagnostic *Diagnostic) {
	if diagnostic.File == "" {
		return
	}
	_, _ = fmt.Fprintf(w, "%v%s", fileFormat, r.relativePath(diagnostic.File))
	if diagnostic.Start.Line != 0 {
		_, _ = fmt.Fprintf(w, ":%d", diagnostic.Start.Line)
		if diagnostic.End.Line != 0 && diagnostic.End.Line != diagnostic.Start.Line {
			_, _ = fmt.Fprintf(w, "-%d", diagnostic.End.Line)
		}
		if diagnostic.Start.Column != 0 {
			_, _ = fmt.Fprintf(w, ":%d", diagnostic.Start.Column)
			if diagnostic.End.Column != 0 && diagnostic.End.Column != diagnostic.Start.Column {
				_, _ = fmt.Fprintf(w, "-%d", diagnostic.End.Column)
			}
		}
	}
}

type line struct {
	Number int
	Text   string

	StartCol, EndCol int
}

func linePrefix(line string) string {
	for i, ch := range line {
		if !unicode.IsSpace(ch) {
			return line[:i]
		}
	}
	return line
}

func resizeLines(content string) []string {
	return strings.Split(content, "\n")
}

type annotation struct {
	line int

	columnStart, columnEnd int
}

func createAnnotations(diagnostic *Diagnostic) []annotation {
	if diagnostic.Start.Line == 0 {
		return nil
	}
	var annotations []annotation

	if diagnostic.Start.Line == diagnostic.End.Line || diagnostic.End.Line == 0 {
		annotations = append(annotations, annotation{
			line:        diagnostic.Start.Line,
			columnStart: diagnostic.Start.Column,
			columnEnd:   diagnostic.End.Column,
		})
	} else {
		last := annotation{
			line:      diagnostic.End.Line,
			columnEnd: diagnostic.End.Column,
		}
		if diagnostic.End.Column != 0 && diagnostic.Start.Column != 0 {
			last.columnStart = 1
		}
		annotations = append(annotations, annotation{
			line:        diagnostic.Start.Line,
			columnStart: diagnostic.Start.Column,
		}, last)
	}
	return annotations
}

func readLines(diagnostic *Diagnostic) []line {
	annotations := createAnnotations(diagnostic)

	file, err := os.Open(diagnostic.File)
	if err != nil {
		return nil
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	var lines []line
	lineNumber := 0
	for scanner.Scan() {
		lineNumber++
		if annotation := getAnnotation(annotations, lineNumber); annotation != nil {
			text := scanner.Text()
			if annotation.columnEnd == 0 {
				annotation.columnEnd = len(text)
			}
			if annotation.columnEnd > len(text) {
				annotation.columnEnd = len(text)
			}
			lines = append(lines, line{
				Number:   lineNumber,
				Text:     text,
				StartCol: annotation.columnStart,
				EndCol:   annotation.columnEnd,
			})
		}
	}
	return lines
}

func getAnnotation(annotations []annotation, line int) *annotation {
	for _, annotation := range annotations {
		if annotation.line == line {
			return &annotation
		}
	}
	return nil
}

func linesChanged(diagnostic *Diagnostic) int {
	if diagnostic.Start.Line == 0 {
		return 0
	}
	if diagnostic.End.Line != 0 {
		return diagnostic.End.Line - diagnostic.Start.Line + 1
	}
	return 1
}

func maxDigits(low, high int) int {
	return max(countDigits(low), countDigits(high))
}

func countDigits(n int) int {
	switch {
	case n < 0:
		return countDigits(-n) + 1
	case n < 10:
		return 1
	case n < 100:
		return 2
	case n < 1000:
		return 3
	case n < 10000:
		return 4
	case n < 100000:
		return 5
	case n < 1000000:
		return 6
	case n < 10000000:
		return 7
	}
	return 0
}

func (r *ANSIReporter) reporter() {}

var _ reporter = (*ANSIReporter)(nil)
