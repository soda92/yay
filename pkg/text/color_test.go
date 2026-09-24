//go:build !integration

package text

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TestColorHash is intentionally not parallel because it mutates the
// package-level UseColor variable. No other parallel test in this
// package reads UseColor, so sequential execution is sufficient.
func TestColorHash(t *testing.T) {
	original := UseColor
	defer func() { UseColor = original }()

	UseColor = true
	require.Equal(t, ColorHash("core"), ColorHash("core"))
	require.NotEqual(t, ColorHash("core"), ColorHash("extra"))

	UseColor = false
	require.Equal(t, "core", ColorHash("core"))
}

// TestSetPalette is intentionally not parallel because it mutates the
// package-level palette variables.
func TestSetPalette(t *testing.T) {
	originalColor := UseColor
	t.Cleanup(func() {
		UseColor = originalColor
		require.NoError(t, SetPalette("default"))
	})

	UseColor = true

	// Default palette uses the standard ANSI codes.
	require.NoError(t, SetPalette("default"))
	require.Equal(t, "\x1b[33mx\x1b[0m", Yellow("x"))
	require.Equal(t, "\x1b[31mx\x1b[0m", Red("x"))

	// Empty name is an alias for the default palette.
	require.NoError(t, SetPalette(""))
	require.Equal(t, "\x1b[33mx\x1b[0m", Yellow("x"))

	// The light palette uses darker 256-color codes.
	require.NoError(t, SetPalette("light"))
	require.Equal(t, "\x1b[38;5;136mx\x1b[0m", Yellow("x"))
	require.Equal(t, "\x1b[38;5;124mx\x1b[0m", Red("x"))
	require.Equal(t, "\x1b[38;5;30mx\x1b[0m", Cyan("x"))

	// ColorHash follows the active palette: "core" hashes to the yellow
	// slot (index 2).
	require.Equal(t, "\x1b[38;5;136mcore\x1b[0m", ColorHash("core"))

	// Colors stay disabled regardless of the selected palette.
	UseColor = false
	require.Equal(t, "x", Yellow("x"))
	UseColor = true

	// An unknown palette errors and leaves the active palette untouched.
	require.Error(t, SetPalette("nope"))
	require.Equal(t, "\x1b[38;5;136mx\x1b[0m", Yellow("x"))
}

func TestSetPaletteRestoresDefault(t *testing.T) {
	require.NoError(t, SetPalette("light"))
	require.NoError(t, SetPalette("default"))
	require.Equal(t, "\x1b[33mx\x1b[0m", Yellow("x"))
	require.Equal(t, "\x1b[33mcore\x1b[0m", ColorHash("core"))
}
