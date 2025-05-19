package commands

import (
	"encoding/json"
	"fmt"
	"github.com/EduardDranca/GoAgent/internal/agent/context"
	"github.com/EduardDranca/GoAgent/internal/logging"
	"github.com/EduardDranca/GoAgent/internal/utils"
	"github.com/fatih/color"
	"strings"
)

// Command interface
type Command interface {
	Process(agentContext context.ProgrammingAgentContext) (string, error)
}

// ReadCommand struct represents a command to read file contents.
type ReadCommand struct {
	Files []string
}

// Process for ReadCommand retrieves the content of specified files.
func (c *ReadCommand) Process(agentContext context.ProgrammingAgentContext) (string, error) {
	logging.Logger.Infof("Executing command: Read files: %v", c.Files)
	var sb strings.Builder
	for _, file := range c.Files {
		content, _ := agentContext.GetFileContent(file)
		sb.WriteString(formatFileContent(file, content))
	}
	return sb.String(), nil
}

// CheckStructureCommand struct represents a command to check the repository structure.
type CheckStructureCommand struct{}

// Process for CheckStructureCommand retrieves the repository structure.
func (c *CheckStructureCommand) Process(agentContext context.ProgrammingAgentContext) (string, error) {
	logging.Logger.Infof("Executing command: Check repository structure")
	structure := agentContext.GetRepoStructure()
	structureJSON, err := json.MarshalIndent(structure, "", "  ")
	if err != nil {
		return "", fmt.Errorf("error converting structure to JSON: %w", err)
	}
	return fmt.Sprintf("The current project structure is as follows:\n%s", string(structureJSON)), nil
}

// SearchCommand struct represents a command to search for code.
type SearchCommand struct {
	Query string
}

// Process for SearchCommand searches for files containing the specified query.
func (c *SearchCommand) Process(agentContext context.ProgrammingAgentContext) (string, error) {
	logging.Logger.Infof("Executing command: Search code for: %s", c.Query)
	searchResults := agentContext.SearchCode(c.Query)
	filesJSON, err := json.MarshalIndent(searchResults, "", "  ")
	if err != nil {
		return "", fmt.Errorf("error converting files to JSON: %w", err)
	}
	return fmt.Sprintf("Files containing %s:\n%s", c.Query, string(filesJSON)), nil
}

// UpdateFileCommand struct represents a command to update a file with an implementation plan.
type UpdateFileCommand struct {
	FilePath           string   `json:"file_path"`
	ImplementationPlan string   `json:"implementation_plan"`
	ContextFiles       []string `json:"context_files"`
}

// Process for UpdateFileCommand returns a message indicating file update.
func (c *UpdateFileCommand) Process(_ context.ProgrammingAgentContext) (string, error) {
	logging.Logger.Infof(color.GreenString("%s:", c.FilePath))
	out, err := utils.RenderWithGlamour(c.ImplementationPlan)
	if err != nil {
		logging.Logger.Infof(color.HiWhiteString(c.ImplementationPlan))
	} else {
		logging.Logger.Infof(out) // Print glamour-rendered output to stdout
	}
	return fmt.Sprintf("The file %s was updated, please carry on with the change request.", c.FilePath), nil

}

// MoveFileCommand struct represents a command to move a file.
type MoveFileCommand struct {
	OldPath string `json:"old_path"`
	NewPath string `json:"new_path"`
}

// Process for MoveFileCommand moves the specified file.
func (c *MoveFileCommand) Process(agentContext context.ProgrammingAgentContext) (string, error) {
	logging.Logger.Infof("Executing command: Move file from %s to %s", c.OldPath, c.NewPath)
	err := agentContext.MoveFile(c.OldPath, c.NewPath)
	if err != nil {
		commandError := fmt.Errorf("error moving file: %w", err).Error()
		return commandError, nil
	}
	return fmt.Sprintf("File moved from %s to %s", c.OldPath, c.NewPath), nil
}

// DeleteFileCommand struct represents a command to delete a file.
type DeleteFileCommand struct {
	FilePath string `json:"file_path"`
}

// Process for DeleteFileCommand deletes the specified file.
func (c *DeleteFileCommand) Process(agentContext context.ProgrammingAgentContext) (string, error) {
	logging.Logger.Infof("Executing command: Delete file %s", c.FilePath)
	err := agentContext.Delete(c.FilePath)
	if err != nil {
		return "", fmt.Errorf("error deleting file: %w", err)
	}
	return fmt.Sprintf("File deleted: %s", c.FilePath), nil
}

// CommitCommand struct represents a command to commit changes.
type CommitCommand struct {
	Message string `json:"commit"`
}

// Process for CommitCommand logs the commit message and returns it.
func (c *CommitCommand) Process(_ context.ProgrammingAgentContext) (string, error) {
	logging.Logger.Infof("Executing command: Commit with message: %s", c.Message)
	return c.Message, nil
}

// RespondCommand struct represents a command to respond to a prompt.
type RespondCommand struct {
	Message string `json:"answer"`
}

// Process for RespondCommand logs the response message and returns it.
func (c *RespondCommand) Process(_ context.ProgrammingAgentContext) (string, error) {
	logging.Logger.Infof("Executing command: Respond with message: %s", c.Message)
	return c.Message, nil
}

