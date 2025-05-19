package config

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/EduardDranca/GoAgent/internal/logging"
	"gopkg.in/yaml.v3"
)

// Config holds all the configuration parameters for the application.
type Config struct {
	Directory          string
	ProgrammingService LLMServiceType
	GeminiApiKey       string
	GroqApiKey         string
	OpenaiApiKey       string
	// RateLimitRPM specifies the rate limit in requests per minute.
	// Defaults to 0, which means no rate limit.
	RateLimitRPM int
	// GlamourStylePath specifies the style to use for glamour output.
	GlamourStylePath GlamourStyleType
	LogLevel         string // Add LogLevel field
	MaxHistoryLength int    // Add MaxHistoryLength field
	MaxProcessLoops  int    `yaml:"max_process_loops"`

	InstructionsModelName string `yaml:"instructions_model"`
	GenerateCodeModelName string `yaml:"generate_code_model"`
	AnalysisModelName     string `yaml:"analysis_model"`
}

// ConfigFile is a struct for YAML parsing, mirroring Config but suitable for file loading.
type ConfigFile struct {
	// Default model names - these are defaults if not specified per service
	Gemini struct {
		InstructionsModelName string `yaml:"instructions_model"`
		GenerateCodeModelName string `yaml:"generate_code_model"`
		AnalysisModelName     string `yaml:"analysis_model"`
	} `yaml:"gemini"`
	Groq struct {
		InstructionsModelName string `yaml:"instructions_model"`
		GenerateCodeModelName string `yaml:"generate_code_model"`
		AnalysisModelName     string `yaml:"analysis_model"`
	} `yaml:"groq"`
	OpenAI struct {
		InstructionsModelName string `yaml:"instructions_model"`
		GenerateCodeModelName string `yaml:"generate_code_model"`
		AnalysisModelName     string `yaml:"analysis_model"`
	} `yaml:"openai"`
	MaxHistoryLength int    `yaml:"max_history_length"`
	MaxProcessLoops  int    `yaml:"max_process_loops"`
	RateLimitRPM     int    `yaml:"rate_limit_rpm"` // Add RateLimitRPM field for config file
	GlamourStylePath string `yaml:"glamour_style"`  // Add GlamourStylePath field for config file
	LogLevel         string `yaml:"log_level"`      // Add LogLevel field for config file
}

