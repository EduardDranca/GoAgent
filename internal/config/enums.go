package config

// LLMServiceType represents the type of LLM service to use.
type LLMServiceType string

const (
	GeminiService  LLMServiceType = "gemini"
	GroqService   LLMServiceType = "groq"
	OpenAIService LLMServiceType = "openai"
)

// GlamourStyleType represents the type of Glamour style to use.
type GlamourStyleType string

const (
	AsciiStyle      GlamourStyleType = "ascii"
	AutoStyle       GlamourStyleType = "auto"
	DarkStyle       GlamourStyleType = "dark"
	DraculaStyle    GlamourStyleType = "dracula"
	TokyoNightStyle GlamourStyleType = "tokyo-night"
	LightStyle      GlamourStyleType = "light"
	NottyStyle      GlamourStyleType = "notty"
	PinkStyle       GlamourStyleType = "pink"
)
