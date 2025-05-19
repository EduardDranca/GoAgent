package utils

import (
	"fmt"
	"github.com/EduardDranca/GoAgent/internal/logging"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"io/fs"
	"os"
	"os/exec" // Add os/exec import
	"path/filepath"
	"strings"
)

// GitUtil interface for Git operations.
type GitUtil interface {
	Add(dir string) error
	Commit(dir string, message string) error
	LsTree(rootDir string) ([]string, error)
	ResetToHead(dir string) error
}

type NoOpGitUtil struct{}

func (g *NoOpGitUtil) Add(_ string) error {
	logging.Logger.Debugf("NoOpGitUtil: Add")
	return nil
}

func (g *NoOpGitUtil) Commit(_ string, _ string) error {
	logging.Logger.Debugf("NoOpGitUtil: Commit")
	return nil
}

func (g *NoOpGitUtil) LsTree(rootDir string) ([]string, error) {
	filePaths := make([]string, 0)
	err := filepath.WalkDir(rootDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			relPath, err := filepath.Rel(rootDir, path)
			if err != nil {
				return err
			}
			filePaths = append(filePaths, relPath)
		}
		return nil
	})
	if err != nil {
		logging.Logger.Errorf("Error walking directory: %v", err)
		return nil, fmt.Errorf("error walking directory: %w", err)
	}
	return filePaths, nil
}

// ResetToHead is a no-op for NoOpGitUtil.
func (g *NoOpGitUtil) ResetToHead(_ string) error {
	logging.Logger.Debugf("NoOpGitUtil: ResetToHead")
	return nil
}

// RealGitUtil implements GitUtil using go-git.
type RealGitUtil struct{}

func (g *RealGitUtil) Add(dir string) error {
	repo, err := git.PlainOpen(dir)
	if err != nil {
		logging.Logger.Errorf("Error opening repository: %v", err)
		return fmt.Errorf("error opening repository: %w", err)
	}

	w, err := repo.Worktree()
	if err != nil {

		logging.Logger.Errorf("Error getting worktree: %v", err)
		return fmt.Errorf("error getting worktree: %w", err)
	}

	err = w.AddGlob(".")
	if err != nil {
		logging.Logger.Errorf("Error adding files to commit: %v", err)
		return fmt.Errorf("error adding files to commit: %w", err)
	}

	logging.Logger.Debugf("Files added to git staging area in %s", dir)

	return nil
}

func (g *RealGitUtil) Commit(dir string, message string) error {
	repo, err := git.PlainOpen(dir)
	if err != nil {
		logging.Logger.Errorf("Error opening repository: %v", err)
		return fmt.Errorf("error opening repository: %w", err)
	}

	w, err := repo.Worktree()
	if err != nil {
		logging.Logger.Errorf("Error getting worktree: %v", err)
		return fmt.Errorf("error getting worktree: %w", err)
	}

	status, err := w.Status()
	if err != nil {
		logging.Logger.Errorf("Error getting worktree status: %v", err)
		return fmt.Errorf("error getting worktree status: %w", err)
	}

	if status.IsClean() {
		logging.Logger.Infof("No changes to commit")
		return nil
	}

	commit, err := w.Commit(message, &git.CommitOptions{})

	if err != nil {
		logging.Logger.Errorf("Error committing changes: %v", err)
		return fmt.Errorf("error committing changes: %w", err)
	}

	obj, err := repo.CommitObject(commit)
	if err != nil {
		logging.Logger.Errorf("Error getting commit object: %v", err)
		return fmt.Errorf("error getting commit object: %w", err)
	}

	logging.Logger.Infof("Changes committed successfully with message: %s. Commit hash: %s", message, obj.String())
	return nil
}

func (g *RealGitUtil) LsTree(rootDir string) ([]string, error) {
	repo, err := git.PlainOpen(rootDir)
	if err != nil {
		// If the repository doesn't exist, return an empty list of files.
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		logging.Logger.Errorf("Error opening repository: %+v", err)
		return nil, fmt.Errorf("error opening repository: %w", err)
	}

	// Get the HEAD reference
	h, err := repo.Head()
	if err != nil {
		logging.Logger.Errorf("Error getting HEAD: %+v", err)
		return nil, fmt.Errorf("error getting HEAD: %w", err)
	}

	// Get the commit object
	commit, err := repo.CommitObject(h.Hash())
	if err != nil {
		logging.Logger.Errorf("Error getting commit object: %+v", err)
		return nil, fmt.Errorf("error getting commit object: %w", err)
	}

	// Get the tree object
	tree, err := commit.Tree()
	if err != nil {
		logging.Logger.Errorf("Error getting tree object: %+v", err)
		return nil, fmt.Errorf("error getting tree object: %w", err)
	}

	var files []string

	tree.Files().ForEach(func(f *object.File) error {
		files = append(files, f.Name)
		return nil
	})

	ws, err := repo.Worktree()
	if err != nil {
		return nil, fmt.Errorf("getting worktree: %w", err)
	}

	status, err := ws.Status()
	if err != nil {
		return nil, fmt.Errorf("getting worktree status: %w", err)
	}

	for file, _ := range status {
		if !strings.Contains(strings.Join(files, "\n"), file) { // Use strings.Contains with joined string
			files = append(files, file)
		}
	}

	return files, nil
}

