package context

import (
	"bufio"
	errors2 "errors"
	"fmt"
	"github.com/EduardDranca/GoAgent/internal/logging"
	"github.com/EduardDranca/GoAgent/internal/utils"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// LocalProgrammingAgentContext retrieves files from a local directory.
type LocalProgrammingAgentContext struct {
	rootDir              string
	currentFileContents  map[string]string // In-memory storage of file contents
	CurrentRepoStructure []string
	changeRequest        string   // You might want to handle this differently for local context
	updatedFiles         []string // Keep track of updated files for flushing
	newFiles             []string // Keep track of new files for flushing

	deletedFiles []string          // Track files to be deleted during FlushChanges
	movedFiles   map[string]string // Track files to be moved during FlushChanges, oldPath -> newPath

	// fileAliasesMutex is used to protect the fileAliases map.
	fileAliasesMutex sync.RWMutex

	// fileAliases is a map of file alias to actual file path.
	fileAliases map[string]string
	gitUtil     utils.GitUtil // GitUtil interface for Git operations
}

// NewLocalProgrammingAgentContext creates a new LocalProgrammingAgentContext.
func NewLocalProgrammingAgentContext(rootDir string, changeRequest string, gitUtil utils.GitUtil) (*LocalProgrammingAgentContext, error) {
	ctx := &LocalProgrammingAgentContext{
		rootDir:              rootDir,
		currentFileContents:  make(map[string]string),
		changeRequest:        changeRequest,
		CurrentRepoStructure: []string{},
		updatedFiles:         []string{},
		newFiles:             []string{},
		deletedFiles:         []string{},              // Initialize deletedFiles
		movedFiles:           make(map[string]string), // Initialize movedFiles
		fileAliases:          make(map[string]string), // Initialize fileAliases
		gitUtil:              gitUtil,
	}

	// Build the repository structure
	if err := ctx.buildRepoStructure(); err != nil {
		return nil, fmt.Errorf("error building repo structure: %w", err)
	}
	return ctx, nil
}

// buildRepoStructure builds the repository structure by walking through the directory.
func (c *LocalProgrammingAgentContext) buildRepoStructure() error {
	files, err := c.gitUtil.LsTree(c.rootDir)
	if err != nil {
		return fmt.Errorf("error executing git ls-tree: %w", err)
	}

	c.CurrentRepoStructure = files
	return nil
}

// GetFileContent retrieves the content of a file from the local file system or from the cache.
func (c *LocalProgrammingAgentContext) GetFileContent(filePath string) (string, bool) {
	// Resolve alias before proceeding
	filePath = c.resolveAlias(filePath)
	// Check if the file was deleted
	for _, deletedFile := range c.deletedFiles {
		if deletedFile == filePath {
			return "The file has been deleted.", false
		}
	}
	// Check if the file content is already cached
	if contents, ok := c.currentFileContents[filePath]; ok {
		return contents, true
	}

	// Read the file content from the file system
	fullPath := filepath.Join(c.rootDir, filePath)
	content, err := os.ReadFile(fullPath)
	if err != nil {
		return "The file does not exist or could not be read.", false
	}

	// Cache the file content
	c.currentFileContents[filePath] = string(content)
	return strings.ToValidUTF8(string(content), " "), true
}

// UpdateFileContent updates the content of a file in the cache.
func (c *LocalProgrammingAgentContext) UpdateFileContent(filePath string, newContents string) {
	filePath = c.resolveAlias(filePath)
	// If the file is not in the cache, add it to the new files
	if _, exists := c.currentFileContents[filePath]; !exists {
		c.newFiles = append(c.newFiles, filePath)
		c.CurrentRepoStructure = append(c.CurrentRepoStructure, filePath)
	} else {
		// If the file is in the cache, add it to the updated files
		c.updatedFiles = append(c.updatedFiles, filePath)
	}
	// Update the file content in the cache
	c.currentFileContents[filePath] = newContents
}

// SearchCode searches for a given query in the repository.
func (c *LocalProgrammingAgentContext) SearchCode(query string) map[string][]int {
	searchResults := make(map[string][]int)
	for _, file := range c.CurrentRepoStructure {
		file = c.resolveAlias(file)
		content, exists := c.GetFileContent(file)
		if !exists {
			continue // Skip if file does not exist or cannot be read
		}

		scanner := bufio.NewScanner(strings.NewReader(content))
		lineNumber := 1
		for scanner.Scan() {
			line := scanner.Text()
			if strings.Contains(line, query) {
				searchResults[file] = append(searchResults[file], lineNumber)
			}
			lineNumber++
		}

		if err := scanner.Err(); err != nil {
			// Handle error, maybe log it or return it as part of the result
			continue // For now, skip to the next file
		}
	}
	return searchResults
}

// GetRepoStructure returns the current repository structure.
func (c *LocalProgrammingAgentContext) GetRepoStructure() []string {
	return c.CurrentRepoStructure
}

// GetChangeRequest returns the current change request.
func (c *LocalProgrammingAgentContext) GetChangeRequest() string {
	return c.changeRequest
}

// Delete marks a file for deletion during FlushChanges.
func (c *LocalProgrammingAgentContext) Delete(filePath string) error {
	c.deletedFiles = append(c.deletedFiles, filePath)
	// Remove from CurrentRepoStructure and currentFileContents
	c.removeFileFromContext(filePath)
	return nil
}

// FlushChanges writes the updated files to the local file system, deletes marked files, and moves marked files.
func (c *LocalProgrammingAgentContext) FlushChanges() error {

	// Delete files marked for deletion
	for _, filePath := range c.deletedFiles {
		// Resolve alias before proceeding
		fullPath := filepath.Join(c.rootDir, filePath)
		err := os.Remove(fullPath)
		if err != nil {
			if errors2.Is(err, os.ErrNotExist) {
				logging.Logger.Infof("File %s does not exist, skipping deletion: %v", filePath, err)
				continue // Skip to the next file
			}
			return fmt.Errorf("error deleting file %s: %w", filePath, err)
		}
		// Remove from CurrentRepoStructure and currentFileContents
		c.removeFileFromContext(filePath)
	}

	// Move files marked for moving
	for oldPath, newPath := range c.movedFiles {
		// Resolve alias before proceeding
		skipFile := false
		for _, deletedPath := range c.deletedFiles {
			if deletedPath == oldPath {
				logging.Logger.Debugf("File %s has been deleted, skipping move: %v", oldPath, deletedPath)
				skipFile = true
			}
		}
		if skipFile {
			continue // Skip to the next file
		}
		oldFullPath := filepath.Join(c.rootDir, oldPath)
		newFullPath := filepath.Join(c.rootDir, newPath)
		logging.Logger.Infof("Moving file from %s to %s", oldFullPath, newFullPath)

		if err := os.MkdirAll(filepath.Dir(newFullPath), 0755); err != nil {
			return fmt.Errorf("error creating directory for new file: %w", err)
		}
		err := os.Rename(oldFullPath, newFullPath)
		if err != nil {
			return fmt.Errorf("error moving file from %s to %s: %w", oldPath, newPath, err)
		}

		// Update CurrentRepoStructure and currentFileContents
		c.updateFilePathsInContext(oldPath, newPath)
	}

	// Create the new files
	for _, path := range c.newFiles {
		fullPath := filepath.Join(c.rootDir, path)
		if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
			return fmt.Errorf("error creating directory for %s: %w", path, err)
		}
		if err := os.WriteFile(fullPath, []byte(c.currentFileContents[path]), 0644); err != nil {
			return fmt.Errorf("error creating file %s: %w", path, err)
		}
	}

	// Write the updated files
	for _, path := range c.updatedFiles {
		fullPath := filepath.Join(c.rootDir, path)
		if err := os.WriteFile(fullPath, []byte(c.currentFileContents[path]), 0644); err != nil {
			return fmt.Errorf("error updating file %s: %w", path, err)
		}
	}

	// Clear the updated and new files
	c.updatedFiles = []string{}
	c.newFiles = []string{}
	c.deletedFiles = []string{}
	c.movedFiles = make(map[string]string)

	// Clear file aliases after flushing
	c.fileAliasesMutex.Lock()
	c.fileAliases = make(map[string]string) // Clear aliases after flushing
	c.fileAliasesMutex.Unlock()

	return nil
}

