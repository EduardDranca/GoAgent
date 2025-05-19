package utils

import (
	"github.com/charmbracelet/glamour"
)

// glamourStylePath is a package-level variable to store the Glamour style path.
var glamourStylePath = "dracula" // Default to dracula

// SetGlamourStylePath sets the package-level glamourStylePath variable.
func SetGlamourStylePath(stylePath string) {
	glamourStylePath = stylePath
}

// RenderWithGlamour renders markdown text using the glamour library.
// It takes a markdown string as input and returns the rendered string.
// It initializes a glamour.TermRenderer with automatic style detection and word wrapping at 80 characters.
// It then renders the markdown to a buffer and returns the resulting string.
// Error handling is included to catch any issues during renderer creation or markdown rendering.
func RenderWithGlamour(markdown string) (string, error) {
	renderer, err := glamour.NewTermRenderer(
		glamour.WithStylePath(glamourStylePath),
		glamour.WithWordWrap(120),
	)
	if err != nil {
		return "", err
	}

	if out, err := renderer.Render(markdown); err != nil {
		return "", err
	} else {
		return out, nil
	}
}
