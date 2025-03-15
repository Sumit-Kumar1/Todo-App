package usersvc

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"

	"todoapp/internal/models"

	"github.com/golang/mock/gomock"
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
		want     *models.UserSession
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
		{name: "user not found", req: &req, wantErr: models.NewConstError("some error"),
			mockCall: func(mus *MockUserStorer, _ *MockSessionStorer) {
				mus.EXPECT().GetUserByEmail(ctx, email).Return(nil, models.NewConstError("some error"))
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
		mockCall *gomock.Call
		want     *models.UserSession
		wantErr  error
	}{
		{name: "nil request", req: nil, want: nil, wantErr: nil},
		{name: "invalid request", req: &models.LoginReq{Email: ""}, wantErr: models.ErrRequired("email")},
		// {name: "user not found", req: &req, wantErr: nil},
		// {name: "invalid password", req: &req, wantErr: nil},
		// {name: "error while creating session", req: &req, wantErr: nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			s := &Service{UserStore: mockUser, SessionStore: mockSession}

			got, err := s.Login(ctx, tt.req)

			if err != tt.wantErr {
				t.Errorf("Service.Login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Service.Login() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServiceLogout(t *testing.T) {
	type fields struct {
		Store UserStorer
	}
	type args struct {
		ctx   context.Context
		token string
	}
	var tests []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				UserStore: tt.fields.Store,
			}
			if err := s.Logout(tt.args.ctx, tt.args.token); (err != nil) != tt.wantErr {
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
