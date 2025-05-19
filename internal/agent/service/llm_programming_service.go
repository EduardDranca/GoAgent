package service

import (
	context2 "context"
	"fmt"
	"github.com/EduardDranca/GoAgent/internal/agent/assistants"
	"github.com/EduardDranca/GoAgent/internal/agent/commands"
	"github.com/EduardDranca/GoAgent/internal/agent/context"
	"github.com/EduardDranca/GoAgent/internal/input"
	"github.com/EduardDranca/GoAgent/internal/logging"
	"strings"
)

const (
	promptAnalysis = `Which context files, if any, (besides '%s' itself) should be provided to the instruction session to correctly generate the code for this update?
Keep in mind that the llm that generates the code does not have access to the analysis history and will need the content of the files to generate the code even if the file to be generated doesn't reference them directly, the llm will still need to know their content for generation.`
	promptInstruction = `Please construct the final update_file command JSON. Use the following details:
File Path: %s
Implementation Plan: %s
Context file analysis session response: %s
Respond only with the complete JSON object for the command.`
	promptGenerateFileContentExisting = `Make the following changes to this file:
%s

based on this implementation plan:
%s

in the context of this change request:
%s

Using the following file contents as context:
%s`
	promptGenerateFileContentNew = `Create the following file: %s

with the following implementation plan:
%s

in the context of this change request:
%s

Using the following file contents as context:
%s`
	initialPromptImplementContext = "The current project structure is as follows:\n%s\n You are tasked with implementing the following: \n%s"
	initialPromptAskContext       = "The current project structure is as follows:\n%s\n User Query: \n%s"
)

// LLMProgrammingService uses the LLMSession interface for interacting with LLMs.
type LLMProgrammingService struct {
	codeAnalysisAssistant      assistants.AnalysisAssistant
	askAnalysisAssistant       assistants.AnalysisAssistant
	codeInstructionAssistant   assistants.InstructionAssistant
	askInstructionAssistant    assistants.InstructionAssistant
	codeGenerateCodeAssistant  assistants.GenerateCodeAssistant
	patchGenerateCodeAssistant assistants.GenerateCodeAssistant
	maxLoops                   int
}

// NewLLMProgrammingService creates a new instance of LLMProgrammingService.
func NewLLMProgrammingService(
	codeAnalysisAssistant assistants.AnalysisAssistant,
	askAnalysisAssistant assistants.AnalysisAssistant,
	codeInstructionAssistant assistants.InstructionAssistant,
	askInstructionAssistant assistants.InstructionAssistant,
	codeGenerateCodeAssistant assistants.GenerateCodeAssistant,
	patchGenerateCodeAssistant assistants.GenerateCodeAssistant,
	maxLoops int,
) *LLMProgrammingService {
	logging.Logger.Infof("Creating new LLMProgrammingService")
	return &LLMProgrammingService{
		codeAnalysisAssistant:      codeAnalysisAssistant,
		askAnalysisAssistant:       askAnalysisAssistant,
		codeInstructionAssistant:   codeInstructionAssistant,
		askInstructionAssistant:    askInstructionAssistant,
		codeGenerateCodeAssistant:  codeGenerateCodeAssistant,
		patchGenerateCodeAssistant: patchGenerateCodeAssistant,
		maxLoops:                   maxLoops,
	}
}

// ImplementWithContext performs implementation using provided context and LLM sessions.
func (s *LLMProgrammingService) ImplementWithContext(agentContext context.ProgrammingAgentContext) (string, error) {
	logging.Logger.Infof("Starting ImplementWithContext")

	defer s.codeInstructionAssistant.ClearHistory()

	// Create the initial prompt
	initialPrompt := fmt.Sprintf(initialPromptImplementContext, agentContext.GetRepoStructure(), agentContext.GetChangeRequest())

	// Process the initial prompt
	response, err := s.processRequest(initialPrompt, agentContext, true)
	if err != nil {
		return "", fmt.Errorf("error processing request in ImplementWithContext: %w", err)
	}
	return response, nil
}

