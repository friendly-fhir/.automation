package ansi

const (
	// Reset resets all attributes.
	Reset Attribute = 0

	// Bold makes text bold.
	Bold Attribute = 1

	// Dim makes text dim.
	Dim Attribute = 2

	// Italic makes text italic.
	Italic Attribute = 3

	// Underline makes text underlined.
	Underline Attribute = 4

	FGBlack   Attribute = 30
	FGRed     Attribute = 31
	FGGreen   Attribute = 32
	FGYellow  Attribute = 33
	FGBlue    Attribute = 34
	FGMagenta Attribute = 35
	FGCyan    Attribute = 36
	FGWhite   Attribute = 37

	FGDefault Attribute = 39

	BGBlack   Attribute = 40
	BGRed     Attribute = 41
	BGGreen   Attribute = 42
	BGYellow  Attribute = 43
	BGBlue    Attribute = 44
	BGMagenta Attribute = 45
	BGCyan    Attribute = 46
	BGWhite   Attribute = 47

	BGDefault Attribute = 49

	FGBrightBlack   Attribute = 90
	FGBrightRed     Attribute = 91
	FGBrightGreen   Attribute = 92
	FGBrightYellow  Attribute = 93
	FGBrightBlue    Attribute = 94
	FGBrightMagenta Attribute = 95
	FGBrightCyan    Attribute = 96
	FGBrightWhite   Attribute = 97

	BGBrightBlack   Attribute = 100
	BGBrightRed     Attribute = 101
	BGBrightGreen   Attribute = 102
	BGBrightYellow  Attribute = 103
	BGBrightBlue    Attribute = 104
	BGBrightMagenta Attribute = 105
	BGBrightCyan    Attribute = 106
	BGBrightWhite   Attribute = 107
)