// formatFileContent formats the content of a file for output.
func formatFileContent(file string, content string) string {
	return fmt.Sprintf("Content of %s:\n%s\n\n", file, content)
}

// NewCommand function constructs a Command from command string and parameters map.
func NewCommand(commandMap map[string]interface{}) (Command, error) {
	switch commandMap["command"] {
	case "read":
		filesRaw, ok := commandMap["files"]
		if !ok {
			return nil, fmt.Errorf("missing 'files' parameter for read command")
		}
		files, ok := filesRaw.([]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid 'files' parameter type for read command")
		}
		var filePaths []string
		for _, f := range files {
			filePath, ok := f.(string)
			if !ok {
				return nil, fmt.Errorf("invalid file path type in 'files' parameter for read command")
			}
			filePaths = append(filePaths, filePath)
		}
		return &ReadCommand{Files: filePaths}, nil

	case "check_structure":
		return &CheckStructureCommand{}, nil

	case "search":
		queryRaw, ok := commandMap["query"]
		if !ok {
			return nil, fmt.Errorf("missing 'query' parameter for search command")
		}
		query, ok := queryRaw.(string)
		if !ok {
			return nil, fmt.Errorf("invalid 'query' parameter type for search command")
		}
		return &SearchCommand{Query: query}, nil

	case "update_file":
		filePathRaw, ok := commandMap["file_path"]
		if !ok {
			return nil, fmt.Errorf("missing 'file_path' parameter for update_file command")
		}
		filePath, ok := filePathRaw.(string)
		if !ok {
			return nil, fmt.Errorf("invalid 'file_path' parameter type for update_file command")
		}

		implementationPlanRaw, ok := commandMap["implementation_plan"]
		if !ok {
			return nil, fmt.Errorf("missing 'implementation_plan' parameter for update_file command")
		}
		implementationPlan, ok := implementationPlanRaw.(string)
		if !ok {
			return nil, fmt.Errorf("invalid 'implementation_plan' parameter type for update_file command")
		}

		var contextFiles []string
		contextFilesRaw, ok := commandMap["context_files"]
		if !ok {
			contextFiles = []string{}
		} else {
			contextFilesSlice, ok := contextFilesRaw.([]interface{})
			if !ok {
				return nil, fmt.Errorf("invalid 'context_files' parameter type for update_file command")
			}
			contextFiles = convertToStringArray(contextFilesSlice)
		}

		return &UpdateFileCommand{
			FilePath:           filePath,
			ImplementationPlan: implementationPlan,
			ContextFiles:       contextFiles,
		}, nil

	case "move_file":
		oldPathRaw, ok := commandMap["old_path"]
		if !ok {
			return nil, fmt.Errorf("missing 'old_path' parameter for move_file command")
		}
		oldPath, ok := oldPathRaw.(string)
		if !ok {
			return nil, fmt.Errorf("invalid 'old_path' parameter type for move_file command")
		}

		newPathRaw, ok := commandMap["new_path"]
		if !ok {
			return nil, fmt.Errorf("missing 'new_path' parameter for move_file command")
		}
		newPath, ok := newPathRaw.(string)
		if !ok {
			return nil, fmt.Errorf("invalid 'new_path' parameter type for move_file command")
		}
		return &MoveFileCommand{OldPath: oldPath, NewPath: newPath}, nil

	case "delete_file":
		filePathRaw, ok := commandMap["file_path"]
		if !ok {
			return nil, fmt.Errorf("missing 'file_path' parameter for delete_file command")
		}
		filePath, ok := filePathRaw.(string)
		if !ok {
			return nil, fmt.Errorf("invalid 'file_path' parameter type for delete_file command")
		}
		return &DeleteFileCommand{FilePath: filePath}, nil

	case "commit":
		messageRaw, ok := commandMap["message"]
		if !ok {
			return nil, fmt.Errorf("missing 'message' parameter for commit command")
		}
		message, ok := messageRaw.(string)
		if !ok {
			return nil, fmt.Errorf("invalid 'message' parameter type for commit command")
		}
		return &CommitCommand{Message: message}, nil

	case "respond":
		messageRaw, ok := commandMap["answer"]
		if !ok {
			return nil, fmt.Errorf("missing 'message' parameter for respond command")
		}
		message, ok := messageRaw.(string)
		if !ok {
			return nil, fmt.Errorf("invalid 'message' parameter type for respond command")
		}
		return &RespondCommand{Message: message}, nil

	default:
		return nil, fmt.Errorf("unknown command: %s", commandMap["command"])
	}
}

// convertToStringArray converts an interface{} to a []string, handling type assertions and errors.
func convertToStringArray(input interface{}) []string {
	if input == nil {
		return nil
	}
	slice, ok := input.([]interface{})
	if !ok {
		logging.Logger.Errorf("Error: 'files' field is not an array")
		return nil
	}

	stringSlice := make([]string, 0, len(slice))
	for _, v := range slice {
		str, ok := v.(string)
		if !ok {
			logging.Logger.Errorf("Error: element in 'files' array is not a string")
			return nil // Or handle the error as appropriate for your application
		}
		stringSlice = append(stringSlice, str)
	}
	return stringSlice
}
