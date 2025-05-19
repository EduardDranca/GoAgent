package utils_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/EduardDranca/GoAgent/internal/utils"
)

func TestIsDirectory_ExistingDirectory(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "testdir")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	isDir, err := utils.IsDirectory(tempDir)
	require.NoError(t, err)
	require.True(t, isDir, "IsDirectory should return true for a directory")
}

func TestIsDirectory_ExistingFile(t *testing.T) {
	tempFile, err := os.CreateTemp("", "testfile")
	require.NoError(t, err)
	tempFile.Close()
	defer os.Remove(tempFile.Name())

	isDir, err := utils.IsDirectory(tempFile.Name())
	require.NoError(t, err)
	require.False(t, isDir, "IsDirectory should return false for a file")
}

func TestIsDirectory_NonExistentPath(t *testing.T) {
	isDir, err := utils.IsDirectory("/nonexistent-path") // Assuming this path doesn't exist
	require.ErrorIs(t, err, os.ErrNotExist)
	require.False(t, isDir, "IsDirectory should return false for a non-existent path")
}

func TestExtractCodeBlock_CodeBlock(t *testing.T) {
	codeStr := "```go\nfunc main() {\n}\n```"
	expectedResult := "func main() {\n}\n"

	result, err := utils.ExtractCodeBlock(codeStr)
	require.NoError(t, err)
	require.Equal(t, expectedResult, result)
}

func TestExtractCodeBlock_NoCodeBlock(t *testing.T) {
	codeStr := "No code block"
	result, err := utils.ExtractCodeBlock(codeStr)
	require.NoError(t, err)
	require.Equal(t, codeStr, result)
}
