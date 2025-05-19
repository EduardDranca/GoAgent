package completer

import (
	"fmt"
	"github.com/EduardDranca/GoAgent/internal/utils"
	"github.com/reeflective/readline"
	"strings"
)

// CurrentWordCompleter implements the AutoCompleter interface. It provides completions
// based on the current word of the input line.
type CurrentWordCompleter struct {
	files []string
}

// NewCurrentWordCompleter creates a new CurrentWordCompleter with a list of files.
// It initializes the completer with the provided list of files to suggest completions from.
func NewCurrentWordCompleter(files []string) (*CurrentWordCompleter, error) {
	return &CurrentWordCompleter{
		files: files,
	}, nil
}

// Do implements the AutoCompleter interface. It takes a line of input and a position,
// extracts the word at the cursor position, and returns suggestions for completing that word.
// The suggestions are returned as a slice of []rune, and the length of the word at the cursor
// position is also returned.
func (c *CurrentWordCompleter) Do(line []rune, pos int) readline.Completions {
	lineStr := string(line) // Convert []rune to string for easier handling

	// Find the start and end of the word at the cursor position
	startPos := 0
	endPos := pos
	for i := 0; i < pos; i++ {
		if line[i] == ' ' {
			startPos = i + 1
		}
	}
	for endPos < len(line) && line[endPos] != ' ' {
		endPos++
	}

	wordAtCursor := lineStr[startPos:endPos]

	// No completion if cursor is in the middle of a word
	if pos != endPos {
		return readline.CompleteValues()
	}

	var suggestions []string
	for _, file := range c.files {
		if strings.HasPrefix(strings.ToLower(file), strings.ToLower(wordAtCursor)) {
			suggestions = append(suggestions, file)
		}
	}

	if len(suggestions) > 0 {
		commonPrefix := longestCommonPrefix(suggestions)
		// If there's a common prefix, and it's not the same as the word already typed (case-insensitive)
		if commonPrefix != "" && strings.ToLower(commonPrefix) != strings.ToLower(wordAtCursor) {
			// Return only the common prefix as a suggestion
			return readline.CompleteValues(commonPrefix).DisplayList()
		}
	}

	// Otherwise, return all matching suggestions
	return readline.CompleteValues(suggestions...).DisplayList()
}

func getCurrentFiles(gitUtil utils.GitUtil, directory string) ([]string, error) {
	files, err := gitUtil.LsTree(directory)

	if err != nil {
		return nil, fmt.Errorf("failed to fetch files in repository: %w", err)
	}

	return files, nil
}

func InitCompleter(directory string, gitUtil utils.GitUtil) (*CurrentWordCompleter, error) {
	files, err := getCurrentFiles(gitUtil, directory)
	if err != nil {
		return nil, fmt.Errorf("failed to get current files for completer: %w", err)
	}

	currentWordCompleter, err := NewCurrentWordCompleter(files)
	if err != nil {
		// This error path from NewCurrentWordCompleter is currently not reachable
		// based on its implementation, but kept for robustness.
		return nil, fmt.Errorf("failed to initialize current word completer: %w", err)
	}
	return currentWordCompleter, nil
}

// longestCommonPrefix finds the longest common prefix among a slice of strings.
func longestCommonPrefix(strs []string) string {
	if len(strs) == 0 {
		return ""
	}
	if len(strs) == 1 {
		return strs[0]
	}

	// Find the shortest string to limit the comparison length
	shortest := strs[0]
	for _, s := range strs {
		if len(s) < len(shortest) {
			shortest = s
		}
	}

	prefix := ""
	for i := 0; i < len(shortest); i++ {
		char := shortest[i]
		match := true
		for j := 1; j < len(strs); j++ {
			if strs[j][i] != char {
				match = false
				break
			}
		}
		if match {
			prefix += string(char)
		} else {
			break
		}
	}
	return prefix
}