// ResetToHead resets the repository to the HEAD commit.
func (g *RealGitUtil) ResetToHead(dir string) error {
	repo, err := git.PlainOpen(dir)
	if err != nil {
		logging.Logger.Errorf("Error opening repository: %v", err)
		return fmt.Errorf("error opening repository: %w", err)
	}

	w, err := repo.Worktree()
	if err != nil {
		logging.Logger.Errorf("Error getting worktree: %v", err)
		return fmt.Errorf("error getting worktree: %w", err)
	}

	err = w.Reset(&git.ResetOptions{
		Mode: git.HardReset,
	})
	if err != nil {
		logging.Logger.Errorf("Error resetting repository to HEAD: %v", err)
		return fmt.Errorf("error resetting repository to HEAD: %w", err)
	}

	return nil
}

// EditCommitMessage opens the default git editor to edit the commit message.
func EditCommitMessage(initialMessage string) (string, error) {
	// Create a temporary file
	tmpFile, err := os.CreateTemp("", "goagent-commit-msg-*.txt")
	if err != nil {
		return "", fmt.Errorf("failed to create temporary file: %w", err)
	}
	tmpFilePath := tmpFile.Name()
	tmpFile.Close() // Close the file handle immediately after creation

	defer os.Remove(tmpFilePath) // Clean up the temporary file

	// Write the initial message to the temporary file
	if err := os.WriteFile(tmpFilePath, []byte(initialMessage), 0600); err != nil {
		return "", fmt.Errorf("failed to write initial message to temporary file: %w", err)
	}

	// Find the git editor
	editor, err := findGitEditor()
	if err != nil {
		// Fallback to VISUAL or EDITOR env vars
		editor = os.Getenv("VISUAL")
		if editor == "" {
			editor = os.Getenv("EDITOR")
		}
		if editor == "" {
			return "", fmt.Errorf("could not determine git editor. Please set GIT_EDITOR, VISUAL, or EDITOR environment variable, or configure core.editor in git config")
		}
	}

	// Execute the editor
	cmd := exec.Command(editor, tmpFilePath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	logging.Logger.Infof("Opening editor %s for commit message...", editor)
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("editor command failed: %w", err)
	}

	// Read the edited message from the temporary file
	editedContent, err := os.ReadFile(tmpFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to read edited message from temporary file: %w", err)
	}

	// Trim comments (lines starting with #) and leading/trailing whitespace
	lines := strings.Split(string(editedContent), "\n")
	var cleanLines []string
	for _, line := range lines {
		// Trim space before checking for '#' to handle indented comments
		trimmedLine := strings.TrimSpace(line)
		if !strings.HasPrefix(trimmedLine, "#") {
			cleanLines = append(cleanLines, line) // Keep original line content if not a comment
		}
	}
	editedMessage := strings.TrimSpace(strings.Join(cleanLines, "\n"))

	if editedMessage == "" {
		return "", fmt.Errorf("commit message is empty after editing")
	}

	return editedMessage, nil
}

// findGitEditor attempts to find the default git editor using git commands.
func findGitEditor() (string, error) {
	// Try `git var GIT_EDITOR` first
	cmdVar := exec.Command("git", "var", "GIT_EDITOR")
	outputVar, errVar := cmdVar.Output()
	if errVar == nil {
		editor := strings.TrimSpace(string(outputVar))
		if editor != "" {
			return editor, nil
		}
	}

	// If that fails, try `git config core.editor`
	cmdConfig := exec.Command("git", "config", "core.editor")
	outputConfig, errConfig := cmdConfig.Output()
	if errConfig == nil {
		editor := strings.TrimSpace(string(outputConfig))
		if editor != "" {
			return editor, nil
		}
	}

	return "", fmt.Errorf("could not find git editor using git commands")
}
