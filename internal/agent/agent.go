package agent

// AgentInterface defines the interface for all types of agents.
// The Implement method takes a generic type T as a request parameter.
type AgentInterface[T any] interface {
	Implement(request T) error
	Ask(request T) (string, error)
}
