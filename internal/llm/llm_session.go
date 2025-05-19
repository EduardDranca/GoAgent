package llm

import (
	"context"
	"github.com/EduardDranca/GoAgent/internal/agent/models"
)

// LLMSession interface defines the common methods for interacting with different LLMs.
type LLMSession interface {
	SendMessage(ctx context.Context, message string, options ...Option) (string, error) // Change this line
	GetHistory() []models.Message
	SetHistory(history []models.Message)
}

// Option is a functional option type for configuring LLM behavior.
type Option func(o *Options)

// Options holds the configurable parameters for LLM interactions.
type Options struct {
	Temperature      *float32
	TopK             *int
	TopP             *float32
	MaxOutputTokens  *int
	ResponseFormat   *string                 // "json" or "text"
	JSONSchema       *map[string]interface{} // For Gemini, this will be used to construct genai.Schema. For Groq, this will be added to the system prompt.
	MaxHistoryLength *int
}

// WithTemperature sets the temperature for the LLM.
func WithTemperature(temperature float32) Option {
	return func(o *Options) {
		o.Temperature = &temperature
	}
}

// WithTopK sets the TopK value for the LLM.
func WithTopK(topK int) Option {
	return func(o *Options) {
		o.TopK = &topK
	}
}

// WithTopP sets the TopP value for the LLM.
func WithTopP(topP float32) Option {
	return func(o *Options) {
		o.TopP = &topP
	}
}

// WithResponseFormat sets the desired response format ("json" or "text").
func WithResponseFormat(format string) Option {
	return func(o *Options) {
		o.ResponseFormat = &format
	}
}

// WithJSONSchema sets the JSON schema for the LLM to follow.
func WithJSONSchema(schema map[string]interface{}) Option {
	return func(o *Options) {
		o.JSONSchema = &schema
	}
}

// WithMaxHistoryLength sets the maximum history length for the LLM session.
func WithMaxHistoryLength(length int) Option {
	return func(o *Options) {
		o.MaxHistoryLength = &length
	}
}

func WithJSON() Option {
	return WithResponseFormat("json")
}

func createOptions(options ...Option) *Options {
	opts := &Options{}

	// Apply function call options, which will override the default options
	for _, opt := range options {
		opt(opts)
	}
	return opts
}
