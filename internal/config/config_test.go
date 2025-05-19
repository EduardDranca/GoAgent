package config_test

import (
	"flag"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/EduardDranca/GoAgent/internal/config"
)

func TestLoadConfig_DefaultValues(t *testing.T) {
	_, err := config.LoadConfig()
	require.Error(t, err, "LoadConfig should return an error if no Gemini API key is provided for default Gemini service")
}

func TestLoadConfig_InvalidServiceFlag(t *testing.T) {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	os.Args = []string{"agent", "-service", "invalid"}
	defer func() {
		os.Args = os.Args[:1]
	}()
	_, err := config.LoadConfig()
	require.Error(t, err, "LoadConfig should return an error for invalid service")
}

func TestLoadConfig_GeminiApiKeyFromEnv(t *testing.T) {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	os.Args = []string{"agent", "-service", "gemini"}
	os.Setenv("GEMINI_API_KEY", "test-gemini-key")
	defer func() {
		os.Unsetenv("GEMINI_API_KEY")
		os.Args = os.Args[:1]
	}()

	cfg, err := config.LoadConfig()
	require.NoError(t, err, "LoadConfig should not return an error")
	require.Equal(t, "test-gemini-key", cfg.GeminiApiKey, "Gemini API key should be loaded from env")
	require.Equal(t, config.GeminiService, cfg.ProgrammingService, "ProgrammingService should be Gemini")
}

func TestLoadConfig_GeminiApiKeyFromFlag(t *testing.T) {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	os.Args = []string{"agent", "-service", "gemini", "-gemini-api-key", "test-gemini-key"}
	defer func() {
		os.Args = os.Args[:1]
	}()

	cfg, err := config.LoadConfig()
	require.NoError(t, err, "LoadConfig should not return an error")
	require.Equal(t, "test-gemini-key", cfg.GeminiApiKey, "Gemini API key should be loaded from flag")
	require.Equal(t, config.GeminiService, cfg.ProgrammingService, "ProgrammingService should be Gemini")
}

func TestLoadConfig_GeminiApiKeyMissing(t *testing.T) {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	os.Args = []string{"agent", "-service", "gemini"}
	defer func() {
		os.Args = os.Args[:1]
	}()
	_, err := config.LoadConfig()
	require.Error(t, err, "LoadConfig should return an error if Gemini API key is missing")
}

func TestLoadConfig_GroqApiKeyMissing(t *testing.T) {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	os.Args = []string{"agent", "-service", "groq"}
	defer func() {
		os.Args = os.Args[:1]
	}()
	_, err := config.LoadConfig()
	require.Error(t, err, "LoadConfig should return an error if Groq API key is missing")
}

func TestLoadConfig_OpenaiApiKeyMissing(t *testing.T) {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	os.Args = []string{"agent", "-service", "openai"}
	defer func() {
		os.Args = os.Args[:1]
	}()
	_, err := config.LoadConfig()
	require.Error(t, err, "LoadConfig should return an error if OpenAI API key is missing")
}

func TestLoadConfig_GroqApiKeyFromEnv(t *testing.T) {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	os.Args = []string{"agent", "-service", "groq"}
	os.Setenv("GROQ_API_KEY", "test-groq-key")
	defer func() {
		os.Unsetenv("GROQ_API_KEY")
		os.Args = os.Args[:1]
	}()

	cfg, err := config.LoadConfig()
	require.NoError(t, err, "LoadConfig should not return an error")
	require.Equal(t, "test-groq-key", cfg.GroqApiKey, "Groq API key should be loaded from env")
	require.Equal(t, config.GroqService, cfg.ProgrammingService, "ProgrammingService should be Groq")
}

func TestLoadConfig_OpenaiApiKeyFromEnv(t *testing.T) {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	os.Args = []string{"agent", "-service", "openai"}
	os.Setenv("OPENAI_API_KEY", "test-openai-key")
	defer func() {
		os.Unsetenv("OPENAI_API_KEY")
		os.Args = os.Args[:1]
	}()

	cfg, err := config.LoadConfig()
	require.NoError(t, err, "LoadConfig should not return an error")
	require.Equal(t, "test-openai-key", cfg.OpenaiApiKey, "OpenAI API key should be loaded from env")
	require.Equal(t, config.OpenAIService, cfg.ProgrammingService, "ProgrammingService should be OpenAI")
}

func TestLoadConfig_DefaultRateLimitRPM(t *testing.T) {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	os.Setenv("GEMINI_API_KEY", "test-gemini-key")
	os.Args = []string{"agent"}
	defer func() {
		os.Args = os.Args[:1]
	}()

	cfg, _ := config.LoadConfig()
	require.Equal(t, 0, cfg.RateLimitRPM, "RateLimitRPM should default to 0")
}

func TestLoadConfig_RateLimitRPMFromFlag(t *testing.T) {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	os.Args = []string{"agent", "-rate-limit", "100"}
	os.Setenv("GEMINI_API_KEY", "test-gemini-key")
	defer func() {
		os.Args = os.Args[:1]
	}()

	cfg, _ := config.LoadConfig()
	require.Equal(t, 100, cfg.RateLimitRPM, "RateLimitRPM should be loaded from flag")
}
