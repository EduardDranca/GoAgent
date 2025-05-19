package context

import (
	"github.com/EduardDranca/GoAgent/internal/utils"
	"os"
	"path/filepath"
	"sort"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLocalAgentContext_MoveFile_ReadWithNewPath(t *testing.T) {
	// Create a temporary directory.
	tempDir, err := os.MkdirTemp("", "test-agent-context")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Instantiate a LocalProgrammingAgentContext with the temporary directory.
	ctx, _ := NewLocalProgrammingAgentContext(tempDir, "change request", &utils.NoOpGitUtil{})

	// Create a test file (testfile.txt) within the temporary directory and write initial content to it.
	initialContent := "This is the initial content of the test file."
	testFilePath := filepath.Join(tempDir, "testfile.txt")
	err = os.WriteFile(testFilePath, []byte(initialContent), 0644)
	require.NoError(t, err)

	// Call ctx.MoveFile("testfile.txt", "renamed_testfile.txt") to simulate renaming the file.
	err = ctx.MoveFile("testfile.txt", "renamed_testfile.txt")
	require.NoError(t, err)

	// Verify that ctx.GetFileContent("renamed_testfile.txt") returns the expected content and that the file exists.
	newContent, exists := ctx.GetFileContent("renamed_testfile.txt")
	require.Equal(t, true, exists)
	require.Equal(t, initialContent, newContent)

	newFilePath := filepath.Join(tempDir, "renamed_testfile.txt")
	_, err = os.Stat(newFilePath)
	require.Error(t, err, "file should not exist at new path before flush")

	// Flush changes
	err = ctx.FlushChanges()
	require.NoError(t, err)

	// Verify that ctx.GetFileContent("renamed_testfile.txt") returns the expected content and that the file exists.
	newContent, exists = ctx.GetFileContent("renamed_testfile.txt")
	require.Equal(t, true, exists)
	require.Equal(t, initialContent, newContent)

	_, err = os.Stat(newFilePath)
	require.NoError(t, err, "file should exist at new path after flush")
}

func TestLocalAgentContext_MoveFile_ReadWithOldPathAlias(t *testing.T) {
	// Create a temporary directory.
	tempDir, err := os.MkdirTemp("", "test-agent-context")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Instantiate a LocalProgrammingAgentContext with the temporary directory.
	ctx, _ := NewLocalProgrammingAgentContext(tempDir, "change request", &utils.NoOpGitUtil{})

	// Create a test file (testfile.txt) within the temporary directory and write initial content to it.
	initialContent := "This is the initial content of the test file."
	testFilePath := filepath.Join(tempDir, "testfile.txt")
	err = os.WriteFile(testFilePath, []byte(initialContent), 0644)
	require.NoError(t, err)

	// Call ctx.MoveFile("testfile.txt", "renamed_testfile.txt") to simulate renaming the file.
	err = ctx.MoveFile("testfile.txt", "renamed_testfile.txt")
	require.NoError(t, err)

	// Verify that ctx.GetFileContent("testfile.txt") (using the old path, which should now be an alias) returns the expected content and that the file exists.
	oldContent, exists := ctx.GetFileContent("testfile.txt")
	require.True(t, exists, "file should exist at old path")
	require.Equal(t, initialContent, oldContent, "file content should be the same at old path")

	newFilePath := filepath.Join(tempDir, "renamed_testfile.txt")
	_, err = os.Stat(newFilePath)
	require.Error(t, err, "file should not exist at new path before flush")

	// Flush changes
	err = ctx.FlushChanges()
	require.NoError(t, err)

	// After flushing, the file should exist at the new path, and the old path should no longer be valid.
	_, exists = ctx.GetFileContent("testfile.txt")
	require.False(t, exists, "file should exist at old path")

	_, err = os.Stat(newFilePath)
	require.NoError(t, err, "file should exist at new path after flush")
}

func TestLocalAgentContext_DeleteFile_ExistingFile(t *testing.T) {
	// Create a temporary directory.
	tempDir, err := os.MkdirTemp("", "test-delete-file")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create a file named "deletable_file.txt" within the temporary directory.
	deletableFilePath := filepath.Join(tempDir, "deletable_file.txt")
	_, err = os.Create(deletableFilePath)
	require.NoError(t, err)

	// Instantiate a LocalProgrammingAgentContext.
	ctx, _ := NewLocalProgrammingAgentContext(tempDir, "change request", &utils.NoOpGitUtil{})

	// Call ctx.DeleteFile("deletable_file.txt").
	err = ctx.Delete("deletable_file.txt")
	require.NoError(t, err)

	// Assert that no error is returned and the file "deletable_file.txt" still exists before flush.
	_, err = os.Stat(deletableFilePath)
	require.NoError(t, err, "file should exist before flush")

	// Flush changes
	err = ctx.FlushChanges()
	require.NoError(t, err)

	// Assert that no error is returned and the file "deletable_file.txt" no longer exists.
	_, err = os.Stat(deletableFilePath)
	require.Error(t, err, "file should not exist")
	require.True(t, os.IsNotExist(err), "error should be 'file not exists'")

	// Alternatively, check using GetFileContent:
	content, exists := ctx.GetFileContent("deletable_file.txt")
	require.Equal(t, "The file does not exist or could not be read.", content, "content should be empty for non-existent file")
	require.False(t, exists, "file should not exist")
}

func TestLocalAgentContext_DeleteFile_NonExistentFile(t *testing.T) {
	// Create a temporary directory.
	tempDir, err := os.MkdirTemp("", "test-delete-file")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Instantiate a LocalProgrammingAgentContext.
	ctx, _ := NewLocalProgrammingAgentContext(tempDir, "change request", &utils.NoOpGitUtil{})

	// Call ctx.DeleteFile("non_existent_file.txt").
	err = ctx.Delete("non_existent_file.txt")
	require.NoError(t, err)

	// Assert that no error is returned even if the file doesn't exist.
	// We just check that no error is returned from DeleteFile.
}

func TestLocalAgentContext_MoveFile_NonExistentFile(t *testing.T) {
	// Create a temporary directory.
	tempDir, err := os.MkdirTemp("", "test-move-file")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Instantiate a LocalProgrammingAgentContext with the temporary directory.
	ctx, _ := NewLocalProgrammingAgentContext(tempDir, "change request", &utils.NoOpGitUtil{})

	// Call ctx.MoveFile("non_existent_file.txt", "renamed_testfile.txt") to simulate renaming a non-existent file.
	err = ctx.MoveFile("non_existent_file.txt", "renamed_testfile.txt")
	require.Error(t, err)
	require.Equal(t, "file non_existent_file.txt does not exist", err.Error(), "error message should indicate non-existent file")
}

func TestLocalAgentContext_MoveFile_ToNonExistentDirectory(t *testing.T) {
	// Create a temporary directory.
	tempDir, err := os.MkdirTemp("", "test-move-file")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Instantiate a LocalProgrammingAgentContext with the temporary directory.
	ctx, _ := NewLocalProgrammingAgentContext(tempDir, "change request", &utils.NoOpGitUtil{})

	// Create a test file (testfile.txt) within the temporary directory and write initial content to it.
	initialContent := "This is the initial content of the test file."
	testFilePath := filepath.Join(tempDir, "testfile.txt")
	err = os.WriteFile(testFilePath, []byte(initialContent), 0644)
	require.NoError(t, err)

	// Call ctx.MoveFile("testfile.txt", "non_existent_dir/renamed_testfile.txt") to simulate renaming the file to a non-existent directory.
	err = ctx.MoveFile("testfile.txt", "non_existent_dir/renamed_testfile.txt")
	require.NoError(t, err)

	// Verify that ctx.GetFileContent("non_existent_dir/renamed_testfile.txt") returns the expected content and that the file does not exist.
	newContent, exists := ctx.GetFileContent("non_existent_dir/renamed_testfile.txt")
	require.True(t, exists)
	require.Equal(t, initialContent, newContent)

	newFilePath := filepath.Join(tempDir, "non_existent_dir", "renamed_testfile.txt")
	_, err = os.Stat(newFilePath)
	require.Error(t, err, "file should not exist at new path before flush")

	// Flush changes
	err = ctx.FlushChanges()
	require.NoError(t, err)

	// Verify that ctx.GetFileContent("non_existent_dir/renamed_testfile.txt") returns the expected content and that the file exists.
	newContent, exists = ctx.GetFileContent("non_existent_dir/renamed_testfile.txt")
	require.True(t, exists)
	require.Equal(t, initialContent, newContent)

	_, err = os.Stat(newFilePath)
	require.NoError(t, err, "file should exist at new path after flush")
}

func TestLocalAgentContext_DeleteFile_BeingMoved(t *testing.T) {
	// Create a temporary directory.
	tempDir, err := os.MkdirTemp("", "test-delete-file")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create a file named "deletable_file.txt" within the temporary directory.
	deletableFilePath := filepath.Join(tempDir, "deletable_file.txt")
	_, err = os.Create(deletableFilePath)
	require.NoError(t, err)

	// Instantiate a LocalProgrammingAgentContext.
	ctx, _ := NewLocalProgrammingAgentContext(tempDir, "change request", &utils.NoOpGitUtil{})

	// Call ctx.MoveFile("deletable_file.txt", "renamed_testfile.txt") to simulate renaming the file.
	err = ctx.MoveFile("deletable_file.txt", "renamed_testfile.txt")
	require.NoError(t, err)

	// Call ctx.Delete("deletable_file.txt") to simulate deleting the file.
	err = ctx.Delete("deletable_file.txt")
	require.NoError(t, err)

	// Verify that ctx.GetFileContent("renamed_testfile.txt") returns the expected content and that the file does not exist.
	newContent, exists := ctx.GetFileContent("renamed_testfile.txt")
	require.False(t, exists)
	require.Equal(t, "The file has been deleted.", newContent)

	newFilePath := filepath.Join(tempDir, "renamed_testfile.txt")
	_, err = os.Stat(newFilePath)
	require.Error(t, err, "file should not exist at new path before flush")

	// Flush changes
	err = ctx.FlushChanges()
	require.NoError(t, err)

	// Verify that ctx.GetFileContent("renamed_testfile.txt") returns the expected content and that the file does not exist.
	newContent, exists = ctx.GetFileContent("renamed_testfile.txt")
	require.False(t, exists)
	require.Equal(t, "The file does not exist or could not be read.", newContent)

	_, err = os.Stat(newFilePath)
	require.Error(t, err, "file should not exist at new path after flush")
}

func TestLocalAgentContext_MoveAndDeleteMultipleFiles(t *testing.T) {
	// Create a temporary directory.
	tempDir, err := os.MkdirTemp("", "test-move-delete-files")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Instantiate a LocalProgrammingAgentContext with the temporary directory.
	ctx, _ := NewLocalProgrammingAgentContext(tempDir, "change request", &utils.NoOpGitUtil{})

	// Create test files (file1.txt, file2.txt) within the temporary directory and write initial content to them.
	initialContent1 := "This is the initial content of file1."
	testFilePath1 := filepath.Join(tempDir, "file1.txt")
	err = os.WriteFile(testFilePath1, []byte(initialContent1), 0644)
	require.NoError(t, err)

	initialContent2 := "This is the initial content of file2."
	testFilePath2 := filepath.Join(tempDir, "file2.txt")
	err = os.WriteFile(testFilePath2, []byte(initialContent2), 0644)
	require.NoError(t, err)

	// Call ctx.MoveFile("file1.txt", "renamed_file1.txt") to simulate renaming the file.
	err = ctx.MoveFile("file1.txt", "renamed_file1.txt")
	require.NoError(t, err)

	// Call ctx.Delete("file2.txt") to simulate deleting the file.
	err = ctx.Delete("file2.txt")
	require.NoError(t, err)

	// Verify that ctx.GetFileContent("renamed_file1.txt") returns the expected content and that the file does not exist.
	newContent1, exists1 := ctx.GetFileContent("renamed_file1.txt")
	require.Equal(t, true, exists1)
	require.Equal(t, initialContent1, newContent1)

	newFilePath1 := filepath.Join(tempDir, "renamed_file1.txt")
	_, err = os.Stat(newFilePath1)
	require.Error(t, err, "file should not exist at new path before flush")

	// Verify that ctx.GetFileContent("file2.txt") does not exist.
	newContent2, exists2 := ctx.GetFileContent("file2.txt")
	require.False(t, exists2)
	require.Equal(t, "The file has been deleted.", newContent2)

	newFilePath2 := filepath.Join(tempDir, "file2.txt")
	_, err = os.Stat(newFilePath2)
	require.NoError(t, err, "file should exist at old path before flush")

	// Flush changes
	err = ctx.FlushChanges()
	require.NoError(t, err)

	// Verify that ctx.GetFileContent("renamed_file1.txt") returns the expected content and that the file exists.
	newContent1, exists1 = ctx.GetFileContent("renamed_file1.txt")
	require.Equal(t, true, exists1)
	require.Equal(t, initialContent1, newContent1)

	_, err = os.Stat(newFilePath1)
	require.NoError(t, err, "file should exist at new path after flush")

	// Verify that ctx.GetFileContent("file2.txt") returns the expected content and that the file does not exist.
	newContent2, exists2 = ctx.GetFileContent("file2.txt")
	require.Equal(t, false, exists2)
	require.Equal(t, "The file does not exist or could not be read.", newContent2)

	_, err = os.Stat(newFilePath2)
	require.Error(t, err, "file should not exist at old path after flush")
}

func TestLocalAgentContext_DeleteFile_ExistingFileBeforeFlush(t *testing.T) {
	// Create a temporary directory.
	tempDir, err := os.MkdirTemp("", "test-delete-file")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create a file named "deletable_file.txt" within the temporary directory.
	deletableFilePath := filepath.Join(tempDir, "deletable_file.txt")
	_, err = os.Create(deletableFilePath)
	require.NoError(t, err)

	// Instantiate a LocalProgrammingAgentContext.
	ctx, _ := NewLocalProgrammingAgentContext(tempDir, "change request", &utils.NoOpGitUtil{})

	// Call ctx.DeleteFile("deletable_file.txt").
	err = ctx.Delete("deletable_file.txt")
	require.NoError(t, err)

	// Assert that no error is returned and the file "deletable_file.txt" still exists before flush.
	_, err = os.Stat(deletableFilePath)
	require.NoError(t, err, "file should exist before flush")

	// Flush changes
	err = ctx.FlushChanges()
	require.NoError(t, err)

	// Assert that no error is returned and the file "deletable_file.txt" no longer exists.
	_, err = os.Stat(deletableFilePath)
	require.Error(t, err, "file should not exist")
	require.True(t, os.IsNotExist(err), "error should be 'file not exists'")

	// Alternatively, check using GetFileContent:
	content, exists := ctx.GetFileContent("deletable_file.txt")
	require.Equal(t, "The file does not exist or could not be read.", content, "content should be empty for non-existent file")
	require.False(t, exists, "file should not exist")
}

func TestLocalAgentContext_Delete_InMemory(t *testing.T) {
	// Create a temporary directory.
	tempDir, err := os.MkdirTemp("", "test-delete-inmemory")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create a test file.
	testFilePath := filepath.Join(tempDir, "testfile.txt")
	_, err = os.Create(testFilePath)
	require.NoError(t, err)

	// Instantiate a LocalProgrammingAgentContext.
	ctx, _ := NewLocalProgrammingAgentContext(tempDir, "change request", &utils.NoOpGitUtil{})

	// Call Delete on the test file.
	err = ctx.Delete("testfile.txt")
	require.NoError(t, err)

	// Assert that the file still exists on the file system.
	_, err = os.Stat(testFilePath)
	require.NoError(t, err, "file should still exist on file system")

	// Assert that the file path is present in the deletedFiles list.
	require.Contains(t, ctx.deletedFiles, "testfile.txt", "deletedFiles should contain the file path")
}

func TestLocalAgentContext_MoveFile_InMemory(t *testing.T) {
	// Create a temporary directory.
	tempDir, err := os.MkdirTemp("", "test-move-inmemory")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create a test file.
	testFilePath := filepath.Join(tempDir, "testfile.txt")
	_, err = os.Create(testFilePath)
	require.NoError(t, err)

	// Instantiate a LocalProgrammingAgentContext.
	ctx, _ := NewLocalProgrammingAgentContext(tempDir, "change request", &utils.NoOpGitUtil{})

	// Call MoveFile on the test file.
	newPath := "renamed_testfile.txt"
	err = ctx.MoveFile("testfile.txt", newPath)
	require.NoError(t, err)

	// Assert that the file still exists at the old path on the file system.
	_, err = os.Stat(testFilePath)
	require.NoError(t, err, "file should still exist at old path on file system")

	// Assert that an entry exists in the movedFiles map.
	require.Contains(t, ctx.movedFiles, "testfile.txt", "movedFiles should contain the old path")
	require.Equal(t, newPath, ctx.movedFiles["testfile.txt"], "movedFiles should map old path to new path")

	// Assert that CurrentRepoStructure has been updated.
	foundOldPath := false
	foundNewPath := false
	for _, path := range ctx.CurrentRepoStructure {
		if path == "testfile.txt" {
			foundOldPath = true
		}
		if path == newPath {
			foundNewPath = true
		}
	}
	require.False(t, foundOldPath, "CurrentRepoStructure should not contain old path")
	require.True(t, foundNewPath, "CurrentRepoStructure should contain new path")
}

func TestLocalAgentContext_GetRepoStructure_UpdatedInMemory(t *testing.T) {
	// Create a temporary directory.
	tempDir, err := os.MkdirTemp("", "test-repo-structure-inmemory")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create a test file.
	testFilePath := filepath.Join(tempDir, "testfile.txt")
	_, err = os.Create(testFilePath)
	require.NoError(t, err)
	testFilePath2 := filepath.Join(tempDir, "testfile2.txt")
	_, err = os.Create(testFilePath2)
	require.NoError(t, err)

	// Instantiate a LocalProgrammingAgentContext.
	ctx, _ := NewLocalProgrammingAgentContext(tempDir, "change request", &utils.NoOpGitUtil{})

	// Initial repo structure check
	initialStructure := ctx.GetRepoStructure()
	require.Contains(t, initialStructure, "testfile.txt")
	require.Contains(t, initialStructure, "testfile2.txt")

	// Call Delete on testfile.txt
	err = ctx.Delete("testfile.txt")
	require.NoError(t, err)

	// GetRepoStructure after Delete
	structureAfterDelete := ctx.GetRepoStructure()
	require.NotContains(t, structureAfterDelete, "testfile.txt", "GetRepoStructure should not contain deleted file")
	require.Contains(t, structureAfterDelete, "testfile2.txt")

	// Call MoveFile on testfile2.txt
	newPath := "renamed_testfile2.txt"
	err = ctx.MoveFile("testfile2.txt", newPath)
	require.NoError(t, err)

	// GetRepoStructure after MoveFile
	structureAfterMove := ctx.GetRepoStructure()
	require.NotContains(t, structureAfterMove, "testfile2.txt", "GetRepoStructure should not contain old path after move")
	require.Contains(t, structureAfterMove, newPath, "GetRepoStructure should contain new path after move")

	// Ensure the structure is sorted (as GetRepoStructure might return sorted list) for consistent comparison in some cases
	sort.Strings(structureAfterMove)
	expectedStructure := []string{"renamed_testfile2.txt"}
	sort.Strings(expectedStructure)

}
