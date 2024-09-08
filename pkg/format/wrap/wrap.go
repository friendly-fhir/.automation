package wrap

import "strings"

// String wraps a string to a maximum column width.
func String(s string) string {
	return defaultWrapper.String(s)
}

// Strings wraps a series of strings representing lines to a maximum column width.
// If any of the 'lines' contain newlines, they will be passed through verbatim.
func Strings(lines ...string) string {
	return defaultWrapper.Strings(lines...)
}

// Lines wraps a slice of strings to a maximum column width.
// If any of the 'lines' contain newlines, they will be passed through verbatim.
func Lines(lines ...string) []string {
	return defaultWrapper.Lines(lines...)
}

type Wrapper struct {
	MaxWidth int
}

func (w *Wrapper) String(s string) string {
	return w.Strings(strings.Split(s, "\n")...)
}

func (w *Wrapper) Strings(lines ...string) string {
	return strings.Join(w.Lines(lines...), "\n")
}

func (w *Wrapper) Lines(lines ...string) []string {
	if len(lines) == 0 {
		return nil
	}
	if len(lines) == 1 && lines[0] == "" {
		return nil
	}
	var newlines []string
	var sb strings.Builder
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			newlines = append(newlines, "")
			continue
		}
		for _, word := range strings.Fields(line) {
			if sb.Len()+len(word)+1 > w.MaxWidth {
				newlines = append(newlines, sb.String())
				sb.Reset()
			}
			if sb.Len() > 0 {
				sb.WriteString(" ")
			}
			sb.WriteString(word)
		}
	}
	if sb.Len() > 0 {
		newlines = append(newlines, sb.String())
		sb.Reset()
	}
	return newlines
}

var defaultWrapper = Wrapper{
	MaxWidth: 120,
}