// AskWithContext performs asking using provided context and LLM sessions, without file modifications.
func (s *LLMProgrammingService) AskWithContext(agentContext context.ProgrammingAgentContext) (string, error) {
	logging.Logger.Infof("Starting AskWithContext")

	defer s.askInstructionAssistant.ClearHistory()

	// Create the initial prompt
	initialPrompt := fmt.Sprintf(initialPromptAskContext, agentContext.GetRepoStructure(), agentContext.GetChangeRequest())

	// Process the initial prompt
	response, err := s.processRequest(initialPrompt, agentContext, false)
	if err != nil {
		return "", fmt.Errorf("error processing request in AskWithContext: %w", err)
	}
	return response, nil
}

// processRequest encapsulates the shared logic for AskWithContext and ImplementWithContext.
func (s *LLMProgrammingService) processRequest(initialPrompt string, agentContext context.ProgrammingAgentContext, useImplementSessions bool) (string, error) {
	logging.Logger.Debugf("Starting processRequest")

	resp, err := s.sendMessage(initialPrompt, useImplementSessions)
	if err != nil {
		return "", fmt.Errorf("error sending initial message in processRequest: %w", err)
	}

	loopCounter := 0
	for {
		loopCounter++
		if loopCounter > s.maxLoops {
			userInput, err := input.UserInputGetter(fmt.Sprintf("The process has run for %d loops. Do you want to continue? [Y]es/[N]o ", s.maxLoops))
			if err != nil {
				logging.Logger.Errorf("Error getting user input: %v. Stopping process.", err)
				return "Process stopped by user after loop limit.", nil
			}

			userInput = strings.ToLower(strings.TrimSpace(userInput))
			if userInput == "yes" || userInput == "y" {
				loopCounter = 0 // Reset loop counter to continue
				continue
			} else if userInput == "no" || userInput == "n" {
				return "Process stopped by user after loop limit.", nil
			} else {
				logging.Logger.Warnf("Invalid user input: %s. Stopping process.", userInput)
				return "Process stopped by user after loop limit.", nil
			}
		}

		processedResponse, isCommand, isFinalCommand, err := s.processCommand(resp, agentContext)
		if err != nil {
			logging.Logger.Errorf("Error processing command in processRequest: %v", err)
		}

		if isFinalCommand {
			logging.Logger.Infof("Received final command, task complete.")
			return processedResponse, nil
		}

		if isCommand {
			resp, err = s.sendMessage(processedResponse, useImplementSessions)
			if err != nil {
				return "", fmt.Errorf("error sending message in processRequest: %w", err)
			}
		} else {
			// Handle non-command responses if needed, for now, just log them or send them back for context
			logging.Logger.Warnf("Received non-command response from LLM: %s. Forwarding as context.", processedResponse)
			resp, err = s.sendMessage(processedResponse, useImplementSessions)
			if err != nil {
				return "", fmt.Errorf("error sending message in processRequest for non-command response: %w", err)
			}
		}
	}
}

// sendMessage sends messages to the appropriate sessions based on sessionType ('ask' or 'implement').
func (s *LLMProgrammingService) sendMessage(processedResponse string, useImplementSessions bool) (commands.Command, error) {
	var analysisAssistant assistants.AnalysisAssistant
	var instructionAssistant assistants.InstructionAssistant

	if useImplementSessions {
		analysisAssistant = s.codeAnalysisAssistant
		instructionAssistant = s.codeInstructionAssistant
	} else {
		analysisAssistant = s.askAnalysisAssistant
		instructionAssistant = s.askInstructionAssistant
	}

	return s.executeLLMAssistant(processedResponse, analysisAssistant, instructionAssistant)
}

// executeLLMAssistant encapsulates the common logic for executing analysis and instruction assistants.
func (s *LLMProgrammingService) executeLLMAssistant(processedResponse string, analysisAssistant assistants.AnalysisAssistant, instructionAssistant assistants.InstructionAssistant) (commands.Command, error) {
	analysisResponse, err := analysisAssistant.Execute(context2.Background(), processedResponse)
	if err != nil {
		return nil, fmt.Errorf("error during analysis agent execution in executeLLMAssistant: %w", err)
	}
	instructionResponseCommand, err := instructionAssistant.Instruct(context2.Background(), analysisResponse)
	if err != nil {
		return nil, fmt.Errorf("error during instruction agent execution in executeLLMAssistant: %w", err)
	}
	return instructionResponseCommand, nil
}

