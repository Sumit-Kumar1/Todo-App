package usersvc

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"todoapp/internal/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
)

var errMock = models.NewConstError("some error")

const testFailFmt = "Test[%d] failed - %s"

func TestServiceRegister(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userMock := NewMockUserStorer(ctrl)
	sessionMock := NewMockSessionStorer(ctrl)
	s := New(userMock, sessionMock)
	email := "abcd@cdef.com"
	ctx := context.Background()
	userData := models.UserData{}
	longPass := strings.Repeat("abcd", 20)
	req := models.RegisterReq{
		Name:     "Hello world",
		LoginReq: &models.LoginReq{Email: email, Password: "abcd@abcd"},
	}

	tests := []struct {
		name     string
		req      *models.RegisterReq
		mockCall func(*MockUserStorer, *MockSessionStorer)
		wantRes  any
		wantErr  error
	}{
		{name: "nil request", req: nil, mockCall: nil, wantRes: nil},
		{
			name:     "invalid request",
			req:      &models.RegisterReq{Name: ""},
			wantErr:  models.ErrRequired("name"),
			wantRes:  nil,
			mockCall: nil,
		},
		{name: "User already exists", req: &req, wantErr: models.ErrUserAlreadyExists,
			mockCall: func(mock *MockUserStorer, _ *MockSessionStorer) {
				mock.EXPECT().GetUserByEmail(ctx, email).Return(&userData, nil)
			}, wantRes: nil,
		},
		{name: "user not found", req: &req, wantErr: errMock, wantRes: nil,
			mockCall: func(mus *MockUserStorer, _ *MockSessionStorer) {
				mus.EXPECT().GetUserByEmail(ctx, email).Return(nil, errMock)
			}},
		{name: "pass encrypt error", req: &models.RegisterReq{Name: req.Name,
			LoginReq: &models.LoginReq{Email: email, Password: longPass}},
			wantErr: bcrypt.ErrPasswordTooLong, wantRes: nil,
			mockCall: func(mus *MockUserStorer, _ *MockSessionStorer) {
				mus.EXPECT().GetUserByEmail(ctx, email).Return(nil, models.ErrNotFound("user"))
			}},
		{name: "error while registering user", req: &req, wantErr: errMock,
			mockCall: func(mus *MockUserStorer, _ *MockSessionStorer) {
				mus.EXPECT().GetUserByEmail(ctx, email).Return(nil, nil)
				mus.EXPECT().RegisterUser(ctx, gomock.Any()).Return(errMock)
			}, wantRes: nil,
		},
		{name: "error while creating session", req: &req, wantErr: errMock,
			mockCall: func(mus *MockUserStorer, mss *MockSessionStorer) {
				mus.EXPECT().GetUserByEmail(ctx, email).Return(nil, nil)
				mus.EXPECT().RegisterUser(ctx, gomock.Any()).Return(nil)
				mss.EXPECT().CreateSession(ctx, gomock.Any()).Return(errMock)
			}, wantRes: nil,
		},
		{name: "valid user register flow", req: &req, wantErr: nil,
			mockCall: func(mus *MockUserStorer, mss *MockSessionStorer) {
				mus.EXPECT().GetUserByEmail(ctx, email).Return(nil, nil)
				mus.EXPECT().RegisterUser(ctx, gomock.Any()).Return(nil)
				mss.EXPECT().CreateSession(ctx, gomock.Any()).Return(nil)
			}, wantRes: req,
		},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockCall != nil {
				tt.mockCall(userMock, sessionMock)
			}

			got, err := s.Register(context.Background(), tt.req)

			assert.Equalf(t, tt.wantErr, err, testFailFmt, i, tt.name)

			if tt.wantRes != nil {
				assert.NotNilf(t, got.Expiry, testFailFmt, i, tt.name)
			}
		})
	}
}

