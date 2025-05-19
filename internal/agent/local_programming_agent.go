package agent

import (
	"fmt"
	"github.com/EduardDranca/GoAgent/internal/agent/context"
	"github.com/EduardDranca/GoAgent/internal/agent/models"
	"github.com/EduardDranca/GoAgent/internal/agent/service"
	"github.com/EduardDranca/GoAgent/internal/input"
	"github.com/EduardDranca/GoAgent/internal/logging"
	"github.com/EduardDranca/GoAgent/internal/utils"
	"strings"
)

// LocalProgrammingAgent implements the Agent interface for local file changes.
type LocalProgrammingAgent struct {
	programmingService service.ProgrammingService
	gitUtil            utils.GitUtil // Inject GitUtil interface
	autoCommit         bool
}

// NewLocalProgrammingAgent creates a new LocalProgrammingAgent.
func NewLocalProgrammingAgent(programmingService service.ProgrammingService, gitUtil utils.GitUtil) AgentInterface[models.AgentRequest] {
	if gitUtil == nil {
		gitUtil = &utils.RealGitUtil{} // Default to RealGitUtil if nil is provided
	}
	return &LocalProgrammingAgent{programmingService: programmingService, gitUtil: gitUtil, autoCommit: false}
}

// SetAutoCommit sets the autoCommit field of the LocalProgrammingAgent.
func (a *LocalProgrammingAgent) SetAutoCommit(autoCommit bool) {
	a.autoCommit = autoCommit
}

func (a *LocalProgrammingAgent) resetAndWrapError(dir string, baseError error, message string, isFlushError bool) error {
	if !isFlushError {
		logging.Logger.Warnf("An error occurred during implementation: %s\n. %v\n", message, baseError)
		skipFlushChoice, err := input.UserInputGetter("Would you like to skip flushing changes and proceed without saving the changes to disk? [Y]es/[N]o.")
		if err != nil {
			logging.Logger.Warnf("Error reading skip flush choice, defaulting to skip: %v", err)
			return err
		}

		if strings.ToLower(skipFlushChoice) == "yes" || strings.ToLower(skipFlushChoice) == "y" {
			logging.Logger.Infof("User chose to skip flushing changes.")
			return fmt.Errorf("error occurred during implementation, user skipped flushing changes: %w", baseError)
		}
		logging.Logger.Infof("User chose not to skip flushing changes, proceeding with flush and potential reset on failure.")
		return nil
	}

	if isFlushError {
		// Handle flush error or if user didn't choose to skip flush in implementation error case
		logging.Logger.Warnf("An error occurred during flushing changes to disk: %s\n. %v\n", message, baseError)
		var resetChoice string
		resetChoice, err := input.UserInputGetter("Do you want to reset to the last commit and discard the changes made by the agent? [Y]es/[N]o.")
		if err != nil {
			logging.Logger.Warnf("Error reading reset choice, defaulting to reset: %v", err)
			resetChoice = "yes" // Default to reset if input reading fails for flush error
		}

		if strings.ToLower(resetChoice) == "yes" || strings.ToLower(resetChoice) == "y" {
			resetErr := a.gitUtil.ResetToHead(dir)
			if resetErr != nil {
				return fmt.Errorf("error resetting repository after user confirmed reset: %w, original flush error: %w", resetErr, baseError)
			}
			logging.Logger.Infof("User chose to reset repository after flush error.")
			return fmt.Errorf("error occurred during flushing changes, user chose to reset repository: %w", baseError)
		}
		logging.Logger.Infof("User chose not to reset repository after flush error.")
		return nil
	}

	return fmt.Errorf("error occurred during implementation or flushing changes: %w", baseError)
}

// Implement implements the Agent interface for LocalProgrammingAgent.
func (a *LocalProgrammingAgent) Implement(request models.AgentRequest) error {
	// Skip empty change requests
	if strings.TrimSpace(request.Query) == "" {
		logging.Logger.Infof("Skipping empty change request.")
		return nil
	}

	// Create a programming context
	programmingAgentContext, err := a.createContext(request.Directory, request.Query, a.gitUtil)
	if err != nil {
		// Use errors.New for consistent error wrapping if desired, or fmt.Errorf
		return fmt.Errorf("error initializing programming context: %w", err)
	}

	// Display working message
	logging.Logger.Infof("Working on request...")

	// Implement the plan with context
	commitMessage, err := a.programmingService.ImplementWithContext(programmingAgentContext)
	if err != nil {
		resetErr := a.resetAndWrapError(request.Directory, err, "error implementing plan with context", false)
		if resetErr != nil {
			return resetErr
		}
	}

	// Flush changes to the file system
	err = programmingAgentContext.FlushChanges()
	if err != nil {
		resetErr := a.resetAndWrapError(request.Directory, err, "error updating files", true)
		if resetErr != nil {
			return resetErr
		}
	}

	// Use commit message from ImplementWithContext or default if empty or error
	finalCommitMessage := commitMessage
	if finalCommitMessage == "" {
		finalCommitMessage = "Automated changes by GoAgent" // Default commit message
	}

	err = a.handleCommit(request.Directory, finalCommitMessage)
	if err != nil {
		return err
	}

	return nil
}