// LoadConfig parses command-line flags, loads environment variables, and reads config file.
func LoadConfig() (*Config, error) {
	defaultDir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get current working directory: %w", err)
	}

	// Define command-line flags with default values
	defaultProgrammingService := string(GeminiService)
	defaultRateLimitRPM := 0
	defaultGlamourStyle := string(DraculaStyle)
	defaultLogLevel := "info"
	defaultMaxHistoryLength := 100
	defaultMaxProcessLoops := 25 // Corrected default value based on analysis history

	directoryFlag := flag.String("directory", "", "Sets the root directory of your Git repository. Defaults to the current working directory if not provided. Must be a Git repository.")
	programmingServiceFlag := flag.String("service", defaultProgrammingService, fmt.Sprintf("Sets the programming service to use (%s, %s, %s). Defaults to %s.", GeminiService, GroqService, OpenAIService, defaultProgrammingService))
	geminiApiKeyFlag := flag.String("gemini-api-key", "", "Sets the Gemini API key. Required when using the Gemini service.")
	groqApiKeyFlag := flag.String("groq-api-key", "", "Sets the Groq API key. Required when using the Groq service.")
	openaiApiKeyFlag := flag.String("openai-api-key", "", "Sets the OpenAI API key. Required when using the OpenAI service.")
	rateLimitRPMFlag := flag.Int("rate-limit", defaultRateLimitRPM, "Sets the rate limit for API requests per minute. Prevents exceeding API usage limits. Defaults to 0, which means no rate limit.")
	glamourStyleFlag := flag.String("glamour-style", defaultGlamourStyle, fmt.Sprintf("Sets the Glamour style for Markdown rendering in the terminal. Defaults to %s.", defaultGlamourStyle))
	logLevelFlag := flag.String("log-level", defaultLogLevel, fmt.Sprintf("Sets the logging level. Allowed values are: %s. Defaults to %s.", strings.Join([]string{"debug", "info", "warning", "error"}, ", "), defaultLogLevel))
	maxHistoryLengthFlag := flag.Int("max-history-length", defaultMaxHistoryLength, "Sets the maximum history length for LLM sessions. Defaults to 100.")
	maxProcessLoopsFlag := flag.Int("max-process-loops", defaultMaxProcessLoops, "Sets the maximum number of process loops. Defaults to 25.")

	// Parse command-line flags
	flag.Parse()

	// Read flag values
	directory := defaultDir
	if *directoryFlag != "" {
		directory = *directoryFlag
	}
	programmingServiceStr := *programmingServiceFlag
	geminiApiKey := *geminiApiKeyFlag
	groqApiKey := *groqApiKeyFlag
	openaiApiKey := *openaiApiKeyFlag
	rateLimitRPM := *rateLimitRPMFlag
	glamourStylePathStr := *glamourStyleFlag
	logLevel := *logLevelFlag
	maxHistoryLength := *maxHistoryLengthFlag
	maxProcessLoops := *maxProcessLoopsFlag

	// Load API keys from environment variables if not provided via flags
	if geminiApiKey == "" {
		geminiApiKey = os.Getenv("GEMINI_API_KEY")
	}
	if groqApiKey == "" {
		groqApiKey = os.Getenv("GROQ_API_KEY")
	}
	if openaiApiKey == "" {
		openaiApiKey = os.Getenv("OPENAI_API_KEY")
	}

	// Validate programming service
	var programmingService LLMServiceType
	switch programmingServiceStr {
	case string(GeminiService):
		programmingService = GeminiService
	case string(GroqService):
		programmingService = GroqService
	case string(OpenAIService):
		programmingService = OpenAIService
	default:
		return nil, fmt.Errorf("invalid programming service: %s", programmingServiceStr)
	}

	logging.Logger.Infof("Service being used: %s", programmingService)

	// Check for API key if required service is selected
	if programmingService == GeminiService && geminiApiKey == "" {
		return nil, fmt.Errorf("gemini-api-key flag or GEMINI_API_KEY environment variable not set for gemini programming service")
	}
	if programmingService == GroqService && groqApiKey == "" {
		return nil, fmt.Errorf("groq-api-key flag or GROQ_API_KEY environment variable not set for groq programming service")
	}
	if programmingService == OpenAIService && openaiApiKey == "" {
		return nil, fmt.Errorf("openai-api-key flag or OPENAI_API_KEY environment variable not set for openai programming service")
	}

	// Initialize config with values from flags (or their defaults)
	cfg := &Config{
		Directory:          directory,
		ProgrammingService: programmingService,
		GeminiApiKey:       geminiApiKey,
		GroqApiKey:         groqApiKey,
		OpenaiApiKey:       openaiApiKey,
		RateLimitRPM:       rateLimitRPM,
		GlamourStylePath:   GlamourStyleType(glamourStylePathStr), // Will be validated later
		LogLevel:           logLevel,                              // Will be validated later
		MaxHistoryLength:   maxHistoryLength,
		MaxProcessLoops:    maxProcessLoops,

		// Default model names - these are defaults if not specified per service
		InstructionsModelName: "",
		GenerateCodeModelName: "",
		AnalysisModelName:     "",
	}

	configFilePath := filepath.Join(".go-agent", "config.yaml")

	_, err = os.Stat(configFilePath)

	defaultConfigFileMap := map[string]map[string]string{
		string(GeminiService): {
			"instructions_model":  "gemini-2.5-flash-preview-04-17",
			"generate_code_model": "gemini-2.5-flash-preview-04-17",
			"analysis_model":      "gemini-2.5-flash-preview-04-17",
		},
		string(GroqService): {
			"instructions_model":  "meta-llama/llama-4-maverick-17b-128e-instruct",
			"generate_code_model": "meta-llama/llama-4-maverick-17b-128e-instruct",
			"analysis_model":      "meta-llama/llama-4-maverick-17b-128e-instruct",
		},
		string(OpenAIService): {
			"instructions_model":  "gpt-4.5-preview",
			"generate_code_model": "gpt-4.5-preview",
			"analysis_model":      "gpt-4.5-preview",
		},
	}

	if os.IsNotExist(err) {
		// Config file does not exist, create directory and file with default values
		err = os.MkdirAll(filepath.Dir(configFilePath), 0755)
		if err != nil {
			return nil, fmt.Errorf("failed to create config directory: %w", err)
		}
		file, err := os.Create(configFilePath)
		if err != nil {
			return nil, fmt.Errorf("failed to create config file: %w", err)
		}
		defer file.Close()

		// Create a temporary ConfigFile struct with default values for writing
		tempDefaultConfigFile := &ConfigFile{
			Gemini: struct {
				InstructionsModelName string `yaml:"instructions_model"`
				GenerateCodeModelName string `yaml:"generate_code_model"`
				AnalysisModelName     string `yaml:"analysis_model"`
			}{
				InstructionsModelName: defaultConfigFileMap[string(GeminiService)]["instructions_model"],
				GenerateCodeModelName: defaultConfigFileMap[string(GeminiService)]["generate_code_model"],
				AnalysisModelName:     defaultConfigFileMap[string(GeminiService)]["analysis_model"],
			},
			Groq: struct {
				InstructionsModelName string `yaml:"instructions_model"`
				GenerateCodeModelName string `yaml:"generate_code_model"`
				AnalysisModelName     string `yaml:"analysis_model"`
			}{
				InstructionsModelName: defaultConfigFileMap[string(GroqService)]["instructions_model"],
				GenerateCodeModelName: defaultConfigFileMap[string(GroqService)]["generate_code_model"],
				AnalysisModelName:     defaultConfigFileMap[string(GroqService)]["analysis_model"],
			},
			OpenAI: struct {
				InstructionsModelName string `yaml:"instructions_model"`
				GenerateCodeModelName string `yaml:"generate_code_model"`
				AnalysisModelName     string `yaml:"analysis_model"`
			}{
				InstructionsModelName: defaultConfigFileMap[string(OpenAIService)]["instructions_model"],
				GenerateCodeModelName: defaultConfigFileMap[string(OpenAIService)]["generate_code_model"],
				AnalysisModelName:     defaultConfigFileMap[string(OpenAIService)]["analysis_model"],
			},
			MaxHistoryLength: defaultMaxHistoryLength,
			MaxProcessLoops:  defaultMaxProcessLoops,
			RateLimitRPM:     defaultRateLimitRPM,
			GlamourStylePath: defaultGlamourStyle,
			LogLevel:         defaultLogLevel,
		}

		yamlData, err := yaml.Marshal(tempDefaultConfigFile) // Use the temporary struct
		if err != nil {
			return nil, fmt.Errorf("failed to marshal default config to YAML: %w", err)
		}
		_, err = file.Write(yamlData)
		if err != nil {
			return nil, fmt.Errorf("failed to write default config to file: %w", err)
		}

	} else if err == nil {
		// Config file exists, load values from it
		yamlFile, err := os.ReadFile(configFilePath)
		if err != nil {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}

		var configFile ConfigFile
		err = yaml.Unmarshal(yamlFile, &configFile)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal config file: %w", err)
		}

		// Override values from config file if they are set AND the corresponding flag was not explicitly set (i.e., still has its default value)

		// Model names
		if programmingService == GroqService {
			if configFile.Groq.InstructionsModelName != "" {
				cfg.InstructionsModelName = configFile.Groq.InstructionsModelName
			}
			if configFile.Groq.GenerateCodeModelName != "" {
				cfg.GenerateCodeModelName = configFile.Groq.GenerateCodeModelName
			}
			if configFile.Groq.AnalysisModelName != "" {
				cfg.AnalysisModelName = configFile.Groq.AnalysisModelName
			}
		}
		if programmingService == GeminiService {
			if configFile.Gemini.InstructionsModelName != "" {
				cfg.InstructionsModelName = configFile.Gemini.InstructionsModelName
			}
			if configFile.Gemini.GenerateCodeModelName != "" {
				cfg.GenerateCodeModelName = configFile.Gemini.GenerateCodeModelName
			}
			if configFile.Gemini.AnalysisModelName != "" {
				cfg.AnalysisModelName = configFile.Gemini.AnalysisModelName
			}
		}
		if programmingService == OpenAIService {
			if configFile.OpenAI.InstructionsModelName != "" {
				cfg.InstructionsModelName = configFile.OpenAI.InstructionsModelName
			}
			if configFile.OpenAI.GenerateCodeModelName != "" { // Corrected field name
				cfg.GenerateCodeModelName = configFile.OpenAI.GenerateCodeModelName
			}
			if configFile.OpenAI.AnalysisModelName != "" {
				cfg.AnalysisModelName = configFile.OpenAI.AnalysisModelName
			}
		}

		// MaxHistoryLength: Override if set in file AND flag is default
		if configFile.MaxHistoryLength != 0 && *maxHistoryLengthFlag == defaultMaxHistoryLength {
			cfg.MaxHistoryLength = configFile.MaxHistoryLength
		}

		// MaxProcessLoops: Override if set in file AND flag is default
		if configFile.MaxProcessLoops != 0 && *maxProcessLoopsFlag == defaultMaxProcessLoops {
			cfg.MaxProcessLoops = configFile.MaxProcessLoops
		}

		// RateLimitRPM: Override if set in file (> 0) AND flag is default (0)
		if configFile.RateLimitRPM > 0 && *rateLimitRPMFlag == defaultRateLimitRPM {
			cfg.RateLimitRPM = configFile.RateLimitRPM
		}

		// GlamourStylePath: Override if set in file AND flag is default
		if configFile.GlamourStylePath != "" && *glamourStyleFlag == defaultGlamourStyle {
			cfg.GlamourStylePath = GlamourStyleType(configFile.GlamourStylePath)
		}

		// LogLevel: Override if set in file AND flag is default
		if configFile.LogLevel != "" && *logLevelFlag == defaultLogLevel {
			cfg.LogLevel = configFile.LogLevel
		}
	}

	// Set default model names if still empty after checking file
	if cfg.InstructionsModelName == "" {
		cfg.InstructionsModelName = defaultConfigFileMap[string(programmingService)]["instructions_model"]
	}
	if cfg.GenerateCodeModelName == "" {
		cfg.GenerateCodeModelName = defaultConfigFileMap[string(programmingService)]["generate_code_model"]
	}
	if cfg.AnalysisModelName == "" {
		cfg.AnalysisModelName = defaultConfigFileMap[string(programmingService)]["analysis_model"]
	}

	// Validate Glamour style after potential override
	var finalGlamourStyle GlamourStyleType
	switch cfg.GlamourStylePath {
	case AsciiStyle, AutoStyle, DarkStyle, DraculaStyle, TokyoNightStyle, LightStyle, NottyStyle, PinkStyle:
		finalGlamourStyle = cfg.GlamourStylePath
	default:
		return nil, fmt.Errorf("invalid glamour style: %s, allowed styles are %s, %s, %s, %s, %s, %s, %s, %s", cfg.GlamourStylePath, AsciiStyle, AutoStyle, DarkStyle, DraculaStyle, TokyoNightStyle, LightStyle, NottyStyle, PinkStyle)
	}
	cfg.GlamourStylePath = finalGlamourStyle // Assign validated style

	// Validate Log Level after potential override
	switch cfg.LogLevel {
	case "debug", "info", "warning", "error":
		// Valid log level
	default:
		return nil, fmt.Errorf("invalid log level: %s, allowed levels are debug, info, warning, error", cfg.LogLevel)
	}

	return cfg, nil
}
