package text

import "fmt"

const (
	ansiRedCode     = "\x1b[31m"
	ansiGreenCode   = "\x1b[32m"
	ansiYellowCode  = "\x1b[33m"
	ansiBlueCode    = "\x1b[34m"
	ansiMagentaCode = "\x1b[35m"
	ansiCyanCode    = "\x1b[36m"

	boldCode = "\x1b[1m"

	ResetCode = "\x1b[0m"
)

// SGR foreground codes of the active palette. They default to the standard
// ANSI 16-color palette and can be swapped with SetPalette.
var (
	redCode     = ansiRedCode
	greenCode   = ansiGreenCode
	yellowCode  = ansiYellowCode
	blueCode    = ansiBlueCode
	magentaCode = ansiMagentaCode

	// CyanCode is the raw SGR start code for cyan, for callers that wrap
	// output manually and pair it with ResetCode.
	CyanCode = ansiCyanCode
)

type paletteCodes struct {
	red, green, yellow, blue, magenta, cyan string
}

// palettes holds the selectable color palettes. The "light" palette uses
// darker 256-color shades that stay legible on light terminal backgrounds,
// where the bright ANSI defaults (yellow in particular) are nearly invisible.
var palettes = map[string]paletteCodes{
	"default": {
		red:     ansiRedCode,
		green:   ansiGreenCode,
		yellow:  ansiYellowCode,
		blue:    ansiBlueCode,
		magenta: ansiMagentaCode,
		cyan:    ansiCyanCode,
	},
	"light": {
		red:     "\x1b[38;5;124m", // #af0000
		green:   "\x1b[38;5;28m",  // #008700
		yellow:  "\x1b[38;5;136m", // #af8700
		blue:    "\x1b[38;5;25m",  // #005faf
		magenta: "\x1b[38;5;129m", // #af00af
		cyan:    "\x1b[38;5;30m",  // #008787
	},
}

// SetPalette selects the active color palette by name. An empty name selects
// the default ANSI palette. An unknown name returns an error and leaves the
// active palette unchanged.
func SetPalette(name string) error {
	if name == "" {
		name = "default"
	}

	p, ok := palettes[name]
	if !ok {
		return fmt.Errorf("unknown color palette %q (valid palettes: default, light)", name)
	}

	redCode = p.red
	greenCode = p.green
	yellowCode = p.yellow
	blueCode = p.blue
	magentaCode = p.magenta
	CyanCode = p.cyan

	return nil
}

// UseColor determines if package will emit colors.
var UseColor = true

func stylize(startCode, in string) string {
	if UseColor {
		return startCode + in + ResetCode
	}

	return in
}

func Red(in string) string {
	return stylize(redCode, in)
}

func Green(in string) string {
	return stylize(greenCode, in)
}

func Yellow(in string) string {
	return stylize(yellowCode, in)
}

func Cyan(in string) string {
	return stylize(CyanCode, in)
}

func Magenta(in string) string {
	return stylize(magentaCode, in)
}

func Blue(in string) string {
	return stylize(blueCode, in)
}

func Bold(in string) string {
	return stylize(boldCode, in)
}

// ColorHash Colors text using a hashing algorithm. The same text will always produce the
// same color while different text will produce a different color.
func ColorHash(name string) (output string) {
	if !UseColor {
		return name
	}

	var hash uint = 5381

	for i := 0; i < len(name); i++ {
		hash = uint(name[i]) + ((hash << 5) + hash)
	}

	// Index order matches the ANSI codes 31-36 used historically.
	codes := [6]string{redCode, greenCode, yellowCode, blueCode, magentaCode, CyanCode}

	return codes[hash%6] + name + ResetCode
}
