package commands

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/EduardDranca/GoAgent/internal/agent/context"
	"reflect"
	"testing"

	"go.uber.org/mock/gomock"
)

func TestReadCommand_Process(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockContext := context.NewMockProgrammingAgentContext(ctrl)
	mockContext.EXPECT().GetFileContent("file1.txt").Return("This is file1 content.", true)
	mockContext.EXPECT().GetFileContent("file2.txt").Return("This is file2 content.", true)

	command := &ReadCommand{Files: []string{"file1.txt", "file2.txt"}}

	output, err := command.Process(mockContext)

	if err != nil {
		t.Fatalf("ReadCommand.Process failed: %v", err)
	}

	expectedOutput := "Content of file1.txt:\nThis is file1 content.\n\nContent of file2.txt:\nThis is file2 content.\n\n"

	if output != expectedOutput {
		t.Errorf("ReadCommand.Process: Unexpected output:\nGot:  %q\nWant: %q", output, expectedOutput)
	}
}

func TestReadCommand_Process_EmptyFiles(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockContext := context.NewMockProgrammingAgentContext(ctrl)
	mockContext.EXPECT().GetFileContent("empty_file1.txt").Return("", true) // Simulate empty file 1
	mockContext.EXPECT().GetFileContent("empty_file2.txt").Return("", true) // Simulate empty file 2

	command := &ReadCommand{Files: []string{"empty_file1.txt", "empty_file2.txt"}}

	output, err := command.Process(mockContext)

	if err != nil {
		t.Fatalf("ReadCommand_Process_EmptyFiles failed: %v", err)
	}

	expectedOutput := "Content of empty_file1.txt:\n\n\nContent of empty_file2.txt:\n\n\n"

	if output != expectedOutput {
		t.Errorf("ReadCommand_Process_EmptyFiles: Unexpected output:\nGot:  %q\nWant: %q", output, expectedOutput)
	}
}

func TestReadCommand_Process_NonExistentFile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockContext := context.NewMockProgrammingAgentContext(ctrl)
	mockContext.EXPECT().GetFileContent("non_existent_file.txt").Return("", false) // Simulate file not found

	command := &ReadCommand{Files: []string{"non_existent_file.txt"}}

	output, err := command.Process(mockContext)

	if err != nil {
		t.Fatalf("ReadCommand_Process_NonExistentFile failed: %v", err)
	}

	expectedOutput := "Content of non_existent_file.txt:\n\n\n"

	if output != expectedOutput {
		t.Errorf("ReadCommand_Process_NonExistentFile: Unexpected output:\nGot:  %q\nWant: %q", output, expectedOutput)
	}
}

func TestReadCommand_Process_SpecialCharacters(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockContext := context.NewMockProgrammingAgentContext(ctrl)
	specialCharsContent := "This file contains special characters: !@#$%^&*()_+=-`~[]\\{}|;':\",./<>?"
	mockContext.EXPECT().GetFileContent("special_chars_file.txt").Return(specialCharsContent, true)

	command := &ReadCommand{Files: []string{"special_chars_file.txt"}}

	output, err := command.Process(mockContext)

	if err != nil {
		t.Fatalf("ReadCommand_Process_SpecialCharacters failed: %v", err)
	}

	expectedOutput := "Content of special_chars_file.txt:\nThis file contains special characters: !@#$%^&*()_+=-`~[]\\{}|;':\",./<>?\n\n"

	if output != expectedOutput {
		t.Errorf("ReadCommand_Process_SpecialCharacters: Unexpected output:\nGot:  %q\nWant: %q", output, expectedOutput)
	}
}

