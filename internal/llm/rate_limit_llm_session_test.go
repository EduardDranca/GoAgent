package llm

import (
	"context"
	"errors"
	"github.com/EduardDranca/GoAgent/internal/agent/models"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"golang.org/x/time/rate"
)

func TestRateLimitSession_SendMessage(t *testing.T) {
	mockSession := NewMockLLMSession("mock response", nil)
	limiter := rate.NewLimiter(rate.Limit(10), 10)
	rlSession := NewRateLimitSession(mockSession, limiter)

	msg := "test message"
	resp, err := rlSession.SendMessage(context.Background(), msg)

	assert.NoError(t, err)
	assert.Equal(t, "mock response", resp)
	assert.Equal(t, []models.Message{{Content: msg}}, mockSession.History)
}

func TestRateLimitSession_SendMessage_RateLimited(t *testing.T) {
	mockSession := NewMockLLMSession("mock response", nil)
	limiter := rate.NewLimiter(rate.Limit(1), 1) // Limit to 1 request per second
	rlSession := NewRateLimitSession(mockSession, limiter)

	// First request should be allowed
	_, err := rlSession.SendMessage(context.Background(), "message 1")
	assert.NoError(t, err)

	// Second request should be rate limited
	start := time.Now()
	_, err = rlSession.SendMessage(context.Background(), "message 2")
	duration := time.Since(start)

	assert.NoError(t, err)                          // No error from rate limiter itself, the wait should just delay
	assert.GreaterOrEqual(t, duration, time.Second) // Check if it waited at least 1 second (due to rate limit of 1 per second)
}

func TestRateLimitSession_GetHistory_SetHistory(t *testing.T) {
	mockSession := NewMockLLMSession("mock response", nil)
	rlSession := NewRateLimitSession(mockSession, rate.NewLimiter(rate.Limit(10), 10))

	history := []models.Message{{Content: "message 1"}, {Content: "message 2"}}
	rlSession.SetHistory(history)
	retrievedHistory := rlSession.GetHistory()

	assert.Equal(t, history, retrievedHistory)
	assert.Equal(t, history, mockSession.History) // Ensure history is passed to mock session
}

func TestRateLimitSession_RateLimit(t *testing.T) {
	mockSession := NewMockLLMSession("mock response", nil)
	rateLimiter := rate.NewLimiter(rate.Limit(1), 1) // 1 request per second
	rlSession := NewRateLimitSession(mockSession, rateLimiter)

	numMessages := 3
	start := time.Now()
	for i := 0; i < numMessages; i++ {
		_, _ = rlSession.SendMessage(context.Background(), "message")
	}
	elapsed := time.Since(start)

	expectedDuration := time.Duration(numMessages-1) * time.Second // Expected time is (numMessages - 1) seconds due to rate limit
	tolerance := 500 * time.Millisecond                            // Allow some tolerance

	assert.InDelta(t, expectedDuration.Seconds(), elapsed.Seconds(), tolerance.Seconds(), "Elapsed time should be approximately equal to expected time")
}

func TestRateLimitSession_ContextCancel(t *testing.T) {
	mockSession := NewMockLLMSession("mock response", nil)
	rateLimiter := rate.NewLimiter(rate.Limit(10), 10)
	rlSession := NewRateLimitSession(mockSession, rateLimiter)

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Immediately cancel the context

	_, err := rlSession.SendMessage(ctx, "test message")

	assert.Error(t, err)
	assert.True(t, errors.Is(err, context.Canceled), "Expected context cancellation error")
}

func TestRateLimitSession_PropagatesLLMSessionError(t *testing.T) {
	mockError := errors.New("mock LLMSession error")
	mockSession := NewMockLLMSession("", mockError) // Mock session returns an error
	rateLimiter := rate.NewLimiter(rate.Limit(10), 10)
	rlSession := NewRateLimitSession(mockSession, rateLimiter)

	_, err := rlSession.SendMessage(context.Background(), "test message")

	assert.Error(t, err)
	assert.EqualError(t, err, mockError.Error(), "Expected error to be propagated from mock LLMSession")
}
