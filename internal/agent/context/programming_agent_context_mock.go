package context

import (
	"github.com/EduardDranca/GoAgent/internal/utils"
	"reflect"

	"go.uber.org/mock/gomock"
)

// MockProgrammingAgentContext is a mock of ProgrammingAgentContext interface.
type MockProgrammingAgentContext struct {
	ctrl     *gomock.Controller
	recorder *MockProgrammingAgentContextMockRecorder
}

// MockProgrammingAgentContextMockRecorder is the mock recorder for MockProgrammingAgentContext.
type MockProgrammingAgentContextMockRecorder struct {
	mock *MockProgrammingAgentContext
}

// NewMockProgrammingAgentContext creates a new mock instance.
func NewMockProgrammingAgentContext(ctrl *gomock.Controller) *MockProgrammingAgentContext {
	mock := &MockProgrammingAgentContext{ctrl: ctrl}
	mock.recorder = &MockProgrammingAgentContextMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockProgrammingAgentContext) EXPECT() *MockProgrammingAgentContextMockRecorder {
	return m.recorder
}

// Delete mocks base method.
func (m *MockProgrammingAgentContext) Delete(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockProgrammingAgentContextMockRecorder) Delete(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockProgrammingAgentContext)(nil).Delete), arg0)
}

// FlushChanges mocks base method.
func (m *MockProgrammingAgentContext) FlushChanges() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FlushChanges")
	ret0, _ := ret[0].(error)
	return ret0
}

// FlushChanges indicates an expected call of FlushChanges.
func (mr *MockProgrammingAgentContextMockRecorder) FlushChanges(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FlushChanges", reflect.TypeOf((*MockProgrammingAgentContext)(nil).FlushChanges), arg0)
}

// GetChangeRequest mocks base method.
func (m *MockProgrammingAgentContext) GetChangeRequest() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetChangeRequest")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetChangeRequest indicates an expected call of GetChangeRequest.
func (mr *MockProgrammingAgentContextMockRecorder) GetChangeRequest() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetChangeRequest", reflect.TypeOf((*MockProgrammingAgentContext)(nil).GetChangeRequest))
}

// GetFileContent mocks base method.
func (m *MockProgrammingAgentContext) GetFileContent(arg0 string) (string, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFileContent", arg0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// GetFileContent indicates an expected call of GetFileContent.
func (mr *MockProgrammingAgentContextMockRecorder) GetFileContent(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFileContent", reflect.TypeOf((*MockProgrammingAgentContext)(nil).GetFileContent), arg0)
}

// GetGitUtil mocks base method.
func (m *MockProgrammingAgentContext) GetGitUtil() utils.GitUtil {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetGitUtil")
	ret0, _ := ret[0].(utils.GitUtil)
	return ret0
}

// GetGitUtil indicates an expected call of GetGitUtil.
func (mr *MockProgrammingAgentContextMockRecorder) GetGitUtil() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGitUtil", reflect.TypeOf((*MockProgrammingAgentContext)(nil).GetGitUtil))
}

// GetRepoStructure mocks base method.
func (m *MockProgrammingAgentContext) GetRepoStructure() []string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRepoStructure")
	ret0, _ := ret[0].([]string)
	return ret0
}

// GetRepoStructure indicates an expected call of GetRepoStructure.
func (mr *MockProgrammingAgentContextMockRecorder) GetRepoStructure() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRepoStructure", reflect.TypeOf((*MockProgrammingAgentContext)(nil).GetRepoStructure))
}

// MoveFile mocks base method.
func (m *MockProgrammingAgentContext) MoveFile(arg0, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MoveFile", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// MoveFile indicates an expected call of MoveFile.
func (mr *MockProgrammingAgentContextMockRecorder) MoveFile(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MoveFile", reflect.TypeOf((*MockProgrammingAgentContext)(nil).MoveFile), arg0, arg1)
}

// SearchCode mocks base method.
func (m *MockProgrammingAgentContext) SearchCode(arg0 string) map[string][]int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SearchCode", arg0)
	ret0, _ := ret[0].(map[string][]int)
	return ret0
}

// SearchCode indicates an expected call of SearchCode.
func (mr *MockProgrammingAgentContextMockRecorder) SearchCode(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SearchCode", reflect.TypeOf((*MockProgrammingAgentContext)(nil).SearchCode), arg0)
}

// UpdateFileContent mocks base method.
func (m *MockProgrammingAgentContext) UpdateFileContent(arg0, arg1 string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "UpdateFileContent", arg0, arg1)
}

// UpdateFileContent indicates an expected call of UpdateFileContent.
func (mr *MockProgrammingAgentContextMockRecorder) UpdateFileContent(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateFileContent", reflect.TypeOf((*MockProgrammingAgentContext)(nil).UpdateFileContent), arg0, arg1)
}
