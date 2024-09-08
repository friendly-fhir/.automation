/*
Package ansi provides ANSI escape sequences for formatting text.
*/
package ansi

import (
	"fmt"
)

// Attribute represents an ANSI SGR attribute.
type Attribute uint8

// Format returns the ANSI SGR escape sequence for the attribute.
func (a Attribute) Format(w fmt.State, verb rune) {
	switch verb {
	case 'v':
		fmt.Fprintf(w, "%s%dm", SGRPrefix, uint8(a))
	}
}

// String returns the ANSI SGR escape sequence for the attribute.
func (a Attribute) String() string {
	return fmt.Sprintf("%v", a)
}

var _ fmt.Formatter = (*Attribute)(nil)

// Format represents a collection of ANSI SGR attributes.
type Format []Attribute

// Format returns the ANSI SGR escape sequence for the format.
func (f Format) Format(w fmt.State, verb rune) {
	switch verb {
	case 'v':
		fmt.Fprint(w, SGRPrefix)
		for i, a := range f {
			if i > 0 {
				fmt.Fprint(w, ";")
			}
			fmt.Fprint(w, uint8(a))
		}
		fmt.Fprint(w, "m")
	}
}

// String returns the ANSI SGR escape sequence for the format.
func (f Format) String() string {
	return fmt.Sprintf("%v", f)
}

var _ fmt.Formatter = (*Format)(nil)

const (
	// SGRPrefix is the ANSI SGR escape sequence prefix.
	SGRPrefix = "\033["
)
