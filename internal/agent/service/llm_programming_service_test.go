package service

import (
	"context"
	"github.com/EduardDranca/GoAgent/internal/agent/commands"
	context2 "github.com/EduardDranca/GoAgent/internal/agent/context"
	"reflect"
	"testing"

	"github.com/EduardDranca/GoAgent/internal/agent/assistants"
	"github.com/EduardDranca/GoAgent/internal/llm"
	"go.uber.org/mock/gomock"
)

// MockAnalysisAssistant is a mock of AnalysisAssistant interface.
type MockAnalysisAssistant struct {
	ctrl     *gomock.Controller
	recorder *MockAnalysisAssistantMockRecorder
}

// MockAnalysisAssistantMockRecorder is the mock recorder for MockAnalysisAssistant.
type MockAnalysisAssistantMockRecorder struct {
	mock *MockAnalysisAssistant
}

// NewMockAnalysisAssistant creates a new mock instance.
func NewMockAnalysisAssistant(ctrl *gomock.Controller) *MockAnalysisAssistant {
	mock := &MockAnalysisAssistant{ctrl: ctrl}
	mock.recorder = &MockAnalysisAssistantMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAnalysisAssistant) EXPECT() *MockAnalysisAssistantMockRecorder {
	return m.recorder
}

// Execute mocks base method.
func (m *MockAnalysisAssistant) Execute(arg0 context.Context, arg1 string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Execute", arg0, arg1)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Execute indicates an expected call of Execute.
func (mr *MockAnalysisAssistantMockRecorder) Execute(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Execute", reflect.TypeOf((*MockAnalysisAssistant)(nil).Execute), arg0, arg1)
}

// ClearHistory mocks base method.
func (m *MockAnalysisAssistant) ClearHistory() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ClearHistory")
}

// ClearHistory indicates an expected call of ClearHistory.
func (mr *MockAnalysisAssistantMockRecorder) ClearHistory() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ClearHistory", reflect.TypeOf((*MockAnalysisAssistant)(nil).ClearHistory))
}

// MockInstructionAssistant is a mock of InstructionAssistant interface.
type MockInstructionAssistant struct {
	ctrl     *gomock.Controller
	recorder *MockInstructionAssistantMockRecorder
}

// MockInstructionAssistantMockRecorder is the mock recorder for MockInstructionAssistant.
type MockInstructionAssistantMockRecorder struct {
	mock *MockInstructionAssistant
}

// NewMockInstructionAssistant creates a new mock instance.
func NewMockInstructionAssistant(ctrl *gomock.Controller) *MockInstructionAssistant {
	mock := &MockInstructionAssistant{ctrl: ctrl}
	mock.recorder = &MockInstructionAssistantMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockInstructionAssistant) EXPECT() *MockInstructionAssistantMockRecorder {
	return m.recorder
}

// Execute mocks base method.
func (m *MockInstructionAssistant) Instruct(arg0 context.Context, arg1 string) (commands.Command, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Execute", arg0, arg1)
	ret0, _ := ret[0].(commands.Command)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Execute indicates an expected call of Execute.
func (mr *MockInstructionAssistantMockRecorder) Execute(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Execute", reflect.TypeOf((*MockInstructionAssistant)(nil).Instruct), arg0, arg1)
}

// ClearHistory mocks base method.
func (m *MockInstructionAssistant) ClearHistory() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ClearHistory")
}

// ClearHistory indicates an expected call of ClearHistory.
func (mr *MockInstructionAssistantMockRecorder) ClearHistory() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ClearHistory", reflect.TypeOf((*MockInstructionAssistant)(nil).ClearHistory))
}

// MockGenerateCodeAssistant is a mock of GenerateCodeAssistant interface.
type MockGenerateCodeAssistant struct {
	ctrl     *gomock.Controller
	recorder *MockGenerateCodeAssistantMockRecorder
}

// MockGenerateCodeAssistantMockRecorder is the mock recorder for MockGenerateCodeAssistant.
type MockGenerateCodeAssistantMockRecorder struct {
	mock *MockGenerateCodeAssistant
}

// NewMockGenerateCodeAssistant creates a new mock instance.
func NewMockGenerateCodeAssistant(ctrl *gomock.Controller) *MockGenerateCodeAssistant {
	mock := &MockGenerateCodeAssistant{ctrl: ctrl}
	mock.recorder = &MockGenerateCodeAssistantMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockGenerateCodeAssistant) EXPECT() *MockGenerateCodeAssistantMockRecorder {
	return m.recorder
}

