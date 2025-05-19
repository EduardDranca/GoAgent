package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/EduardDranca/GoAgent/internal/agent/models"
	"github.com/EduardDranca/GoAgent/internal/logging"
	"github.com/avast/retry-go"
	"github.com/jpoz/groq"
	"time"
)

// GroqSession implements the LLMSession interface for Groq models.
type GroqSession struct {
	client           *groq.Client
	model            string
	history          []models.Message // Store the history
	systemPrompt     string           // Store the system prompt
	defaultOptions   *Options
	maxHistoryLength int
}

// NewGroqSession creates a new GroqSession. It now accepts the groq.Client as a parameter.
func NewGroqSession(client *groq.Client, modelName string, systemPrompt string, options ...Option) *GroqSession {

	defaultOptions := createOptions(options...)

	s := &GroqSession{
		client:         client,
		model:          modelName,
		history:        []models.Message{},
		systemPrompt:   systemPrompt, // Store the system prompt
		defaultOptions: defaultOptions,
	}

	s.maxHistoryLength = *defaultOptions.MaxHistoryLength
	if s.maxHistoryLength == -1 {
		s.maxHistoryLength = 100 // Default value if not set
	}

	if defaultOptions != nil && defaultOptions.JSONSchema != nil {
		schemaJSON, err := json.MarshalIndent(*defaultOptions.JSONSchema, "", "  ")
		if err != nil {
			logging.Logger.Errorf("failed to marshal JSON schema: %v", err)
		} else {
			s.systemPrompt = fmt.Sprintf("%s\n\nPlease respond using the following JSON schema:\n%s", s.systemPrompt, string(schemaJSON))
		}
	}
	return s
}

// SendMessage sends a message to the Groq model and returns the response.
func (s *GroqSession) SendMessage(_ context.Context, message string, options ...Option) (string, error) {
	opts := &Options{}
	for _, opt := range options {
		opt(opts)
	}

	// Add the incoming message to the history
	userMessage := models.Message{Role: "user", Content: message}
	s.history = append(s.history, userMessage)

	messages := []groq.Message{
		{
			Role:    "system",
			Content: s.systemPrompt, // Use the stored system prompt
		},
	}
	// Add the history to the messages
	for _, msg := range s.history {
		messages = append(messages, groq.Message{Role: msg.Role, Content: msg.Content})
	}

	req := groq.CompletionCreateParams{
		Model:    s.model,
		Messages: messages,
	}

	// Apply default options
	if s.defaultOptions.Temperature != nil {
		req.Temperature = *s.defaultOptions.Temperature
	}
	if s.defaultOptions.TopP != nil {
		req.TopP = *s.defaultOptions.TopP
	}
	if s.defaultOptions.MaxOutputTokens != nil {
		req.MaxTokens = *s.defaultOptions.MaxOutputTokens
	}
	if s.defaultOptions.ResponseFormat != nil && *s.defaultOptions.ResponseFormat == "json" {
		req.ResponseFormat = groq.ResponseFormat{
			Type: "json_object",
		}
	}

	//Override options with any passed into the function call
	if opts.Temperature != nil {
		req.Temperature = *opts.Temperature
	}
	if opts.TopP != nil {
		req.TopP = *opts.TopP
	}
	if opts.MaxOutputTokens != nil {
		req.MaxTokens = *opts.MaxOutputTokens
	}
	if opts.ResponseFormat != nil && *opts.ResponseFormat == "json" {
		req.ResponseFormat = groq.ResponseFormat{
			Type: "json_object",
		}
	}

	// Handle JSON schema for Groq (add to system prompt)
	if opts.JSONSchema != nil {
		schemaJSON, err := json.MarshalIndent(*opts.JSONSchema, "", "  ")
		if err != nil {
			return "", fmt.Errorf("failed to marshal JSON schema: %w", err)
		}
		// Modify the system message content
		messages[0].Content = fmt.Sprintf("%s\n\nPlease respond using the following JSON schema:\n%s", messages[0].Content, string(schemaJSON))
	}

	var resp *groq.ChatCompletion
	// TODO: Added retry since the groq session sometimes fails, should investigate at a later date.
	err := retry.Do(
		func() error {
			var rErr error
			resp, rErr = s.client.CreateChatCompletion(req)
			return rErr
		},
		retry.Attempts(3), // Maximum 3 attempts
		retry.DelayType(retry.BackOffDelay),
		retry.Delay(100*time.Millisecond),
		retry.MaxDelay(5*time.Second), // Maximum delay between retries is 5 seconds
	)
	if err != nil {
		return "", fmt.Errorf("failed to create chat completion after multiple retries: %w", err)
	}

	assistantMessage := models.Message{Role: "assistant", Content: resp.Choices[0].Message.Content}
	//Append to history after sending the request
	s.history = append(s.history, assistantMessage)

	if len(s.history) > s.maxHistoryLength*2 {
		s.history = s.history[2:]
	}

	return resp.Choices[0].Message.Content, nil
}

// GetHistory returns the conversation history.
func (s *GroqSession) GetHistory() []models.Message {
	return s.history
}

func (s *GroqSession) SetHistory(history []models.Message) {
	s.history = history
}
