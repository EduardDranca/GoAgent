package assistants

import (
	context2 "context"
	"github.com/EduardDranca/GoAgent/internal/llm"
)

type AnalysisAssistant interface {
	Execute(ctx context2.Context, message string) (string, error)
}

type defaultAnalysisAssistant struct {
	session llm.LLMSession
}

func NewAnalysisAssistant(session llm.LLMSession) AnalysisAssistant {
	return &defaultAnalysisAssistant{
		session: session,
	}
}

func (a *defaultAnalysisAssistant) Execute(ctx context2.Context, message string) (string, error) {
	response, err := a.session.SendMessage(ctx, message)
	if err != nil {
		return "", err
	}
	return response, nil
}
