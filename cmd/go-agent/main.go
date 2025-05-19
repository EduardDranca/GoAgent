package main

import (
	"context"
	"errors"
	"strings"

	"github.com/fatih/color"
	"github.com/reeflective/readline" // Use the new library

	"github.com/EduardDranca/GoAgent/internal/agent"
	"github.com/EduardDranca/GoAgent/internal/agent/initialize"
	"github.com/EduardDranca/GoAgent/internal/agent/models"
	"github.com/EduardDranca/GoAgent/internal/agent/service"
	"github.com/EduardDranca/GoAgent/internal/config"
	"github.com/EduardDranca/GoAgent/internal/input"
	"github.com/EduardDranca/GoAgent/internal/input/completer"
	"github.com/EduardDranca/GoAgent/internal/logging"
	"github.com/EduardDranca/GoAgent/internal/utils"
)

const (
	CommandAsk       = "/ask"
	CommandImplement = "/implement"
)

func main() {
	logging.Logger.Info("Starting GoAgent...")
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		logging.Logger.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize logging with log level from config
	err = logging.InitializeLogging(cfg.LogLevel)
	logging.Logger.Debugf("Log level set to %s", cfg.LogLevel)
	if err != nil {
		logging.Logger.Warn("Failed to initialize logging: %v", err)
	}
	defer logging.CloseLogger()

	// Set Glamour style path from config to utils package
	utils.SetGlamourStylePath(string(cfg.GlamourStylePath)) // Set the global glamourStylePath

	logging.Logger.Infof("Configuration loaded successfully. Log level: %s", cfg.LogLevel)
	// Initialize context with cancel for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize services
	programmingService, err := initialize.InitProgrammingService(ctx, cfg)
	if err != nil {
		logging.Logger.Fatalf("Failed to initialize services: %v", err)
	}
	logging.Logger.Infof("Programming service initialized successfully.")

	// Run the application based on the specified mode
	runService(programmingService, cfg.Directory)
}

// runService runs the application in local mode
func runService(programmingService service.ProgrammingService, directory string) {
	logging.Logger.Infof("Starting runService in directory: %s", directory)
	if directory == "" {
		logging.Logger.Fatalf("Error: -d option missing; directory must be a git repository")
	}

	isDir, err := utils.IsDirectory(directory)
	if err != nil {
		logging.Logger.Fatalf("Error: failed to access directory %s with error %v", directory, err)
	}

	if !isDir {
		logging.Logger.Fatalf("Error: value provided for -d argument must be a valid directory")
	}

	// Initialize gitUtil based on whether the directory is a git repository
	gitUtil := initGitUtil(directory)

	currentWordCompleter, err := completer.InitCompleter(directory, gitUtil)
	if err != nil {
		logging.Logger.Errorf("Failed to initialize current word completer: %v. File name completion won't be available", err)
	}

	programmingAgent := agent.NewLocalProgrammingAgent(programmingService, gitUtil)

	runChangeRequestLoop(directory, programmingAgent, currentWordCompleter, gitUtil)
}

func runChangeRequestLoop(directory string, programmingAgent agent.AgentInterface[models.AgentRequest], currentWordCompleter *completer.CurrentWordCompleter, gitUtil utils.GitUtil) {
	for {
		logging.Logger.Infof("Waiting for change request...")
		changeRequest, err := input.GetLocalChangeRequest(currentWordCompleter)
		if err != nil {
			if errors.Is(err, readline.ErrInterrupt) {
				logging.Logger.Infof("Program interrupted by user.")
				break
			}
			logging.Logger.Errorf("Error reading change request: %v", err)
			continue // Go to the next iteration
		}
		logging.Logger.Debugf("Received change request: %s", changeRequest)

		processLocalChangeRequest(directory, changeRequest, programmingAgent)

		// Re-initialize completer in case files changed
		currentWordCompleter, err = completer.InitCompleter(directory, gitUtil)
		if err != nil {
			logging.Logger.Errorf("Failed to initialize current word completer: %v", err)
			// Continue even if completer fails, just without completion
		}
	}
}

// initGitUtil initializes GitUtil based on whether the directory is a git repository.
func initGitUtil(directory string) utils.GitUtil {
	// Check if the directory is a git repository
	var gitUtil utils.GitUtil
	isGitRepo, err := utils.IsGitRepository(directory)
	if err != nil {
		logging.Logger.Infof("Error checking if directory is a git repository: %v", err)
		isGitRepo = false
	}

	if isGitRepo {
		gitUtil = &utils.RealGitUtil{}
	} else {
		gitUtil = &utils.NoOpGitUtil{}
	}
	return gitUtil
}

// parseCommand extracts the command and argument from the input string.
// If the string starts with '/', it's treated as a command. Otherwise, it's
// treated as an implicit '/implement' command.
func parseCommand(input string) (command string, argument string) {
	input = strings.TrimSpace(input)
	if strings.HasPrefix(input, "/") {
		parts := strings.FieldsFunc(input, func(r rune) bool {
			return r == ' ' || r == '\t'
		})
		if len(parts) > 0 {
			command = parts[0]
			argument = strings.TrimSpace(strings.TrimPrefix(input, command))
		} else {
			// Should not happen with TrimSpace, but handle defensively
			command = input
			argument = ""
		}
	} else {
		command = CommandImplement // Default command
		argument = input
	}
	return command, argument
}

// processLocalChangeRequest parses the change request and dispatches to the appropriate handler.
func processLocalChangeRequest(directory string, changeRequest string, programmingAgent agent.AgentInterface[models.AgentRequest]) {
	logging.Logger.Debugf("Processing change request: %s", changeRequest)

	command, argument := parseCommand(changeRequest)

	switch command {
	case CommandAsk:
		handleAskCommand(directory, argument, programmingAgent)
	case CommandImplement:
		handleImplementCommand(directory, argument, programmingAgent)
	default:
		logging.Logger.Errorf("Error: unknown command '%s'. Supported commands: %s, %s", command, CommandAsk, CommandImplement)
	}
}

// handleAskCommand processes the /ask command.
func handleAskCommand(directory string, query string, programmingAgent agent.AgentInterface[models.AgentRequest]) {
	logging.Logger.Debugf("Handling %s command with query: %s", CommandAsk, query)
	result, err := programmingAgent.Ask(models.AgentRequest{
		Query:     query,
		Directory: directory,
	})
	if err != nil {
		logging.Logger.Errorf("Error: failed to ask agent: %v", err)
		return
	}

	out, err := utils.RenderWithGlamour(result)
	if err != nil {
		// Fallback to printing raw result if glamour rendering fails
		logging.Logger.Infof(color.HiWhiteString(result))
	} else {
		logging.Logger.Infof(out) // Print glamour-rendered output to stdout
	}
}

// handleImplementCommand processes the /implement command (default).
func handleImplementCommand(directory string, changeRequest string, programmingAgent agent.AgentInterface[models.AgentRequest]) {
	logging.Logger.Debugf("Handling %s command with request: %s", CommandImplement, changeRequest)
	err := programmingAgent.Implement(models.AgentRequest{
		Query:     changeRequest,
		Directory: directory,
	})
	if err != nil {
		logging.Logger.Errorf("Error: failed to implement change request: %v", err)
		return
	}
	logging.Logger.Infof("Change request processed successfully.")
}
