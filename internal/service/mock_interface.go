// Code generated by MockGen. DO NOT EDIT.
// Source: interface.go

// Package service is a generated GoMock package.
package service

import (
	context "context"
	reflect "reflect"
	models "todoapp/internal/models"

	gomock "github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
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
func (m *MockStorer) Create(ctx context.Context, id, title string, userID *uuid.UUID) (*models.Task, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, id, title, userID)
	ret0, _ := ret[0].(*models.Task)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockStorerMockRecorder) Create(ctx, id, title, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockStorer)(nil).Create), ctx, id, title, userID)
}

// Delete mocks base method.
func (m *MockStorer) Delete(ctx context.Context, id string, userID *uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, id, userID)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockStorerMockRecorder) Delete(ctx, id, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockStorer)(nil).Delete), ctx, id, userID)
}

// GetAll mocks base method.
func (m *MockStorer) GetAll(ctx context.Context, userID *uuid.UUID) ([]models.Task, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAll", ctx, userID)
	ret0, _ := ret[0].([]models.Task)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAll indicates an expected call of GetAll.
func (mr *MockStorerMockRecorder) GetAll(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAll", reflect.TypeOf((*MockStorer)(nil).GetAll), ctx, userID)
}

// GetByEmail mocks base method.
func (m *MockStorer) GetByEmail(ctx context.Context, email string) (*models.UserData, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByEmail", ctx, email)
	ret0, _ := ret[0].(*models.UserData)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByEmail indicates an expected call of GetByEmail.
func (mr *MockStorerMockRecorder) GetByEmail(ctx, email interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByEmail", reflect.TypeOf((*MockStorer)(nil).GetByEmail), ctx, email)
}

// GetSessionByID mocks base method.
func (m *MockStorer) GetSessionByID(ctx context.Context, userID *uuid.UUID) (*models.UserSession, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSessionByID", ctx, userID)
	ret0, _ := ret[0].(*models.UserSession)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSessionByID indicates an expected call of GetSessionByID.
func (mr *MockStorerMockRecorder) GetSessionByID(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSessionByID", reflect.TypeOf((*MockStorer)(nil).GetSessionByID), ctx, userID)
}

// Logout mocks base method.
func (m *MockStorer) Logout(ctx context.Context, token *uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Logout", ctx, token)
	ret0, _ := ret[0].(error)
	return ret0
}

// Logout indicates an expected call of Logout.
func (mr *MockStorerMockRecorder) Logout(ctx, token interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Logout", reflect.TypeOf((*MockStorer)(nil).Logout), ctx, token)
}

// MarkDone mocks base method.
func (m *MockStorer) MarkDone(ctx context.Context, id string, userID *uuid.UUID) (*models.Task, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MarkDone", ctx, id, userID)
	ret0, _ := ret[0].(*models.Task)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// MarkDone indicates an expected call of MarkDone.
func (mr *MockStorerMockRecorder) MarkDone(ctx, id, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MarkDone", reflect.TypeOf((*MockStorer)(nil).MarkDone), ctx, id, userID)
}

// RefreshSession mocks base method.
func (m *MockStorer) RefreshSession(ctx context.Context, newSession *models.UserSession) (*models.UserSession, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RefreshSession", ctx, newSession)
	ret0, _ := ret[0].(*models.UserSession)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RefreshSession indicates an expected call of RefreshSession.
func (mr *MockStorerMockRecorder) RefreshSession(ctx, newSession interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RefreshSession", reflect.TypeOf((*MockStorer)(nil).RefreshSession), ctx, newSession)
}

// RegisterUser mocks base method.
func (m *MockStorer) RegisterUser(ctx context.Context, data *models.UserData, session *models.UserSession) (*models.UserSession, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RegisterUser", ctx, data, session)
	ret0, _ := ret[0].(*models.UserSession)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RegisterUser indicates an expected call of RegisterUser.
func (mr *MockStorerMockRecorder) RegisterUser(ctx, data, session interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterUser", reflect.TypeOf((*MockStorer)(nil).RegisterUser), ctx, data, session)
}

// Update mocks base method.
func (m *MockStorer) Update(ctx context.Context, id, title string, userID *uuid.UUID) (*models.Task, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, id, title, userID)
	ret0, _ := ret[0].(*models.Task)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Update indicates an expected call of Update.
func (mr *MockStorerMockRecorder) Update(ctx, id, title, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockStorer)(nil).Update), ctx, id, title, userID)
}

// MockUserStorer is a mock of UserStorer interface.
type MockUserStorer struct {
	ctrl     *gomock.Controller
	recorder *MockUserStorerMockRecorder
}

// MockUserStorerMockRecorder is the mock recorder for MockUserStorer.
type MockUserStorerMockRecorder struct {
	mock *MockUserStorer
}

// NewMockUserStorer creates a new mock instance.
func NewMockUserStorer(ctrl *gomock.Controller) *MockUserStorer {
	mock := &MockUserStorer{ctrl: ctrl}
	mock.recorder = &MockUserStorerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserStorer) EXPECT() *MockUserStorerMockRecorder {
	return m.recorder
}

// GetByEmail mocks base method.
func (m *MockUserStorer) GetByEmail(ctx context.Context, email string) (*models.UserData, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByEmail", ctx, email)
	ret0, _ := ret[0].(*models.UserData)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByEmail indicates an expected call of GetByEmail.
func (mr *MockUserStorerMockRecorder) GetByEmail(ctx, email interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByEmail", reflect.TypeOf((*MockUserStorer)(nil).GetByEmail), ctx, email)
}

// GetSessionByID mocks base method.
func (m *MockUserStorer) GetSessionByID(ctx context.Context, userID *uuid.UUID) (*models.UserSession, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSessionByID", ctx, userID)
	ret0, _ := ret[0].(*models.UserSession)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSessionByID indicates an expected call of GetSessionByID.
func (mr *MockUserStorerMockRecorder) GetSessionByID(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSessionByID", reflect.TypeOf((*MockUserStorer)(nil).GetSessionByID), ctx, userID)
}

// Logout mocks base method.
func (m *MockUserStorer) Logout(ctx context.Context, token *uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Logout", ctx, token)
	ret0, _ := ret[0].(error)
	return ret0
}

// Logout indicates an expected call of Logout.
func (mr *MockUserStorerMockRecorder) Logout(ctx, token interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Logout", reflect.TypeOf((*MockUserStorer)(nil).Logout), ctx, token)
}

// RefreshSession mocks base method.
func (m *MockUserStorer) RefreshSession(ctx context.Context, newSession *models.UserSession) (*models.UserSession, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RefreshSession", ctx, newSession)
	ret0, _ := ret[0].(*models.UserSession)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RefreshSession indicates an expected call of RefreshSession.
func (mr *MockUserStorerMockRecorder) RefreshSession(ctx, newSession interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RefreshSession", reflect.TypeOf((*MockUserStorer)(nil).RefreshSession), ctx, newSession)
}

// RegisterUser mocks base method.
func (m *MockUserStorer) RegisterUser(ctx context.Context, data *models.UserData, session *models.UserSession) (*models.UserSession, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RegisterUser", ctx, data, session)
	ret0, _ := ret[0].(*models.UserSession)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RegisterUser indicates an expected call of RegisterUser.
func (mr *MockUserStorerMockRecorder) RegisterUser(ctx, data, session interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterUser", reflect.TypeOf((*MockUserStorer)(nil).RegisterUser), ctx, data, session)
}

// MockTodoStorer is a mock of TodoStorer interface.
type MockTodoStorer struct {
	ctrl     *gomock.Controller
	recorder *MockTodoStorerMockRecorder
}

// MockTodoStorerMockRecorder is the mock recorder for MockTodoStorer.
type MockTodoStorerMockRecorder struct {
	mock *MockTodoStorer
}

// NewMockTodoStorer creates a new mock instance.
func NewMockTodoStorer(ctrl *gomock.Controller) *MockTodoStorer {
	mock := &MockTodoStorer{ctrl: ctrl}
	mock.recorder = &MockTodoStorerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTodoStorer) EXPECT() *MockTodoStorerMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockTodoStorer) Create(ctx context.Context, id, title string, userID *uuid.UUID) (*models.Task, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, id, title, userID)
	ret0, _ := ret[0].(*models.Task)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockTodoStorerMockRecorder) Create(ctx, id, title, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockTodoStorer)(nil).Create), ctx, id, title, userID)
}

// Delete mocks base method.
func (m *MockTodoStorer) Delete(ctx context.Context, id string, userID *uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, id, userID)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockTodoStorerMockRecorder) Delete(ctx, id, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockTodoStorer)(nil).Delete), ctx, id, userID)
}

