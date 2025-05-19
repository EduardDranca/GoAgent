package initialize

import (
	"context"
	"github.com/EduardDranca/GoAgent/internal/config"
	"testing"
)

func TestInitLLMService(t *testing.T) {
	ctx := context.Background()

	cfg := &config.Config{
		ProgrammingService: config.GeminiService,
		GeminiApiKey:       "valid-api-key",
	}

	_, err := InitProgrammingService(ctx, cfg)
	if err != nil {
		t.Errorf("InitLLMService returned an error: %v", err)
	}
}

func TestInitLLMServiceInvalidServiceType(t *testing.T) {
	ctx := context.Background()
	cfg := &config.Config{
		ProgrammingService: config.LLMServiceType("invalid-service-type"),
		GeminiApiKey:       "valid-api-key",
	}

	_, err := InitProgrammingService(ctx, cfg)
	if err == nil {
		t.Errorf("InitLLMService did not return an error for invalid service type")
	}
}

func TestInitLLMServiceEmptyAPIKey(t *testing.T) {
	ctx := context.Background()
	cfg := &config.Config{
		ProgrammingService: config.GeminiService,
		GeminiApiKey:       "",
	}

	_, err := InitProgrammingService(ctx, cfg)
	if err == nil {
		t.Errorf("InitLLMService did not return an error for empty API key")
	}
}
