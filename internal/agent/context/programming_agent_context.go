package context

// ProgrammingAgentContext defines the interface for interacting with the project's context.
type ProgrammingAgentContext interface {
	GetFileContent(filePath string) (string, bool)
	UpdateFileContent(filePath string, newContents string)
	SearchCode(query string) map[string][]int
	GetRepoStructure() []string
	GetChangeRequest() string
	Delete(filePath string) error
	FlushChanges() error
	MoveFile(oldPath string, newPath string) error
}