// SetChangeRequest sets the change request.
func (c *LocalProgrammingAgentContext) SetChangeRequest(changeRequest string) {
	c.changeRequest = changeRequest
}

// MoveFile marks a file for moving during FlushChanges.
func (c *LocalProgrammingAgentContext) MoveFile(oldPath string, newPath string) error {
	// Check if the old file exists

	oldFilePath := filepath.Join(c.rootDir, oldPath)
	if _, err := os.Stat(oldFilePath); errors2.Is(err, os.ErrNotExist) {
		return errors2.New(fmt.Sprintf("file %s does not exist", oldPath))
	}
	c.movedFiles[oldPath] = newPath

	// In-memory update of CurrentRepoStructure and currentFileContents to reflect the move immediately for context consistency.
	c.updateFilePathsInContext(oldPath, newPath)

	// Add alias for moved file
	c.fileAliasesMutex.Lock()
	c.fileAliases[newPath] = oldPath
	c.fileAliasesMutex.Unlock()

	return nil
}

// resolveAlias resolves a file path alias if it exists.
func (c *LocalProgrammingAgentContext) resolveAlias(filePath string) string {
	c.fileAliasesMutex.RLock()
	actualPath, isAlias := c.fileAliases[filePath]
	c.fileAliasesMutex.RUnlock()
	if isAlias {
		return actualPath
	}
	return filePath
}

// removeFileFromContext removes a file from CurrentRepoStructure and currentFileContents.
func (c *LocalProgrammingAgentContext) removeFileFromContext(filePath string) {
	// Remove from CurrentRepoStructure
	for i, path := range c.CurrentRepoStructure {
		if path == filePath {
			c.CurrentRepoStructure = append(c.CurrentRepoStructure[:i], c.CurrentRepoStructure[i+1:]...)
			break
		}
	}
	// Remove from currentFileContents
	delete(c.currentFileContents, filePath)
}

// updateFilePathsInContext updates file paths in CurrentRepoStructure, currentFileContents, updatedFiles, and newFiles after a move.
func (c *LocalProgrammingAgentContext) updateFilePathsInContext(oldPath string, newPath string) {
	// Update CurrentRepoStructure
	for i, path := range c.CurrentRepoStructure {
		if path == oldPath {
			c.CurrentRepoStructure[i] = newPath
			break
		}
	}

	// Update currentFileContents
	if content, exists := c.currentFileContents[oldPath]; exists {
		c.currentFileContents[newPath] = content
		delete(c.currentFileContents, oldPath)
	}

	// Update updatedFiles
	for i, path := range c.updatedFiles {
		if path == oldPath {
			c.updatedFiles[i] = newPath
			break
		}
	}

	// Update newFiles
	for i, path := range c.newFiles {
		if path == oldPath {
			c.newFiles[i] = newPath
			break
		}
	}
}

// Ensure that LocalProgrammingAgentContext implements ProgrammingAgentContext
var _ ProgrammingAgentContext = (*LocalProgrammingAgentContext)(nil)
