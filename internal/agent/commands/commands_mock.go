// Code generated by MockGen. DO NOT EDIT.
// Source: internal/agent/commands/commands.go

// Package commands is a generated GoMock package.
package commands

import (
	context "github.com/EduardDranca/GoAgent/internal/agent/context"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockCommand is a mock of Command interface.
type MockCommand struct {
	ctrl     *gomock.Controller
	recorder *MockCommandMockRecorder
}

// MockCommandMockRecorder is the mock recorder for MockCommand.
type MockCommandMockRecorder struct {
	mock *MockCommand
}

// NewMockCommand creates a new mock instance.
func NewMockCommand(ctrl *gomock.Controller) *MockCommand {
	mock := &MockCommand{ctrl: ctrl}
	mock.recorder = &MockCommandMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCommand) EXPECT() *MockCommandMockRecorder {
	return m.recorder
}

// Process mocks base method.
func (m *MockCommand) Process(agentContext context.ProgrammingAgentContext) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Process", agentContext)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Process indicates an expected call of Process.
func (mr *MockCommandMockRecorder) Process(agentContext interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Process", reflect.TypeOf((*MockCommand)(nil).Process), agentContext)
}
