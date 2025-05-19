package llm

import (
	"context"
	"fmt"
	"github.com/EduardDranca/GoAgent/internal/agent/models"
	"github.com/google/generative-ai-go/genai"
)

// GeminiSession implements the LLMSession interface for Google's Gemini models.
type GeminiSession struct {
	client           *genai.Client
	model            *genai.GenerativeModel
	chat             *genai.ChatSession // Store the chat session
	defaultOptions   *Options
	maxHistoryLength int
}

// NewGeminiSession creates a new GeminiSession. It now accepts the genai.Client as a parameter.
func NewGeminiSession(client *genai.Client, modelName string, systemPrompt string, options ...Option) (*GeminiSession, error) {
	model := client.GenerativeModel(modelName)
	chat := model.StartChat() // Start the chat session here

	// Set the system instruction.
	model.SystemInstruction = &genai.Content{
		Parts: []genai.Part{
			genai.Text(systemPrompt),
		},
	}

	defaultOptions := createOptions(options...)

	s := &GeminiSession{
		client:         client,
		model:          model,
		chat:           chat, // Store the chat session
		defaultOptions: defaultOptions,
	}

	s.maxHistoryLength = *defaultOptions.MaxHistoryLength
	if s.maxHistoryLength == -1 {
		s.maxHistoryLength = 100 // Default value if not set
	}

	return s, nil
}

// SendMessage sends a message to the Gemini model and returns the response.
func (s *GeminiSession) SendMessage(ctx context.Context, message string, options ...Option) (string, error) {
	if s.defaultOptions != nil {
		// Apply default options first
		if s.defaultOptions.Temperature != nil {
			s.model.SetTemperature(*s.defaultOptions.Temperature)
		}
		if s.defaultOptions.TopK != nil {
			s.model.SetTopK(int32(*s.defaultOptions.TopK))
		}
		if s.defaultOptions.TopP != nil {
			s.model.SetTopP(*s.defaultOptions.TopP)
		}
		if s.defaultOptions.MaxOutputTokens != nil {
			s.model.SetMaxOutputTokens(int32(*s.defaultOptions.MaxOutputTokens))
		}
		if s.defaultOptions.ResponseFormat != nil {
			s.model.ResponseMIMEType = "application/" + *s.defaultOptions.ResponseFormat
		}

		// Handle JSON schema for Gemini
		if s.defaultOptions.JSONSchema != nil {
			s.model.ResponseSchema = convertToGeminiSchema(*s.defaultOptions.JSONSchema)
		}
	}

	opts := createOptions(options...)

	// Set model parameters based on options
	if opts.Temperature != nil {
		s.model.SetTemperature(*opts.Temperature)
	}
	if opts.TopK != nil {
		s.model.SetTopK(int32(*opts.TopK))
	}
	if opts.TopP != nil {
		s.model.SetTopP(*opts.TopP)
	}
	if opts.MaxOutputTokens != nil {
		s.model.SetMaxOutputTokens(int32(*opts.MaxOutputTokens))
	}
	if opts.ResponseFormat != nil {
		s.model.ResponseMIMEType = "application/" + *opts.ResponseFormat
	}

	// Handle JSON schema for Gemini
	if opts.JSONSchema != nil {
		s.model.ResponseSchema = convertToGeminiSchema(*opts.JSONSchema)
	}

	resp, err := s.chat.SendMessage(ctx, genai.Text(message))
	if err != nil {
		return "", fmt.Errorf("failed to send message to Gemini: %w", err)
	}
	responseText := extractResponseText(resp)

	currentHistory := s.GetHistory()
	if s.maxHistoryLength > 0 && len(currentHistory) > s.maxHistoryLength*2 {
		s.SetHistory(currentHistory[2:])
	}

	return responseText, nil
}

// GetHistory returns the conversation history.
func (s *GeminiSession) GetHistory() []models.Message {
	return historyToMessages(s.chat.History)
}

// extractResponseText extracts the text from a GenerateContentResponse.
func extractResponseText(resp *genai.GenerateContentResponse) string {
	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return ""
	}
	response := ""
	for _, part := range resp.Candidates[0].Content.Parts {
		text, ok := part.(genai.Text)
		if ok {
			response = response + string(text)
		}
	}
	return response
}

// convertToGeminiSchema converts a generic JSON schema to a genai.Schema.
func convertToGeminiSchema(schema map[string]interface{}) *genai.Schema {
	genaiSchema := &genai.Schema{
		Type:       genai.TypeObject, // Assuming the top level is always an object
		Properties: make(map[string]*genai.Schema),
	}

	for key, value := range schema["properties"].(map[string]interface{}) {
		if fieldSchema, ok := value.(map[string]interface{}); ok {
			genaiSchema.Properties[key] = convertSchemaField(fieldSchema)
		}
	}
	requiredFields := schema["required"]
	if required, ok := requiredFields.([]interface{}); ok {
		for _, field := range required {
			if fieldStr, ok := field.(string); ok {
				genaiSchema.Required = append(genaiSchema.Required, fieldStr)
			}
		}
	}
	return genaiSchema
}

// convertSchemaField recursively converts a field within the schema.
func convertSchemaField(fieldSchema map[string]interface{}) *genai.Schema {
	s := &genai.Schema{}

	if typ, ok := fieldSchema["type"].(string); ok {
		switch typ {
		case "string":
			s.Type = genai.TypeString
		case "integer":
			s.Type = genai.TypeInteger
		case "number":
			s.Type = genai.TypeNumber
		case "boolean":
			s.Type = genai.TypeBoolean
		case "array":
			s.Type = genai.TypeArray
			if items, ok := fieldSchema["items"].(map[string]interface{}); ok {
				s.Items = convertSchemaField(items)
			}
		case "object":
			s.Type = genai.TypeObject
			if properties, ok := fieldSchema["properties"].(map[string]interface{}); ok {
				s.Properties = make(map[string]*genai.Schema)
				for key, value := range properties {
					if propSchema, ok := value.(map[string]interface{}); ok {
						s.Properties[key] = convertSchemaField(propSchema)
					}
				}
			}
		}
		//TODO: Add more types as needed
	}
	if description, ok := fieldSchema["description"].(string); ok {
		s.Description = description
	}
	if enum, ok := fieldSchema["enum"].([]interface{}); ok {
		for _, e := range enum {
			s.Enum = append(s.Enum, fmt.Sprintf("%v", e))
		}
	}
	if properties, ok := fieldSchema["properties"].(map[string]interface{}); ok {
		s.Properties = make(map[string]*genai.Schema)
		for key, value := range properties {
			if propSchema, ok := value.(map[string]interface{}); ok {
				s.Properties[key] = convertSchemaField(propSchema)
			}
		}
	}

	return s
}

// historyToMessages converts []*genai.Content to []models.Message.  Moved here.
func historyToMessages(history []*genai.Content) []models.Message {
	var messages []models.Message
	for _, h := range history {
		if len(h.Parts) > 0 {
			// Assuming only Text parts, as that's what we use
			if text, ok := h.Parts[0].(genai.Text); ok {
				messages = append(messages, models.Message{
					Role:    h.Role,
					Content: string(text),
				})
			}
		}
	}
	return messages
}

func (s *GeminiSession) SetHistory(history []models.Message) {
	var genaiHistory []*genai.Content
	for _, msg := range history {
		genaiHistory = append(genaiHistory, &genai.Content{
			Role:  msg.Role,
			Parts: []genai.Part{genai.Text(msg.Content)},
		})
	}
	s.chat.History = genaiHistory
}