// GetAll mocks base method.
func (m *MockTodoStorer) GetAll(ctx context.Context, userID *uuid.UUID) ([]models.Task, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAll", ctx, userID)
	ret0, _ := ret[0].([]models.Task)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAll indicates an expected call of GetAll.
func (mr *MockTodoStorerMockRecorder) GetAll(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAll", reflect.TypeOf((*MockTodoStorer)(nil).GetAll), ctx, userID)
}

// MarkDone mocks base method.
func (m *MockTodoStorer) MarkDone(ctx context.Context, id string, userID *uuid.UUID) (*models.Task, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MarkDone", ctx, id, userID)
	ret0, _ := ret[0].(*models.Task)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// MarkDone indicates an expected call of MarkDone.
func (mr *MockTodoStorerMockRecorder) MarkDone(ctx, id, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MarkDone", reflect.TypeOf((*MockTodoStorer)(nil).MarkDone), ctx, id, userID)
}

// Update mocks base method.
func (m *MockTodoStorer) Update(ctx context.Context, id, title string, userID *uuid.UUID) (*models.Task, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, id, title, userID)
	ret0, _ := ret[0].(*models.Task)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Update indicates an expected call of Update.
func (mr *MockTodoStorerMockRecorder) Update(ctx, id, title, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockTodoStorer)(nil).Update), ctx, id, title, userID)
}
