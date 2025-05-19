package models

type AgentRequest struct {
	Query     string `json:"query"`
	Directory string `json:"directory"`
}