// Execute mocks base method.
func (m *MockGenerateCodeAssistant) GenerateCode(arg0 context.Context, arg1 string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Execute", arg0, arg1)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Execute indicates an expected call of Execute.
func (mr *MockGenerateCodeAssistantMockRecorder) Execute(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Execute", reflect.TypeOf((*MockGenerateCodeAssistant)(nil).GenerateCode), arg0, arg1)
}

// ClearHistory mocks base method.
func (m *MockGenerateCodeAssistant) ClearHistory() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ClearHistory")
}

// ClearHistory indicates an expected call of ClearHistory.
func (mr *MockGenerateCodeAssistantMockRecorder) ClearHistory() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ClearHistory", reflect.TypeOf((*MockGenerateCodeAssistant)(nil).ClearHistory))
}

func TestNewLLMProgrammingService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	instructionsSessionMock := llm.NewMockLLMSession("", nil)
	generateCodeSessionMock := llm.NewMockLLMSession("", nil)
	analysisSessionMock := llm.NewMockLLMSession("", nil)
	askAnalysisSessionMock := llm.NewMockLLMSession("", nil)
	askInstructionSessionMock := llm.NewMockLLMSession("", nil)
	patchApplySessionMock := llm.NewMockLLMSession("", nil)

	codeAnalysisAssistant := assistants.NewAnalysisAssistant(analysisSessionMock)
	askAnalysisAssistant := assistants.NewAnalysisAssistant(askAnalysisSessionMock)
	codeInstructionAssistant := assistants.NewInstructionAssistant(instructionsSessionMock)
	askInstructionAssistant := assistants.NewInstructionAssistant(askInstructionSessionMock)
	codeGenerateCodeAssistant := assistants.NewGenerateCodeAssistant(generateCodeSessionMock)
	patchGenerateCodeAssistant := assistants.NewGenerateCodeAssistant(patchApplySessionMock)

	service := NewLLMProgrammingService(
		codeAnalysisAssistant,
		askAnalysisAssistant,
		codeInstructionAssistant,
		askInstructionAssistant,
		codeGenerateCodeAssistant,
		patchGenerateCodeAssistant,
		10,
	)

	if service == nil {
		t.Errorf("NewLLMProgrammingService returned nil")
	}
}

func TestLLMProgrammingService_ImplementWithContext_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAnalysisAssistant := NewMockAnalysisAssistant(ctrl)
	mockAskAnalysisAssistant := NewMockAnalysisAssistant(ctrl)
	mockInstructionAssistant := NewMockInstructionAssistant(ctrl)
	mockAskInstructionAssistant := NewMockInstructionAssistant(ctrl)
	mockGenerateCodeAssistant := NewMockGenerateCodeAssistant(ctrl)
	mockPatchGenerateCodeAssistant := NewMockGenerateCodeAssistant(ctrl)
	mockContext := context2.NewMockProgrammingAgentContext(ctrl)
	mockCommand := commands.NewMockCommand(ctrl)
	commitCommand := &commands.CommitCommand{Message: "Commit message"}

	// Set expectations for mock agents and context
	mockAnalysisAssistant.EXPECT().Execute(gomock.Any(), gomock.Any()).Return("context_files: []", nil).Times(1)
	mockInstructionAssistant.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(mockCommand, nil).Times(1)
	mockCommand.EXPECT().Process(mockContext).Return("File updated successfully", nil).Times(1)
	mockAnalysisAssistant.EXPECT().Execute(gomock.Any(), gomock.Any()).Return("Commit message: Commit message", nil).Times(1)
	mockInstructionAssistant.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(commitCommand, nil).Times(1)
	mockInstructionAssistant.EXPECT().ClearHistory().Times(1)

	mockContext.EXPECT().GetRepoStructure().Return([]string{"/"}).Times(1)
	mockContext.EXPECT().GetChangeRequest().Return("Implement feature X").Times(1)
	mockContext.EXPECT().GetFileContent(gomock.Any()).Return("", false).AnyTimes() // Assuming no context files for now
	mockContext.EXPECT().UpdateFileContent(gomock.Any(), gomock.Any()).Return().AnyTimes()

	service := NewLLMProgrammingService(
		mockAnalysisAssistant,
		mockAskAnalysisAssistant,
		mockInstructionAssistant,
		mockAskInstructionAssistant,
		mockGenerateCodeAssistant,
		mockPatchGenerateCodeAssistant,
		10,
	)

	response, err := service.ImplementWithContext(mockContext)

	if err != nil {
		t.Errorf("ImplementWithContext returned an error: %v", err)
	}
	if response != "Commit message" {
		t.Errorf("ImplementWithContext returned unexpected response: %v, want: %v", response, "File updated successfully")
	}
}
