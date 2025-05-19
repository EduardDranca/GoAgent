package assistants

import (
	context2 "context"
	"encoding/json"
	errors2 "errors"
	"fmt"
	"github.com/EduardDranca/GoAgent/internal/agent/commands"
	"github.com/EduardDranca/GoAgent/internal/agent/models"
	"github.com/EduardDranca/GoAgent/internal/llm"
	"github.com/EduardDranca/GoAgent/internal/logging"
	"github.com/avast/retry-go"
	"time"
)

type InstructionAssistant interface {
	Instruct(ctx context2.Context, message string) (commands.Command, error)
	ClearHistory()
}

type defaultInstructionAssistant struct {
	session llm.LLMSession
}

func NewInstructionAssistant(session llm.LLMSession) InstructionAssistant {
	return &defaultInstructionAssistant{
		session: session,
	}
}

// ClearHistory clears the history of the instruction agent.
func (a *defaultInstructionAssistant) ClearHistory() {
	a.session.SetHistory([]models.Message{})
}

func (a *defaultInstructionAssistant) Instruct(ctx context2.Context, message string) (commands.Command, error) {
	logging.Logger.Debugf("Starting Instruct function with message: %s", message)

	var commandMap map[string]interface{}
	var response string

	// Use retry-go to handle potential JSON unmarshalling failures
	err := retry.Do(
		func() error {
			logging.Logger.Debug("Attempting to get and unmarshal LLM response")

			var sendErr error
			// Send the message and request JSON format
			response, sendErr = a.session.SendMessage(ctx, message, llm.WithJSON())
			if sendErr != nil {
				logging.Logger.Errorf("SendMessage failed: %v", sendErr)
				// Return the error to retry the SendMessage call
				return sendErr
			}
			logging.Logger.Debugf("Raw response from SendMessage: %s", response)

			// Attempt to unmarshal the response
			unmarshalErr := json.Unmarshal([]byte(response), &commandMap)
			if unmarshalErr == nil {
				logging.Logger.Debug("Successfully unmarshalled JSON")
				return nil // Success, stop retrying
			}

			// If unmarshalling fails, log and return the error to trigger a retry
			logging.Logger.Errorf("Failed to unmarshal command response: %s. Error: %v", response, unmarshalErr)
			return unmarshalErr
		},
		retry.Attempts(3), // Maximum 3 attempts
		retry.DelayType(retry.BackOffDelay),
		retry.Delay(100*time.Millisecond),
		retry.MaxDelay(5*time.Second), // Maximum delay between retries is 5 seconds
		retry.OnRetry(func(n uint, err error) {
			logging.Logger.Warnf("Retry attempt %d failed: %v", n+1, err)
		}),
	)

	// Check the error after the retry loop
	if err != nil {
		// Check if the final error was a JSON unmarshalling error
		var jsonSyntaxErr *json.SyntaxError
		var jsonUnmarshalTypeErr *json.UnmarshalTypeError
		if errors2.As(err, &jsonSyntaxErr) || errors2.As(err, &jsonUnmarshalTypeErr) {
			return nil, fmt.Errorf("failed to unmarshal command response into JSON after multiple retries. Raw response: %s. %w", response, err)
		}
		// Handle other types of errors (e.g., from SendMessage)
		return nil, fmt.Errorf("failed to get valid command response from LLM after multiple retries. %w", err)
	}

	// Proceed with creating the command now that we have a valid commandMap
	cmd, err := commands.NewCommand(commandMap)
	if err != nil {
		logging.Logger.Errorf("Failed to create command: %v", err)
		return nil, fmt.Errorf("failed to create command from response. %w", err)
	}

	logging.Logger.Debugf("Successfully created command of type: %T", cmd)
	return cmd, nil
}