func TestServiceLogin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSession := NewMockSessionStorer(ctrl)
	mockUser := NewMockUserStorer(ctrl)
	ctx := context.Background()
	id := uuid.New()
	pass := "abcd@abcd"
	email := "abcd@cdef.com"
	req := models.LoginReq{Email: email, Password: pass}
	invalidUsr := models.UserData{ID: id, Password: pass, Name: "hello", Email: email}
	encPass, _ := encryptedPassword(pass)
	usr := models.UserData{ID: id, Name: "Hello world", Email: email, Password: encPass}
	ss := models.SessionData{
		ID:     uuid.New(),
		UserID: usr.ID,
		Token:  uuid.NewString(),
		Expiry: time.Now().Add(1 * time.Minute),
	}

	tests := []struct {
		name     string
		req      *models.LoginReq
		mockCall func(*MockUserStorer, *MockSessionStorer)
		want     *models.SessionData
		wantErr  error
	}{
		{name: "nil request", req: nil, want: nil, wantErr: models.ErrRequired("login request")},
		{
			name:    "invalid request",
			req:     &models.LoginReq{Email: ""},
			wantErr: models.ErrRequired("email"),
		},
		{
			name: "user get error",
			req:  &req,
			mockCall: func(mus *MockUserStorer, _ *MockSessionStorer) {
				mus.EXPECT().GetUserByEmail(ctx, email).Return(nil, errMock)
			},
			wantErr: errMock,
		},
		{
			name: "nil user in get",
			req:  &req,
			mockCall: func(mus *MockUserStorer, _ *MockSessionStorer) {
				mus.EXPECT().GetUserByEmail(ctx, email).Return(nil, nil)
			},
			wantErr: models.ErrNotFound("user"),
		},
		{
			name: "passwd not matching",
			req:  &req,
			mockCall: func(mus *MockUserStorer, _ *MockSessionStorer) {
				mus.EXPECT().GetUserByEmail(ctx, email).Return(&invalidUsr, nil)
			},
			wantErr: models.ErrPsswdNotMatch,
		},
		{
			name: "valid login flow",
			req:  &req,
			mockCall: func(mus *MockUserStorer, mss *MockSessionStorer) {
				mus.EXPECT().GetUserByEmail(ctx, email).Return(&usr, nil)
				mss.EXPECT().GetSessionByID(ctx, &usr.ID).Return(&ss, nil)
			},
			wantErr: nil,
			want:    &ss,
		},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{UserStore: mockUser, SessionStore: mockSession}

			if tt.mockCall != nil {
				tt.mockCall(mockUser, mockSession)
			}

			got, err := s.Login(ctx, tt.req)

			assert.Equalf(t, tt.wantErr, err, testFailFmt, i, tt.name)
			assert.Equalf(t, tt.want, got, testFailFmt, i, tt.want)
		})
	}
}

func TestServiceLogout(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
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
		{
			name:     "invalid case",
			password: strings.Repeat("a", 100),
			wantErr:  bcrypt.ErrPasswordTooLong,
		},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := encryptedPassword(tt.password)

			assert.Equalf(t, tt.wantErr, err, testFailFmt, i, tt.name)
			assert.NotEqualf(t, tt.password, got, testFailFmt, i, tt.name)
		})
	}
}

func TestServiceHandleLoginSession(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSession := NewMockSessionStorer(ctrl)
	s := &Service{SessionStore: mockSession}
	ctx := context.Background()
	uid := uuid.New()
	tt := time.Now().Add(time.Minute * 10).UTC()
	ss := models.SessionData{UserID: uid, ID: uuid.New(), Token: uuid.NewString(), Expiry: tt}

	tests := []struct {
		name    string
		user    *models.UserData
		mockFxn func(*MockSessionStorer)
		want    *models.SessionData
		wantErr error
	}{
		{name: "valid case: session exists", user: &models.UserData{ID: uid},
			mockFxn: func(mss *MockSessionStorer) {
				mss.EXPECT().GetSessionByID(ctx, &uid).Return(&ss, nil)
			},
			want: &models.SessionData{UserID: uid}, wantErr: nil,
		},
		{name: "valid case: expired session exists", user: &models.UserData{ID: uid},
			mockFxn: func(mss *MockSessionStorer) {
				mss.EXPECT().GetSessionByID(ctx, &uid).Return(&models.SessionData{UserID: uid}, nil)
				mss.EXPECT().RefreshSession(ctx, gomock.Any()).Return(nil)
			},
			want: &models.SessionData{UserID: uid}, wantErr: nil,
		},
		{name: "refresh err case: expired session exists", user: &models.UserData{ID: uid},
			mockFxn: func(mss *MockSessionStorer) {
				mss.EXPECT().GetSessionByID(ctx, &uid).Return(&models.SessionData{UserID: uid}, nil)
				mss.EXPECT().RefreshSession(ctx, gomock.Any()).Return(errMock)
			},
			want: nil, wantErr: errMock,
		},
		{name: "err while getting session", user: &models.UserData{ID: uid},
			mockFxn: func(mss *MockSessionStorer) {
				mss.EXPECT().GetSessionByID(ctx, &uid).Return(nil, errMock)
			},
			want: nil, wantErr: errMock,
		},
		{
			name: "valid case: session not found creating new session",
			user: &models.UserData{ID: uid},
			mockFxn: func(mss *MockSessionStorer) {
				mss.EXPECT().GetSessionByID(ctx, &uid).Return(nil, models.ErrNotFound("user ID"))
				mss.EXPECT().CreateSession(ctx, gomock.Any()).Return(nil)
			},
			want:    &ss,
			wantErr: nil,
		},
		{name: "err case: session not found creating new session", user: &models.UserData{ID: uid},
			mockFxn: func(mss *MockSessionStorer) {
				mss.EXPECT().GetSessionByID(ctx, &uid).Return(nil, models.ErrNotFound("user ID"))
				mss.EXPECT().CreateSession(ctx, gomock.Any()).Return(errMock)
			},
			want: nil, wantErr: errMock,
		},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFxn(mockSession)

			got, err := s.handleLoginSession(ctx, tt.user)

			assert.Equalf(t, tt.wantErr, err, "Test[%d] failed - %s", i, tt.name)

			if got != nil {
				assert.Equalf(t, tt.want.UserID, got.UserID, "Test[%d] failed", i, tt.name)
			}
		})
	}
}
