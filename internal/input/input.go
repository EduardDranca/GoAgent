package input

import (
	"errors"
	"fmt"
	"github.com/EduardDranca/GoAgent/internal/input/completer"
	"github.com/reeflective/readline"
	"io" // Import io for io.EOF
	"os"
	"path/filepath"
)

var (
	// Keep a single, package-level shell instance.
	shell *readline.Shell

	// UserInputGetter is a package-level variable that holds the function to get user input.
	// It defaults to GetUserInput but can be overridden for testing purposes to inject mock user input.
	UserInputGetter func(string) (string, error) = GetUserInput
)

// init initializes the shared readline shell instance when the package is loaded.
func init() {
	// Determine a good history file path
	histFile, err := defaultHistoryFile()
	if err != nil {
		// Fallback or handle error appropriately
		fmt.Fprintf(os.Stderr, "Warning: could not determine history file path: %v. Using temporary file.\n", err)
		histFile = filepath.Join(os.TempDir(), "goagent-history.log")
	}

	// Check if the history file is writable
	if _, err := os.Stat(histFile); err != nil {
		if os.IsNotExist(err) {
			// File doesn't exist, create it
			file, err := os.Create(histFile)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Warning: could not create history file %s: %v. History will not be saved.\n", histFile, err)
				histFile = "" // Reset to empty string to indicate no history file
			} else {
				file.Close()
			}
		} else {
			fmt.Fprintf(os.Stderr, "Warning: could not access history file %s: %v. History will not be saved.\n", histFile, err)
			histFile = "" // Reset to empty string to indicate no history file
		}
	}

	// Create the history source
	// Note: NewHistoryFile automatically creates the file and parent directories if they don't exist.
	history, err := readline.NewHistoryFromFile(histFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: could not create history file %s: %v. History will not be saved.\n", histFile, err)
		// Proceed without history persistence if file creation fails
	}

	shell = readline.NewShell()

	// Add the history source to the shell
	if history != nil {
		// Using "default" as the history source name. You could use different names
		// if you needed multiple, separate histories.
		shell.History.Add("default", history)
		fmt.Printf("Using history file: %s\n", histFile) // Inform user
	}
	shell.Completer = nil
}

// defaultHistoryFile provides a sensible default path for the history file.
func defaultHistoryFile() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	// Using XDG Base Directory Specification guidelines if possible
	dataDir := os.Getenv("XDG_DATA_HOME")
	if dataDir == "" {
		dataDir = filepath.Join(home, ".local", "share")
	}
	appDir := filepath.Join(dataDir, "goagent") // Use your app name
	if err := os.MkdirAll(appDir, 0700); err != nil {
		return "", err
	}
	return filepath.Join(appDir, "history"), nil
}

// GetLocalChangeRequest gets user input specifically for a change request,
// potentially using a specific completer.
func GetLocalChangeRequest(currentWordCompleter *completer.CurrentWordCompleter) (string, error) {
	if shell == nil {
		return "", fmt.Errorf("readline shell not initialized")
	}

	shell.Prompt.Primary(func() string { return "\033[31m>>\033[0m " })

	shell.Completer = currentWordCompleter.Do

	// Read the line using the shared shell instance
	line, err := shell.Readline()

	if err != nil {
		// Check for specific errors handled by readline
		if errors.Is(err, readline.ErrInterrupt) {
			fmt.Println("^C")
			return "", err // Pass interrupt error up
		}
		if err == io.EOF {
			// Optional: Print newline for cleaner EOF handling
			fmt.Println("exit") // Or use shell.Config.GetString("eof-prompt") if set
			return "", err      // Pass EOF error up
		}
		// Wrap other errors for context
		return "", fmt.Errorf("error reading line: %w", err)
	}
	return line, nil
}

// GetUserInput gets generic user input using the shared readline shell.
func GetUserInput(prompt string) (string, error) {
	if shell == nil {
		return "", fmt.Errorf("readline shell not initialized")
	}

	// Read the line using the shared shell instance\
	shell.Prompt.Primary(func() string { return "\033[31m>>\033[0m " + prompt }) // Set prompt
	line, err := shell.Readline()
	if err != nil {
		if errors.Is(err, readline.ErrInterrupt) {
			fmt.Println("^C")
			return "", err
		}
		return "", fmt.Errorf("error reading line: %w", err)
	}

	return line, nil
}