// createContext initializes the programming context.
func (a *LocalProgrammingAgent) createContext(dir string, request string, gitUtil utils.GitUtil) (*context.LocalProgrammingAgentContext, error) {
	programmingAgentContext, err := context.NewLocalProgrammingAgentContext(dir, request, gitUtil)
	if err != nil {
		return nil, fmt.Errorf("error initializing programming context: %w", err)
	}
	return programmingAgentContext, nil
}

// Ask implements the Ask method for LocalProgrammingAgent.
func (a *LocalProgrammingAgent) Ask(req models.AgentRequest) (string, error) {
	logging.Logger.Infof("Starting Ask function with request: %s", req.Query)

	// Create local agent context
	agentContext, err := context.NewLocalProgrammingAgentContext(req.Directory, req.Query, a.gitUtil)
	if err != nil {
		logging.Logger.Errorf("Error creating agent context: %v", err)
		return "", fmt.Errorf("error creating agent context: %w", err)
	}

	// Create LLM programming service
	llmService := a.programmingService

	// Call AskWithContext
	answer, err := llmService.AskWithContext(agentContext)
	if err != nil {
		logging.Logger.Errorf("Error in AskWithContext: %v", err)
		return "", fmt.Errorf("error in AskWithContext: %w", err)
	}

	return answer, nil
}

// handleCommit encapsulates all commit related logic including prompting, reading input, and committing changes.
func (a *LocalProgrammingAgent) handleCommit(directory string, commitMessage string) error {
	if a.autoCommit {
		logging.Logger.Infof("Auto-committing changes...")
		return a.executeGitCommit(directory, commitMessage)
	}

	choice := a.promptForCommit()

	switch choice {
	case "Y":
		logging.Logger.Infof("User chose to commit.")
		return a.executeGitCommit(directory, commitMessage)
	case "A":
		logging.Logger.Infof("User chose to always commit (auto-commit enabled).")
		a.autoCommit = true
		return a.executeGitCommit(directory, commitMessage)
	case "E":
		logging.Logger.Infof("User chose to edit commit message.")
		editedMessage, err := utils.EditCommitMessage(commitMessage)
		if err != nil {
			logging.Logger.Errorf("Error editing commit message: %v", err)
			// Ask user if they want to commit with original message or skip
			retryChoice, inputErr := input.UserInputGetter("Error editing message. Commit with original message? [Y]es/[N]o ")
			if inputErr != nil || strings.ToUpper(retryChoice) != "Y" {
				logging.Logger.Infof("User chose not to commit after edit error.")
				// Return specific error about edit failure + user skip
				return fmt.Errorf("failed to edit commit message and user chose not to commit: %w", err)
			}
			logging.Logger.Infof("User chose to commit with original message after edit error.")
			// Proceed with original message
			return a.executeGitCommit(directory, commitMessage)
		}
		logging.Logger.Infof("Commit message edited successfully.")
		return a.executeGitCommit(directory, editedMessage)
	case "N":
		logging.Logger.Infof("User chose not to commit.")
		return nil // User explicitly chose not to commit, not an error state.
	default:
		logging.Logger.Warnf("Invalid choice '%s', changes not committed.", choice)
		return nil // Treat invalid choice like choosing 'No'.
	}
}

// promptForCommit prompts the user to commit changes and returns their choice.
// Returns "Y" for Yes, "N" for No, "A" for Always, "E" for Edit.
func (a *LocalProgrammingAgent) promptForCommit() string {
	logging.Logger.Infof("Waiting for commit confirmation...")

	// Prompt user for commit choice
	commitChoice, err := input.UserInputGetter("Do you want to commit the changes? [Y]es/[N]o/[A]llways/[E]dit ")
	if err != nil {
		logging.Logger.Errorf("Error reading commit choice, defaulting to no commit: %v", err)
		return "N" // Default to no on input error
	}

	// Return the uppercase choice
	return strings.ToUpper(commitChoice)
}

// executeGitCommit executes the git add and git commit commands.
func (a *LocalProgrammingAgent) executeGitCommit(directory string, commitMessage string) error {
	err := a.gitUtil.Add(directory)
	if err != nil {
		return err
	}
	err = a.gitUtil.Commit(directory, commitMessage)

	if err != nil {
		return err
	}
	return nil
}
