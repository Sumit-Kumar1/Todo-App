package usersvc

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"

	"todoapp/internal/models"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	errMock = models.NewConstError("some error")
)

func TestServiceRegister(t *testing.T) {
	ctrl := gomock.NewController(t)
	userMock := NewMockUserStorer(ctrl)
	sessionMock := NewMockSessionStorer(ctrl)
	s := New(userMock, sessionMock)

	email := "abcd@cdef.com"
	ctx := context.Background()
	userData := models.UserData{}
	longPass := strings.Repeat("abcd", 20)
	req := models.RegisterReq{Name: "Hello world", LoginReq: &models.LoginReq{Email: email, Password: "abcd@abcd"}}

	tests := []struct {
		name     string
		req      *models.RegisterReq
		want     *models.SessionData
		mockCall func(*MockUserStorer, *MockSessionStorer)
		wantErr  error
	}{
		{name: "nil request", req: nil, want: nil, mockCall: func(_ *MockUserStorer, _ *MockSessionStorer) {}},
		{name: "invalid request", req: &models.RegisterReq{Name: ""}, wantErr: models.ErrRequired("name"), mockCall: func(_ *MockUserStorer, _ *MockSessionStorer) {}},
		{name: "User already exists", req: &req, wantErr: models.ErrUserAlreadyExists,
			mockCall: func(mock *MockUserStorer, _ *MockSessionStorer) {
				mock.EXPECT().GetUserByEmail(ctx, email).Return(&userData, nil)
			},
		},
		{name: "user not found", req: &req, wantErr: errMock,
			mockCall: func(mus *MockUserStorer, _ *MockSessionStorer) {
				mus.EXPECT().GetUserByEmail(ctx, email).Return(nil, errMock)
			}},
		{name: "pass encrypt error", req: &models.RegisterReq{Name: req.Name, LoginReq: &models.LoginReq{Email: email, Password: longPass}},
			wantErr: errors.New("bcrypt: password length exceeds 72 bytes"), mockCall: func(mus *MockUserStorer, _ *MockSessionStorer) {
				mus.EXPECT().GetUserByEmail(ctx, email).Return(nil, models.ErrNotFound("user"))
			}},
		{name: "error while creating session", req: &models.RegisterReq{Name: req.Name, LoginReq: &models.LoginReq{Email: email, Password: longPass}},
			wantErr: errors.New("bcrypt: password length exceeds 72 bytes"), mockCall: func(mus *MockUserStorer, _ *MockSessionStorer) {
				mus.EXPECT().GetUserByEmail(ctx, email).Return(nil, models.ErrNotFound("user"))
			}},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockCall(userMock, sessionMock)

			got, err := s.Register(context.Background(), tt.req)
			if err != nil && err.Error() != tt.wantErr.Error() {
				t.Errorf("Test[%d] Failed:\nerror:\t%+v\nwantErr: %+v", i, err, tt.wantErr)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Service.Register() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServiceLogin(t *testing.T) {
	// req := models.LoginReq{Email: "abcd@cdef.com", Password: "abcd@abcd"}
	ctrl := gomock.NewController(t)
	mockSession := NewMockSessionStorer(ctrl)
	mockUser := NewMockUserStorer(ctrl)

	tests := []struct {
		name     string
		req      *models.LoginReq
		mockCall func(*MockUserStorer, *MockSessionStorer)
		want     *models.SessionData
		wantErr  error
	}{
		{name: "nil request", req: nil, want: nil, wantErr: models.ErrRequired("login request")},
		{name: "invalid request", req: &models.LoginReq{Email: ""}, wantErr: models.ErrRequired("email")},
		// {name: "user not found", req: &req, wantErr: nil},
		// {name: "invalid password", req: &req, wantErr: nil},
		// {name: "error while creating session", req: &req, wantErr: nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			s := &Service{UserStore: mockUser, SessionStore: mockSession}

			tt.mockCall(mockUser, mockSession)

			got, err := s.Login(ctx, tt.req)

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Service.Login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != nil && got != tt.want {
				t.Errorf("Service.Login() = %v, want %v", got, tt.want)
			}
		})
	}
	ctrl.Finish()
}

func TestServiceLogout(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockSession := NewMockSessionStorer(ctrl)
	token := uuid.New()
	ctx := context.Background()

	_, uidErr := uuid.Parse("123")

	tests := []struct {
		name     string
		token    string
		mockCall *gomock.Call
		wantErr  error
	}{
		{name: "valid case", token: token.String(),
			mockCall: mockSession.EXPECT().Logout(ctx, &token).Return(nil), wantErr: nil},
		{name: "invalid token", token: "123", wantErr: uidErr},
		{name: "error while logging out", token: token.String(),
			mockCall: mockSession.EXPECT().Logout(ctx, &token).Return(errMock), wantErr: errMock},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				SessionStore: mockSession,
			}

			err := s.Logout(ctx, tt.token)

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Service.Logout() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEncryptedPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  error
	}{
		{name: "valid case", password: "abcd"},
		{name: "valid case 2", password: "abcd124@adgbalje"},
		{name: "invalid case", password: strings.Repeat("a", 100), wantErr: bcrypt.ErrPasswordTooLong},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := encryptedPassword(tt.password)

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Test[%d] Failed - %s\nGOT:\t%v\nWANT:\t%v", i, tt.name, err, tt.wantErr)
			}

			if tt.password == got {
				t.Errorf("Test[%d] Failed - %s\nGOT:%v\nWANT:%v", i, tt.name, got, tt.password)
			}
		})
	}
}
