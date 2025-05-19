package utils

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRenderWithGlamour(t *testing.T) {
	testCases := []struct {
		name              string
		style             string
		input             string
		expectedSubstring string
		expectANSI        bool
	}{
		{
			name:              "Ascii Style",
			style:             "ascii",
			input:             "# Heading\n**Bold**",
			expectedSubstring: "Heading",
			expectANSI:        false,
		},
		{
			name:              "Auto Style",
			style:             "auto",
			input:             "# Heading\n**Bold**",
			expectedSubstring: "Heading",
			expectANSI:        false,
		},
		{
			name:              "Dark Style",
			style:             "dark",
			input:             "# Heading\n**Bold**",
			expectedSubstring: "Heading",
			expectANSI:        true,
		},
		{
			name:              "Dracula Style",
			style:             "dracula",
			input:             "# Heading\n**Bold**",
			expectedSubstring: "Heading",
			expectANSI:        true,
		},
		{
			name:              "Tokyo Night Style",
			style:             "tokyo-night",
			input:             "# Heading\n**Bold**",
			expectedSubstring: "Heading",
			expectANSI:        true,
		},
		{
			name:              "Light Style",
			style:             "light",
			input:             "# Heading\n**Bold**",
			expectedSubstring: "Heading",
			expectANSI:        true,
		},
		{
			name:              "Notty Style",
			style:             "notty",
			input:             "# Heading\n**Bold**",
			expectedSubstring: "Heading",
			expectANSI:        false,
		},
		{
			name:              "Pink Style",
			style:             "pink",
			input:             "# Heading\n**Bold**",
			expectedSubstring: "Heading",
			expectANSI:        true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			SetGlamourStylePath(tc.style)

			renderedOutput, err := RenderWithGlamour(tc.input)
			if err != nil {
				t.Fatalf("RenderWithGlamour failed: %v", err)
			}

			assert.Contains(t, renderedOutput, tc.expectedSubstring)

			if tc.expectANSI {
				assert.True(t, strings.Contains(renderedOutput, "\x1b["), "Expected ANSI escape codes in rendered output for style %s", tc.style)
			} else {
				assert.False(t, strings.Contains(renderedOutput, "\x1b["), "Did not expect ANSI escape codes in rendered output for style %s", tc.style)
			}
		})
	}
}
