// Code generated by MockGen. DO NOT EDIT.
// Source: interface.go

// Package todohttp is a generated GoMock package.
package todohttp

import (
	context "context"
	reflect "reflect"
	models "todoapp/internal/models"

	gomock "github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
)

// MockTodoServicer is a mock of TodoServicer interface.
type MockTodoServicer struct {
	ctrl     *gomock.Controller
	recorder *MockTodoServicerMockRecorder
}

// MockTodoServicerMockRecorder is the mock recorder for MockTodoServicer.
type MockTodoServicerMockRecorder struct {
	mock *MockTodoServicer
}

// NewMockTodoServicer creates a new mock instance.
func NewMockTodoServicer(ctrl *gomock.Controller) *MockTodoServicer {
	mock := &MockTodoServicer{ctrl: ctrl}
	mock.recorder = &MockTodoServicerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTodoServicer) EXPECT() *MockTodoServicerMockRecorder {
	return m.recorder
}

// AddTask mocks base method.
func (m *MockTodoServicer) AddTask(ctx context.Context, task string, userID *uuid.UUID) (*models.Task, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddTask", ctx, task, userID)
	ret0, _ := ret[0].(*models.Task)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddTask indicates an expected call of AddTask.
func (mr *MockTodoServicerMockRecorder) AddTask(ctx, task, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddTask", reflect.TypeOf((*MockTodoServicer)(nil).AddTask), ctx, task, userID)
}

// DeleteTask mocks base method.
func (m *MockTodoServicer) DeleteTask(ctx context.Context, id string, userID *uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteTask", ctx, id, userID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteTask indicates an expected call of DeleteTask.
func (mr *MockTodoServicerMockRecorder) DeleteTask(ctx, id, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteTask", reflect.TypeOf((*MockTodoServicer)(nil).DeleteTask), ctx, id, userID)
}

// GetAll mocks base method.
func (m *MockTodoServicer) GetAll(ctx context.Context, userID *uuid.UUID) ([]models.Task, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAll", ctx, userID)
	ret0, _ := ret[0].([]models.Task)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAll indicates an expected call of GetAll.
func (mr *MockTodoServicerMockRecorder) GetAll(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAll", reflect.TypeOf((*MockTodoServicer)(nil).GetAll), ctx, userID)
}

// MarkDone mocks base method.
func (m *MockTodoServicer) MarkDone(ctx context.Context, id string, userID *uuid.UUID) (*models.Task, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MarkDone", ctx, id, userID)
	ret0, _ := ret[0].(*models.Task)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// MarkDone indicates an expected call of MarkDone.
func (mr *MockTodoServicerMockRecorder) MarkDone(ctx, id, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MarkDone", reflect.TypeOf((*MockTodoServicer)(nil).MarkDone), ctx, id, userID)
}

// UpdateTask mocks base method.
func (m *MockTodoServicer) UpdateTask(ctx context.Context, id, title string, isDone bool, userID *uuid.UUID) (*models.Task, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateTask", ctx, id, title, isDone, userID)
	ret0, _ := ret[0].(*models.Task)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateTask indicates an expected call of UpdateTask.
func (mr *MockTodoServicerMockRecorder) UpdateTask(ctx, id, title, isDone, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateTask", reflect.TypeOf((*MockTodoServicer)(nil).UpdateTask), ctx, id, title, isDone, userID)
}
