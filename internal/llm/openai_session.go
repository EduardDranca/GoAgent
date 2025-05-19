package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/EduardDranca/GoAgent/internal/agent/models"
	"github.com/EduardDranca/GoAgent/internal/logging"
	"github.com/openai/openai-go"
)

// OpenAISession implements the LLMSession interface for OpenAI models.
type OpenAISession struct {
	client           *openai.Client
	model            string
	history          []models.Message // Store the history
	systemPrompt     string           // Store the system prompt
	defaultOptions   *Options
	maxHistoryLength int
}

// NewOpenAISession creates a new OpenAISession.
func NewOpenAISession(client *openai.Client, modelName, systemPrompt string, options ...Option) *OpenAISession {
	defaultOptions := createOptions(options...)

	s := &OpenAISession{
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

// SendMessage sends a message to the OpenAI model and returns the response.
func (s *OpenAISession) SendMessage(ctx context.Context, message string, options ...Option) (string, error) {
	opts := &Options{}
	for _, opt := range options {
		opt(opts)
	}

	// Add the incoming message to the history
	userMessage := models.Message{Role: "user", Content: message}
	s.history = append(s.history, userMessage)

	messages := []openai.ChatCompletionMessageParamUnion{
		openai.SystemMessage(s.systemPrompt), // Use the stored system prompt
	}
	// Add the history to the messages
	for _, msg := range s.history {
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    openai.ChatCompletionMessageRole(msg.Role),
			Content: msg.Content,
		})
	}

	req := openai.ChatCompletionNewParams{
		Model:    openai.F(s.model),
		Messages: openai.F(messages),
	}

	// Apply default options
	if s.defaultOptions.Temperature != nil {
		req.Temperature = openai.F(float64(*s.defaultOptions.Temperature))
	}
	if s.defaultOptions.TopP != nil {
		req.TopP = openai.F(float64(*s.defaultOptions.TopP))
	}
	if s.defaultOptions.MaxOutputTokens != nil {
		req.MaxTokens = openai.Int(int64(*s.defaultOptions.MaxOutputTokens))
	}
	if s.defaultOptions.ResponseFormat != nil && *s.defaultOptions.ResponseFormat == "json" {
		if s.defaultOptions.JSONSchema == nil {
			req.ResponseFormat = openai.F[openai.ChatCompletionNewParamsResponseFormatUnion](
				openai.ResponseFormatJSONObjectParam{
					Type: openai.F(openai.ResponseFormatJSONObjectTypeJSONObject),
				},
			)
		} else {
			req.ResponseFormat = openai.F[openai.ChatCompletionNewParamsResponseFormatUnion](
				openai.ResponseFormatJSONSchemaParam{
					Type: openai.F(openai.ResponseFormatJSONSchemaTypeJSONSchema),
					JSONSchema: openai.F(openai.ResponseFormatJSONSchemaJSONSchemaParam{
						Schema: openai.F[interface{}](*s.defaultOptions.JSONSchema),
					}),
				},
			)
		}
	}

	// Override options with any passed into the function call
	if opts.Temperature != nil {
		req.Temperature = openai.F(float64(*opts.Temperature))
	}
	if opts.TopP != nil {
		req.TopP = openai.F(float64(*opts.TopP))
	}
	if opts.MaxOutputTokens != nil {
		req.MaxTokens = openai.Int(int64(*opts.MaxOutputTokens))
	}
	if opts.ResponseFormat != nil && *opts.ResponseFormat == "json" {
		req.Functions = openai.F([]openai.ChatCompletionNewParamsFunction{
			{
				Parameters: openai.F(openai.FunctionParameters(*opts.JSONSchema)),
			},
		})
	}

	resp, err := s.client.Chat.Completions.New(ctx, req)
	if err != nil {
		return "", fmt.Errorf("failed to create chat completion: %w", err)
	}

	if len(resp.Choices) > 0 {
		// Append to history after sending the request
		assistantMessage := models.Message{Role: "assistant", Content: resp.Choices[0].Message.Content}
		s.history = append(s.history, assistantMessage)
		if len(s.history) > s.maxHistoryLength*2 {
			s.history = s.history[2:]
		}
		return resp.Choices[0].Message.Content, nil
	}
	return "", fmt.Errorf("no response from API")
}

// GetHistory returns the conversation history.
func (s *OpenAISession) GetHistory() []models.Message {
	return s.history
}

// SetHistory sets the conversation history.
func (s *OpenAISession) SetHistory(history []models.Message) {
	s.history = history
}
