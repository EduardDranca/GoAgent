package service

import (
	"github.com/EduardDranca/GoAgent/internal/agent/context"
)

// ProgrammingService interface for programming agent service.
type ProgrammingService interface {
	// ImplementWithContext implements the change request using the given context.
	ImplementWithContext(agentContext context.ProgrammingAgentContext) (string, error)
	AskWithContext(ctx context.ProgrammingAgentContext) (string, error)
}
