package agent

import (
	"github.com/EduardDranca/GoAgent/internal/agent/models"
	"github.com/EduardDranca/GoAgent/internal/agent/service"
	"github.com/EduardDranca/GoAgent/internal/input"
	"github.com/EduardDranca/GoAgent/internal/utils"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func NewLocalProgrammingAgentWithAutoCommit(programmingService service.ProgrammingService, gitUtil utils.GitUtil, autoCommit bool) AgentInterface[models.AgentRequest] {
	if gitUtil == nil {
		gitUtil = &utils.RealGitUtil{} // Default to RealGitUtil if nil is provided
	}
	return &LocalProgrammingAgent{programmingService: programmingService, gitUtil: gitUtil, autoCommit: autoCommit}
}

func TestLocalProgrammingAgent_Implement(t *testing.T) {
	// Create a temporary directory.
	tempDir, err := os.MkdirTemp("", "test-local-agent")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := service.NewMockService(ctrl) // Use generated mock
	mockGitUtil := utils.NewMockGitUtil(ctrl)   // Use generated mock

	// Create LocalProgrammingAgent with MockGitUtil
	agent := NewLocalProgrammingAgentWithAutoCommit(mockService, mockGitUtil, true)

	// Define test AgentRequest
	request := models.AgentRequest{
		Directory: tempDir,
		Query:     "test change request",
	}

	mockGitUtil.EXPECT().LsTree(tempDir).Return(make([]string, 0), nil).Times(1)                   // Expect LsTree call
	mockService.EXPECT().ImplementWithContext(gomock.Any()).Return("commit message", nil).Times(1) // Expect ImplementWithContext call
	mockGitUtil.EXPECT().Add(tempDir).Return(nil).Times(1)                                         // Expect Add call
	mockGitUtil.EXPECT().Commit(tempDir, "commit message").Return(nil).Times(1)                    // Expect Commit call

	// Call Implement method
	err = agent.Implement(request)
	require.NoError(t, err)
}

func testInputGetter(responses []string) func(string) (string, error) {
	inputIndex := 0
	return func(prompt string) (string, error) {
		if inputIndex < len(responses) {
			response := responses[inputIndex]
			inputIndex++
			return response, nil
		}
		return "", io.EOF // Simulate EOF if no more responses
	}
}

func TestLocalProgrammingAgent_Implement_CommitChoice_Y(t *testing.T) {
	// Create a temporary directory.
	tempDir, err := os.MkdirTemp("", "test-local-agent")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := service.NewMockService(ctrl)
	mockGitUtil := utils.NewMockGitUtil(ctrl)

	// Create LocalProgrammingAgent with MockGitUtil and autoCommit=false
	agent := NewLocalProgrammingAgentWithAutoCommit(mockService, mockGitUtil, false)

	// Mock user input for commit choice 'Y'
	input.UserInputGetter = testInputGetter([]string{"Y"})

	// Define test AgentRequest
	request := models.AgentRequest{
		Directory: tempDir,
		Query:     "test change request",
	}

	mockGitUtil.EXPECT().LsTree(tempDir).Return(make([]string, 0), nil).Times(1)
	mockService.EXPECT().ImplementWithContext(gomock.Any()).Return("commit message", nil).Times(1)
	mockGitUtil.EXPECT().Add(tempDir).Return(nil).Times(1)                      // Expect Add call 1 time for 'Y'
	mockGitUtil.EXPECT().Commit(tempDir, "commit message").Return(nil).Times(1) // Expect Commit call 1 time for 'Y'

	// Call Implement method
	err = agent.Implement(request)
	require.NoError(t, err)
}

func TestLocalProgrammingAgent_Implement_CommitChoice_N(t *testing.T) {
	// Create a temporary directory.
	tempDir, err := os.MkdirTemp("", "test-local-agent")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := service.NewMockService(ctrl)
	mockGitUtil := utils.NewMockGitUtil(ctrl)

	// Create LocalProgrammingAgent with MockGitUtil and autoCommit=false
	agent := NewLocalProgrammingAgentWithAutoCommit(mockService, mockGitUtil, false)

	// Mock user input for commit choice 'N'
	input.UserInputGetter = testInputGetter([]string{"N"})

	// Define test AgentRequest
	request := models.AgentRequest{
		Directory: tempDir,
		Query:     "test change request",
	}

	mockGitUtil.EXPECT().LsTree(tempDir).Return(make([]string, 0), nil).Times(1)
	mockService.EXPECT().ImplementWithContext(gomock.Any()).Return("commit message", nil).Times(1)
	mockGitUtil.EXPECT().Add(tempDir).Times(0)                  // Expect no Add call for 'N'
	mockGitUtil.EXPECT().Commit(tempDir, gomock.Any()).Times(0) // Expect no Commit call for 'N'

	// Call Implement method
	err = agent.Implement(request)
	require.NoError(t, err)
}

func TestLocalProgrammingAgent_Implement_CommitChoice_A(t *testing.T) {
	// Create a temporary directory.
	tempDir, err := os.MkdirTemp("", "test-local-agent")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := service.NewMockService(ctrl)
	mockGitUtil := utils.NewMockGitUtil(ctrl)

	// Create LocalProgrammingAgent with MockGitUtil and autoCommit=false
	agent := NewLocalProgrammingAgentWithAutoCommit(mockService, mockGitUtil, false)

	// Mock user input for commit choice 'A'
	input.UserInputGetter = testInputGetter([]string{"A"})

	// Define test AgentRequest
	request := models.AgentRequest{
		Directory: tempDir,
		Query:     "test change request",
	}

	mockGitUtil.EXPECT().LsTree(tempDir).Return(make([]string, 0), nil).Times(1)
	mockService.EXPECT().ImplementWithContext(gomock.Any()).Return("commit message", nil).Times(1)
	mockGitUtil.EXPECT().Add(tempDir).Return(nil).Times(1)                      // Expect Add call 1 time for 'A'
	mockGitUtil.EXPECT().Commit(tempDir, "commit message").Return(nil).Times(1) // Expect Commit call 1 time for 'A'

	// Call Implement method
	err = agent.Implement(request)
	require.NoError(t, err)
	require.True(t, agent.(*LocalProgrammingAgent).autoCommit, "autoCommit should be set to true") // Assert autoCommit is true
}