// processCommand processes the LLM response and returns the processed response, a boolean indicating if it's a command, and an error
func (s *LLMProgrammingService) processCommand(command commands.Command, agentContext context.ProgrammingAgentContext) (string, bool, bool, error) {
	logging.Logger.Debugf("Starting processCommand")
	commandType := fmt.Sprintf("%T", command)
	logging.Logger.Debugf("Processing command of type: %s", commandType)

	// Execute the command
	processedResponse, err := s.executeCommand(command, agentContext) // Pass commandMap directly

	isFinalCommand := false
	_, isCommit := command.(*commands.CommitCommand)
	_, isRespond := command.(*commands.RespondCommand)
	if isCommit || isRespond {
		isFinalCommand = true
		logging.Logger.Debugf("Command is a final command (commit or respond)")
		return processedResponse, true, isFinalCommand, nil
	}

	if err != nil {
		// In case of an error, set processedResponse to an empty string.
		// The error itself should be checked by the caller to determine the failure.
		return processedResponse, true, isFinalCommand, fmt.Errorf("error executing command in processCommand: %w", err)
	}
	return processedResponse, true, isFinalCommand, nil
}

// executeCommand executes a given command.
func (s *LLMProgrammingService) executeCommand(command commands.Command, agentContext context.ProgrammingAgentContext) (string, error) { // Changed to accept commandMap
	commandType := fmt.Sprintf("%T", command)
	logging.Logger.Debugf("Starting executeCommand for command type: %s", commandType)

	updateCommand, isUpdate := command.(*commands.UpdateFileCommand)

	if isUpdate {
		var err error
		command, err = s.handleFileUpdate(updateCommand, agentContext) // Pass commandMap to handleFileUpdate
		if err != nil {
			return "File update failed, please retry.", fmt.Errorf("error handling file update in executeCommand: %w", err)
		}
	}
	processedResponse, err := command.Process(agentContext)
	if err != nil {
		// Replace errors.New with fmt.Errorf
		wrappedErr := fmt.Errorf("commandError processing command %s: %w", command, err)
		return wrappedErr.Error(), wrappedErr
	}
	return processedResponse, nil
}

// buildContextFilePromptComponent constructs the context file prompt component.
func (s *LLMProgrammingService) buildContextFilePromptComponent(agentContext context.ProgrammingAgentContext, contextFiles []string, file string) string {
	var contextFilePromptComponent string

	for _, contextFile := range contextFiles {
		if contextFile == file {
			continue
		}
		content, exists := agentContext.GetFileContent(contextFile)
		if !exists {
			logging.Logger.Infof("Context file does not exist: %s", contextFile)
			continue // Skip to the next context file if this one doesn't exist
		}
		if contextFilePromptComponent == "" {
			contextFilePromptComponent = "" // Start empty, add header later if needed
		}
		contextFilePromptComponent += fmt.Sprintf("File: %s\n%s\n", contextFile, content)
	}
	return contextFilePromptComponent
}

// handleFileUpdate handles the file update command.
func (s *LLMProgrammingService) handleFileUpdate(updateCommand *commands.UpdateFileCommand, agentContext context.ProgrammingAgentContext) (commands.Command, error) { // Changed to accept commandMap
	filePath := updateCommand.FilePath                     // Extract file_path from commandMap
	implementationPlan := updateCommand.ImplementationPlan // Extract implementation_plan from commandMap
	logging.Logger.Debugf("Starting handleFileUpdate for file: %s", filePath)

	analysisPrompt := fmt.Sprintf(promptAnalysis, filePath)
	analysisResponse, err := s.codeAnalysisAssistant.Execute(context2.Background(), analysisPrompt)

	if err != nil {
		return nil, fmt.Errorf("error prompting analysis LLM for context files in handleFileUpdate: %w", err)
	}

	instructionPrompt := fmt.Sprintf(promptInstruction, filePath, implementationPlan, analysisResponse)
	instructionResponse, err := s.codeInstructionAssistant.Instruct(context2.Background(), instructionPrompt)

	if err != nil {
		return nil, fmt.Errorf("error prompting instruction LLM to construct final update_file command in handleFileUpdate: %w", err)
	}

	finalUpdateCmd, ok := instructionResponse.(*commands.UpdateFileCommand)
	if !ok {
		return nil, fmt.Errorf("error creating final update_file command in handleFileUpdate: %w", err)
	}

	err = s.generateFileContent(finalUpdateCmd.ImplementationPlan, finalUpdateCmd.FilePath, finalUpdateCmd.ContextFiles, agentContext)
	if err != nil {
		return nil, fmt.Errorf("error generating file content in handleFileUpdate: %w", err) // Return error from generateFileContent
	}
	return instructionResponse, nil
}

