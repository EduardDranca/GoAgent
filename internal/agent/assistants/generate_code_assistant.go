package assistants

import (
	context2 "context"
	"github.com/EduardDranca/GoAgent/internal/agent/models"
	"github.com/EduardDranca/GoAgent/internal/llm"
	"github.com/EduardDranca/GoAgent/internal/utils"
)

type GenerateCodeAssistant interface {
	GenerateCode(ctx context2.Context, message string) (string, error)
}

type defaultGenerateCodeAssistant struct {
	session llm.LLMSession
}

func NewGenerateCodeAssistant(session llm.LLMSession) GenerateCodeAssistant {
	return &defaultGenerateCodeAssistant{
		session: session,
	}
}

func (a *defaultGenerateCodeAssistant) GenerateCode(ctx context2.Context, message string) (string, error) {
	response, err := a.session.SendMessage(ctx, message)
	a.session.SetHistory([]models.Message{})
	if err != nil {
		return "", err
	}

	extractedContent, err := utils.ExtractCodeBlock(response)
	if err != nil {
		return response, nil
	}

	return extractedContent, nil
}
