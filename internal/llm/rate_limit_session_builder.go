package llm

import (
	"context"
	"fmt"

	"github.com/EduardDranca/GoAgent/internal/config"

	"golang.org/x/time/rate"

	"github.com/google/generative-ai-go/genai"
	"github.com/jpoz/groq"
	"github.com/openai/openai-go"
	option2 "github.com/openai/openai-go/option"
	"google.golang.org/api/option"
)

// NewRateLimitSessionBuilder builds a rate-limited LLM session based on the given LLM service type.
func NewRateLimitSessionBuilder(ctx context.Context, llmType config.LLMServiceType, apiKey string, modelName string, rateLimiter *rate.Limiter, systemMessage string, options ...Option) (LLMSession, error) {
	var baseSession LLMSession

	switch llmType {
	case config.GeminiService:
		genaiClient, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
		if err != nil {
			return nil, fmt.Errorf("failed to create Gemini client: %w", err)
		}
		baseSession, err = NewGeminiSession(genaiClient, modelName, systemMessage, options...)
		if err != nil {
			return nil, fmt.Errorf("failed to create Gemini session: %w", err)
		}
	case config.GroqService:
		groqClient := groq.NewClient(groq.WithAPIKey(apiKey))
		baseSession = NewGroqSession(groqClient, modelName, systemMessage, options...)
	case config.OpenAIService:
		openaiClient := openai.NewClient(option2.WithAPIKey(apiKey))
		baseSession = NewOpenAISession(openaiClient, modelName, systemMessage, options...)
	}

	rateLimitedSession := NewRateLimitSession(baseSession, rateLimiter)
	return rateLimitedSession, nil
}