func TestCheckStructureCommand_Process(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockContext := context.NewMockProgrammingAgentContext(ctrl)
	mockContext.EXPECT().GetRepoStructure().Return([]string{"file1.txt", "dir1/", "dir1/file2.txt"})

	command := &CheckStructureCommand{}
	output, err := command.Process(mockContext)

	if err != nil {
		t.Fatalf("CheckStructureCommand.Process failed: %v", err)
	}

	expectedStructure := []string{"file1.txt", "dir1/", "dir1/file2.txt"}
	expectedStructureJSON, _ := json.MarshalIndent(expectedStructure, "", "  ")

	expectedOutput := "The current project structure is as follows:\n" + string(expectedStructureJSON)

	if output != expectedOutput {
		t.Errorf("CheckStructureCommand.Process: Unexpected output:\nGot:  %q\nWant: %q", output, expectedOutput)
	}
}

func TestSearchCommand_Process(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockContext := context.NewMockProgrammingAgentContext(ctrl)
	mockContext.EXPECT().SearchCode("test").Return(map[string][]int{
		"file1.txt": {1, 5, 10},
		"file2.txt": {2, 7},
	})

	command := &SearchCommand{Query: "test"}
	output, err := command.Process(mockContext)

	if err != nil {
		t.Fatalf("SearchCommand.Process failed: %v", err)
	}

	expectedResultsJSON, _ := json.MarshalIndent(map[string][]int{
		"file1.txt": {1, 5, 10},
		"file2.txt": {2, 7},
	}, "", "  ")
	expectedOutput := "Files containing test:\n" + string(expectedResultsJSON)

	if !reflect.DeepEqual(output, expectedOutput) {
		t.Errorf("SearchCommand.Process: Unexpected output:\nGot:  %q\nWant: %q", output, expectedOutput)
	}
}

// UpdateFileCommand does not use context
func TestUpdateFileCommand_Process(t *testing.T) {
	command := &UpdateFileCommand{
		FilePath:           "file1.txt",
		ImplementationPlan: "Update file1.txt with new content.",
	}
	output, err := command.Process(nil) // context is not used in UpdateFileCommand

	if err != nil {
		t.Fatalf("UpdateFileCommand.Process failed: %v", err)
	}

	expectedOutput := "The file file1.txt was updated, please carry on with the change request."

	if output != expectedOutput {
		t.Errorf("UpdateFileCommand.Process: Unexpected output:\nGot:  %q\nWant: %q", output, expectedOutput)
	}
}

func TestMoveFileCommand_Process(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockContext := context.NewMockProgrammingAgentContext(ctrl)
	oldPath := "old/path/file.txt"
	newPath := "new/path/file.txt"

	mockContext.EXPECT().MoveFile(oldPath, newPath).Return(nil).Times(1)

	command := &MoveFileCommand{OldPath: oldPath, NewPath: newPath}
	output, err := command.Process(mockContext)

	if err != nil {
		t.Fatalf("MoveFileCommand.Process failed: %v", err)
	}

	expectedOutput := fmt.Sprintf("File moved from %s to %s", oldPath, newPath)

	if output != expectedOutput {
		t.Errorf("MoveFileCommand.Process: Unexpected output:\nGot:  %q\nWant: %q", output, expectedOutput)
	}
}

func TestMoveFileCommand_Process_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockContext := context.NewMockProgrammingAgentContext(ctrl)
	oldPath := "old/path/file.txt"
	newPath := "new/path/file.txt"
	moveErr := errors.New("simulated move error")

	mockContext.EXPECT().MoveFile(oldPath, newPath).Return(moveErr).Times(1)

	command := &MoveFileCommand{OldPath: oldPath, NewPath: newPath}
	output, err := command.Process(mockContext)

	if err != nil {
		t.Fatalf("MoveFileCommand.Process failed: %v", err)
	}

	// The command is designed to return the error message as the output string
	expectedOutput := fmt.Errorf("error moving file: %w", moveErr).Error()

	if output != expectedOutput {
		t.Errorf("MoveFileCommand.Process: Unexpected output:\nGot:  %q\nWant: %q", output, expectedOutput)
	}
}

