package llm

import (
	"context"
	"fmt"
	"github.com/EduardDranca/GoAgent/internal/agent/models"

	"golang.org/x/time/rate"
)

type RateLimitSession struct {
	llmSession        LLMSession
	requestsPerMinute float64
	rateLimiter       *rate.Limiter
}

func NewRateLimitSession(llmSession LLMSession, rateLimiter *rate.Limiter) *RateLimitSession {
	return &RateLimitSession{
		llmSession:        llmSession,
		rateLimiter:       rateLimiter,
		requestsPerMinute: float64(rateLimiter.Limit()),
	}
}

func (rl *RateLimitSession) GetHistory() []models.Message {
	return rl.llmSession.GetHistory()
}

func (rl *RateLimitSession) SetHistory(history []models.Message) {
	rl.llmSession.SetHistory(history)
}

func (rl *RateLimitSession) SendMessage(ctx context.Context, message string, options ...Option) (string, error) {
	if rl.requestsPerMinute > 0 {
		err := rl.rateLimiter.Wait(ctx)
		if err != nil {
			return "", fmt.Errorf("rate limiter wait error: %w", err)
		}
	}

	response, err := rl.llmSession.SendMessage(ctx, message, options...)
	return response, err
}
