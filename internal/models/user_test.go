package models

import (
	"errors"
	"testing"
)

const (
	validEmail  = "abcd@abcd.com"
	validPasswd = "abcd@1234"
	validName   = "sumit kumar"
)

func TestLoginReqValidate(t *testing.T) {
	tests := []struct {
		name    string
		email   string
		passwd  string
		wantErr error
	}{
		{name: "valid case", email: validEmail, passwd: validPasswd, wantErr: nil},
		{name: "missing email", email: "", passwd: validPasswd, wantErr: ErrRequired("email")},
		{name: "invalid email", email: "acbcd@abc", passwd: validPasswd, wantErr: ErrInvalid("email")},
		{name: "missing password", email: validEmail, passwd: "", wantErr: ErrRequired("password")},
		{name: "invalid password", email: validEmail, passwd: "abcd", wantErr: ErrInvalid("password is too short")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &LoginReq{
				Email:    tt.email,
				Password: tt.passwd,
			}

			if err := l.Validate(); !errors.Is(err, tt.wantErr) {
				t.Errorf("LoginReq.Validate()::\nGOT:\t%v\nWant:\t%v", err, tt.wantErr)
			}
		})
	}
}

func TestRegisterReqValidate(t *testing.T) {
	tests := []struct {
		name     string
		userName string
		wantErr  error
	}{
		{name: "valid info", userName: "sumit kumar"},
		{name: "missing user name", userName: "", wantErr: ErrRequired("name")},
		{name: "invalid user name", userName: "a", wantErr: ErrInvalid("name is too short")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &RegisterReq{
				Name:     tt.userName,
				LoginReq: &LoginReq{Email: validEmail, Password: validPasswd},
			}

			if err := r.Validate(); !errors.Is(err, tt.wantErr) {
				t.Errorf("RegisterReq.Validate()::\nGOT:\t%v\nWant:\t%v", err, tt.wantErr)
			}
		})
	}
}
