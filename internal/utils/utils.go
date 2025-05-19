package utils

import (
	"fmt"
	"os"
	"strings"
)

// IsDirectory checks if the given path is a directory.
// It takes a string parameter 'path' which is the path to check.
// It returns a boolean indicating whether the path is a directory and an error if any occurs during the check.
func IsDirectory(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			// Path doesn't exist, return specific error
			return false, fmt.Errorf("path '%s' does not exist: %w", path, err)
		}
		// Some other error occurred, return the error.
		return false, fmt.Errorf("error statting path '%s': %w", path, err)
	}

	return fileInfo.IsDir(), nil
}

// IsGitRepository checks if the given path is a git repository.
func IsGitRepository(path string) (bool, error) {
	_, err := os.Stat(path + "/.git")
	if err != nil {
		if os.IsNotExist(err) {
			// Path doesn't exist, return specific error
			return false, fmt.Errorf("path '%s' is not a git repository because the directory does not exist: %w", path, err)
		}
		// Some other error occurred, return the error.
		return false, fmt.Errorf("error statting path '%s': %w", path, err)
	}
	return true, nil
}

// ExtractCodeBlock extracts the string within a markdown code block.
// It takes a string parameter 'markdown' which is the markdown string to extract code block from.
// It returns the extracted code block and an error if any occurs during extraction.
func ExtractCodeBlock(markdown string) (string, error) {
	// Trim the input string to remove leading/trailing whitespace
	markdown = strings.TrimSpace(markdown)

	// Find the position of the first occurrence of ```
	firstStart := strings.Index(markdown, "```")
	if firstStart == -1 {
		return markdown, nil
	}
	firstLineEnd := strings.Index(markdown[firstStart:], "\n")
	if firstLineEnd == -1 {
		return markdown, nil
	}
	firstStart += firstLineEnd

	// Find the position of the last occurrence of ```
	lastEnd := strings.LastIndex(markdown, "```")
	if lastEnd <= firstStart {
		return markdown, nil
	}

	// Extract the content between the first and last ```
	content := markdown[firstStart:lastEnd]

	// Trim any leading/trailing whitespace from the extracted content
	content = strings.TrimSpace(content)

	// Add an empty line at the end of the content
	content += "\n"

	return content, nil
}