func TestDeleteFileCommand_Process(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockContext := context.NewMockProgrammingAgentContext(ctrl)
	filePath := "path/to/delete/file.txt"

	mockContext.EXPECT().Delete(filePath).Return(nil).Times(1)

	command := &DeleteFileCommand{FilePath: filePath}
	output, err := command.Process(mockContext)

	if err != nil {
		t.Fatalf("DeleteFileCommand.Process failed: %v", err)
	}

	expectedOutput := fmt.Sprintf("File deleted: %s", filePath)

	if output != expectedOutput {
		t.Errorf("DeleteFileCommand.Process: Unexpected output:\nGot:  %q\nWant: %q", output, expectedOutput)
	}
}

func TestDeleteFileCommand_Process_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockContext := context.NewMockProgrammingAgentContext(ctrl)
	filePath := "path/to/delete/file.txt"
	deleteErr := errors.New("simulated delete error")

	mockContext.EXPECT().Delete(filePath).Return(deleteErr).Times(1)

	command := &DeleteFileCommand{FilePath: filePath}
	output, err := command.Process(mockContext)

	if err == nil {
		t.Fatalf("DeleteFileCommand.Process did not return an error")
	}

	expectedError := fmt.Errorf("error deleting file: %w", deleteErr)

	if err.Error() != expectedError.Error() {
		t.Errorf("DeleteFileCommand.Process: Unexpected error:\nGot:  %v\nWant: %v", err, expectedError)
	}

	// The command returns an empty string on error
	if output != "" {
		t.Errorf("DeleteFileCommand.Process: Unexpected output on error:\nGot:  %q\nWant: %q", output, "")
	}
}

