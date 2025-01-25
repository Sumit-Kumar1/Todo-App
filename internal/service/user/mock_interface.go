// Code generated by MockGen. DO NOT EDIT.
// Source: interface.go

// Package service is a generated GoMock package.
package usersvc

import (
	context "context"
	reflect "reflect"
	models "todoapp/internal/models"

	gomock "github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
)

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

// CreateSession mocks base method.
func (m *MockUserStorer) CreateSession(ctx context.Context, session *models.UserSession) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateSession", ctx, session)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateSession indicates an expected call of CreateSession.
func (mr *MockUserStorerMockRecorder) CreateSession(ctx, session interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateSession", reflect.TypeOf((*MockUserStorer)(nil).CreateSession), ctx, session)
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

// GetUserByEmail mocks base method.
func (m *MockUserStorer) GetUserByEmail(ctx context.Context, email string) (*models.UserData, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByEmail", ctx, email)
	ret0, _ := ret[0].(*models.UserData)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByEmail indicates an expected call of GetUserByEmail.
func (mr *MockUserStorerMockRecorder) GetUserByEmail(ctx, email interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByEmail", reflect.TypeOf((*MockUserStorer)(nil).GetUserByEmail), ctx, email)
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
func (m *MockUserStorer) RefreshSession(ctx context.Context, newSession *models.UserSession) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RefreshSession", ctx, newSession)
	ret0, _ := ret[0].(error)
	return ret0
}

// RefreshSession indicates an expected call of RefreshSession.
func (mr *MockUserStorerMockRecorder) RefreshSession(ctx, newSession interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RefreshSession", reflect.TypeOf((*MockUserStorer)(nil).RefreshSession), ctx, newSession)
}

// RegisterUser mocks base method.
func (m *MockUserStorer) RegisterUser(ctx context.Context, data *models.UserData) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RegisterUser", ctx, data)
	ret0, _ := ret[0].(error)
	return ret0
}

// RegisterUser indicates an expected call of RegisterUser.
func (mr *MockUserStorerMockRecorder) RegisterUser(ctx, data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterUser", reflect.TypeOf((*MockUserStorer)(nil).RegisterUser), ctx, data)
}
