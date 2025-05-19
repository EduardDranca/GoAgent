package llm

import (
	"context"
	"github.com/EduardDranca/GoAgent/internal/agent/models"
)

// MockLLMSession is a mock implementation of the LLMSession interface.
type MockLLMSession struct {
	SendMessageReturnValue string
	SendMessageError       error
	History                []models.Message
}

func (m *MockLLMSession) SendMessage(ctx context.Context, message string, options ...Option) (string, error) {
	// Append the message to the history
	m.History = append(m.History, models.Message{Content: message})
	// Return the stored return value and error
	return m.SendMessageReturnValue, m.SendMessageError
}

func (m *MockLLMSession) GetHistory() []models.Message {
	return m.History
}

func (m *MockLLMSession) SetHistory(history []models.Message) {
	m.History = history
}

// NewMockLLMSession is a constructor for MockLLMSession.
func NewMockLLMSession(sendMessageReturnValue string, sendMessageError error) *MockLLMSession {
	return &MockLLMSession{
		SendMessageReturnValue: sendMessageReturnValue,
		SendMessageError:       sendMessageError,
		History:                []models.Message{},
	}
}