func TestNewCommand(t *testing.T) {
	tests := []struct {
		name            string
		commandMap      map[string]interface{}
		expectedCommand Command
		expectedError   error
	}{
		{
			name: "ReadCommand with valid parameters",
			commandMap: map[string]interface{}{
				"command": "read",
				"files":   []interface{}{"file1.txt", "file2.txt"},
			},
			expectedCommand: &ReadCommand{Files: []string{"file1.txt", "file2.txt"}},
			expectedError:   nil,
		},
		{
			name: "ReadCommand with missing files parameter",
			commandMap: map[string]interface{}{
				"command": "read",
			},
			expectedCommand: nil,
			expectedError:   errors.New("missing 'files' parameter for read command"),
		},
		{
			name: "ReadCommand with invalid files parameter type",
			commandMap: map[string]interface{}{
				"command": "read",
				"files":   "file1.txt",
			},
			expectedCommand: nil,
			expectedError:   errors.New("invalid 'files' parameter type for read command"),
		},
		{
			name: "ReadCommand with invalid file path type",
			commandMap: map[string]interface{}{
				"command": "read",
				"files":   []interface{}{"file1.txt", 123},
			},
			expectedCommand: nil,
			expectedError:   errors.New("invalid file path type in 'files' parameter for read command"),
		},
		{
			name: "CheckStructureCommand",
			commandMap: map[string]interface{}{
				"command": "check_structure",
			},
			expectedCommand: &CheckStructureCommand{},
			expectedError:   nil,
		},
		{
			name: "SearchCommand with valid parameters",
			commandMap: map[string]interface{}{
				"command": "search",
				"query":   "test",
			},
			expectedCommand: &SearchCommand{Query: "test"},
			expectedError:   nil,
		},
		{
			name: "SearchCommand with missing query parameter",
			commandMap: map[string]interface{}{
				"command": "search",
			},
			expectedCommand: nil,
			expectedError:   errors.New("missing 'query' parameter for search command"),
		},
		{
			name: "SearchCommand with invalid query parameter type",
			commandMap: map[string]interface{}{
				"command": "search",
				"query":   123,
			},
			expectedCommand: nil,
			expectedError:   errors.New("invalid 'query' parameter type for search command"),
		},
		{
			name: "UpdateFileCommand with valid parameters",
			commandMap: map[string]interface{}{
				"command":             "update_file",
				"file_path":           "file1.txt",
				"implementation_plan": "Update file1.txt with new content.",
				"context_files":       []interface{}{"file1.txt", "file2.txt"},
			},
			expectedCommand: &UpdateFileCommand{
				FilePath:           "file1.txt",
				ImplementationPlan: "Update file1.txt with new content.",
				ContextFiles:       []string{"file1.txt", "file2.txt"},
			},
			expectedError: nil,
		},
		{
			name: "UpdateFileCommand with missing file_path parameter",
			commandMap: map[string]interface{}{
				"command":             "update_file",
				"implementation_plan": "Update file1.txt with new content.",
				"context_files":       []interface{}{"file1.txt", "file2.txt"},
			},
			expectedCommand: nil,
			expectedError:   errors.New("missing 'file_path' parameter for update_file command"),
		},
		{
			name: "UpdateFileCommand with invalid file_path parameter type",
			commandMap: map[string]interface{}{
				"command":             "update_file",
				"file_path":           123,
				"implementation_plan": "Update file1.txt with new content.",
				"context_files":       []interface{}{"file1.txt", "file2.txt"},
			},
			expectedCommand: nil,
			expectedError:   errors.New("invalid 'file_path' parameter type for update_file command"),
		},
		{
			name: "UpdateFileCommand with missing implementation_plan parameter",
			commandMap: map[string]interface{}{
				"command":       "update_file",
				"file_path":     "file1.txt",
				"context_files": []interface{}{"file1.txt", "file2.txt"},
			},
			expectedCommand: nil,
			expectedError:   errors.New("missing 'implementation_plan' parameter for update_file command"),
		},
		{
			name: "UpdateFileCommand with invalid implementation_plan parameter type",
			commandMap: map[string]interface{}{
				"command":             "update_file",
				"file_path":           "file1.txt",
				"implementation_plan": 123,
				"context_files":       []interface{}{"file1.txt", "file2.txt"},
			},
			expectedCommand: nil,
			expectedError:   errors.New("invalid 'implementation_plan' parameter type for update_file command"),
		},
		{
			name: "UpdateFileCommand with missing context_files parameter",
			commandMap: map[string]interface{}{
				"command":             "update_file",
				"file_path":           "file1.txt",
				"implementation_plan": "Update file1.txt with new content.",
			},
			expectedCommand: &UpdateFileCommand{
				FilePath:           "file1.txt",
				ImplementationPlan: "Update file1.txt with new content.",
				ContextFiles:       []string{},
			},
			expectedError: nil,
		},
		{
			name: "UpdateFileCommand with invalid context_files parameter type",
			commandMap: map[string]interface{}{
				"command":             "update_file",
				"file_path":           "file1.txt",
				"implementation_plan": "Update file1.txt with new content.",
				"context_files":       "file1.txt",
			},
			expectedCommand: nil,
			expectedError:   errors.New("invalid 'context_files' parameter type for update_file command"),
		},
		{
			name: "MoveFileCommand with valid parameters",
			commandMap: map[string]interface{}{
				"command":  "move_file",
				"old_path": "old/file.txt",
				"new_path": "new/file.txt",
			},
			expectedCommand: &MoveFileCommand{OldPath: "old/file.txt", NewPath: "new/file.txt"},
			expectedError:   nil,
		},
		{
			name: "MoveFileCommand with missing old_path parameter",
			commandMap: map[string]interface{}{
				"command":  "move_file",
				"new_path": "new/file.txt",
			},
			expectedCommand: nil,
			expectedError:   errors.New("missing 'old_path' parameter for move_file command"),
		},
		{
			name: "MoveFileCommand with invalid old_path parameter type",
			commandMap: map[string]interface{}{
				"command":  "move_file",
				"old_path": 123,
				"new_path": "new/file.txt",
			},
			expectedCommand: nil,
			expectedError:   errors.New("invalid 'old_path' parameter type for move_file command"),
		},
		{
			name: "MoveFileCommand with missing new_path parameter",
			commandMap: map[string]interface{}{
				"command":  "move_file",
				"old_path": "old/file.txt",
			},
			expectedCommand: nil,
			expectedError:   errors.New("missing 'new_path' parameter for move_file command"),
		},
		{
			name: "MoveFileCommand with invalid new_path parameter type",
			commandMap: map[string]interface{}{
				"command":  "move_file",
				"old_path": "old/file.txt",
				"new_path": 123,
			},
			expectedCommand: nil,
			expectedError:   errors.New("invalid 'new_path' parameter type for move_file command"),
		},
		{
			name: "DeleteFileCommand with valid parameters",
			commandMap: map[string]interface{}{
				"command":   "delete_file",
				"file_path": "file1.txt",
			},
			expectedCommand: &DeleteFileCommand{FilePath: "file1.txt"},
			expectedError:   nil,
		},
		{
			name: "DeleteFileCommand with missing file_path parameter",
			commandMap: map[string]interface{}{
				"command": "delete_file",
			},
			expectedCommand: nil,
			expectedError:   errors.New("missing 'file_path' parameter for delete_file command"),
		},
		{
			name: "DeleteFileCommand with invalid file_path parameter type",
			commandMap: map[string]interface{}{
				"command":   "delete_file",
				"file_path": 123,
			},
			expectedCommand: nil,
			expectedError:   errors.New("invalid 'file_path' parameter type for delete_file command"),
		},
		{
			name: "CommitCommand with valid parameters",
			commandMap: map[string]interface{}{
				"command": "commit",
				"message": "Initial commit",
			},
			expectedCommand: &CommitCommand{Message: "Initial commit"},
			expectedError:   nil,
		},
		{
			name: "CommitCommand with missing message parameter",
			commandMap: map[string]interface{}{
				"command": "commit",
			},
			expectedCommand: nil,
			expectedError:   errors.New("missing 'message' parameter for commit command"),
		},
		{
			name: "CommitCommand with invalid message parameter type",
			commandMap: map[string]interface{}{
				"command": "commit",
				"message": 123,
			},
			expectedCommand: nil,
			expectedError:   errors.New("invalid 'message' parameter type for commit command"),
		},
		{
			name: "RespondCommand with valid parameters",
			commandMap: map[string]interface{}{
				"command": "respond",
				"answer":  "Response to prompt",
			},
			expectedCommand: &RespondCommand{Message: "Response to prompt"},
			expectedError:   nil,
		},
		{
			name: "RespondCommand with missing answer parameter",
			commandMap: map[string]interface{}{
				"command": "respond",
			},
			expectedCommand: nil,
			expectedError:   errors.New("missing 'message' parameter for respond command"),
		},
		{
			name: "RespondCommand with invalid answer parameter type",
			commandMap: map[string]interface{}{
				"command": "respond",
				"answer":  123,
			},
			expectedCommand: nil,
			expectedError:   errors.New("invalid 'message' parameter type for respond command"),
		},
		{
			name: "Unknown command",
			commandMap: map[string]interface{}{
				"command": "unknown",
			},
			expectedCommand: nil,
			expectedError:   errors.New("unknown command: unknown"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			command, err := NewCommand(test.commandMap)

			if test.expectedError != nil {
				if err == nil {
					t.Errorf("NewCommand(%v) = %v, %v; want %v, %v", test.commandMap, command, err, test.expectedCommand, test.expectedError)
				} else if err.Error() != test.expectedError.Error() {
					t.Errorf("NewCommand(%v) = %v, %v; want %v, %v", test.commandMap, command, err, test.expectedCommand, test.expectedError)
				}
			} else {
				if err != nil {
					t.Errorf("NewCommand(%v) = %v, %v; want %v, %v", test.commandMap, command, err, test.expectedCommand, test.expectedError)
				} else if !reflect.DeepEqual(command, test.expectedCommand) {
					t.Errorf("NewCommand(%v) = %v, %v; want %v, %v", test.commandMap, command, err, test.expectedCommand, test.expectedError)
				}
			}
		})
	}
}
