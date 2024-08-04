// Code generated by MockGen. DO NOT EDIT.
// Source: interface.go

// Package service is a generated GoMock package.
package service

import (
	context "context"
	reflect "reflect"
	models "todoapp/internal/models"

	gomock "github.com/golang/mock/gomock"
)

// MockStorer is a mock of Storer interface.
type MockStorer struct {
	ctrl     *gomock.Controller
	recorder *MockStorerMockRecorder
}

// MockStorerMockRecorder is the mock recorder for MockStorer.
type MockStorerMockRecorder struct {
	mock *MockStorer
}

// NewMockStorer creates a new mock instance.
func NewMockStorer(ctrl *gomock.Controller) *MockStorer {
	mock := &MockStorer{ctrl: ctrl}
	mock.recorder = &MockStorerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStorer) EXPECT() *MockStorerMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockStorer) Create(ctx context.Context, id, title string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, id, title)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockStorerMockRecorder) Create(ctx, id, title interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockStorer)(nil).Create), ctx, id, title)
}

// Delete mocks base method.
func (m *MockStorer) Delete(ctx context.Context, id string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockStorerMockRecorder) Delete(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockStorer)(nil).Delete), ctx, id)
}

// GetAll mocks base method.
func (m *MockStorer) GetAll(ctx context.Context) ([]models.Task, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAll", ctx)
	ret0, _ := ret[0].([]models.Task)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAll indicates an expected call of GetAll.
func (mr *MockStorerMockRecorder) GetAll(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAll", reflect.TypeOf((*MockStorer)(nil).GetAll), ctx)
}

// MarkDone mocks base method.
func (m *MockStorer) MarkDone(ctx context.Context, id string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MarkDone", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// MarkDone indicates an expected call of MarkDone.
func (mr *MockStorerMockRecorder) MarkDone(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MarkDone", reflect.TypeOf((*MockStorer)(nil).MarkDone), ctx, id)
}

// Update mocks base method.
func (m *MockStorer) Update(ctx context.Context, id, title, isDone string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, id, title, isDone)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockStorerMockRecorder) Update(ctx, id, title, isDone interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockStorer)(nil).Update), ctx, id, title, isDone)
}