// generateFileContent generates the content of each file based on the implementation plan.
func (s *LLMProgrammingService) generateFileContent(implementationPlan, file string, contextFiles []string, agentContext context.ProgrammingAgentContext) error {
	logging.Logger.Infof("Starting generateFileContent for file: %s", file)

	contextFilePromptComponent := s.buildContextFilePromptComponent(agentContext, contextFiles, file)

	existingFileContent, exists := agentContext.GetFileContent(file)

	var prompt string

	if exists {
		prompt = fmt.Sprintf(
			promptGenerateFileContentExisting,
			existingFileContent,
			implementationPlan,
			agentContext.GetChangeRequest(),
			contextFilePromptComponent)
	} else {
		prompt = fmt.Sprintf(
			promptGenerateFileContentNew,
			file,
			implementationPlan,
			agentContext.GetChangeRequest(),
			contextFilePromptComponent)
	}

	logging.Logger.Debugf("Generating content for file: %s", file)

	codeGenerated, err := s.codeGenerateCodeAssistant.GenerateCode(context2.Background(), prompt)

	if err != nil {
		logging.Logger.Errorf("Error from generateCodeAgent.Execute: %v", err)
		return err
	} else {

		appliedPatch, patchErr := s.applyPatch(existingFileContent, codeGenerated, file, agentContext)
		if patchErr != nil {
			logging.Logger.Errorf("Error applying patch for file %s: %v", file, patchErr)
		}
		if appliedPatch {
			return nil // Patch applied successfully, content updated in applyPatch
		}

		logging.Logger.Debugf("Generated content for file: %s", file)
		agentContext.UpdateFileContent(file, codeGenerated)
	}

	return nil
}

func (s *LLMProgrammingService) applyPatch(existingFileContent string, extractedContent string, file string, agentContext context.ProgrammingAgentContext) (bool, error) {
	isPatch := strings.HasPrefix(extractedContent, "--- a/")
	if !isPatch {
		return false, nil
	}

	logging.Logger.Infof("Detected git patch response, attempting to apply patch for file: %s", file)
	if existingFileContent == "" {
		logging.Logger.Warnf("Cannot apply patch to a non-existent file: %s. Falling back to code block extraction.", file)
		return false, nil
	}

	logging.Logger.Debugf("Calling patchGenerateCodeAssistant.GenerateCode for file: %s", file)
	// The patch assistant is specifically designed to take existing content and a patch and return the new content.
	// The prompt format here is crucial for the patch assistant to understand the input.
	patchPrompt := fmt.Sprintf("Apply the following git patch to the file content:\n\nFile Content:\n```\n%s\n```\n\nGit Patch:\n```diff\n%s\n```\n\nProvide only the resulting file content.", existingFileContent, extractedContent)
	patchedContent, patchErr := s.patchGenerateCodeAssistant.GenerateCode(context2.Background(), patchPrompt)
	if patchErr != nil {
		logging.Logger.Errorf("Error applying patch using patchGenerateCodeAssistant.GenerateCode for file %s: %v.", file, patchErr)
		return false, patchErr
	}

	logging.Logger.Infof("Successfully applied patch using patchGenerateCodeAssistant.GenerateCode for file %s.", file)
	agentContext.UpdateFileContent(file, patchedContent)

	return true, nil
}
