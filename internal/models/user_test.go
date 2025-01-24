package models

import (
	"errors"
	"testing"
)

func TestLoginReq_Validate(t *testing.T) {
	tests := []struct {
		name    string
		email   string
		passwd  string
		wantErr error
	}{
		{name: "valid case", email: "abcd@abcd.com", passwd: "abcd@1234", wantErr: nil},
		{name: "missing email", email: "", passwd: "abcd@1234", wantErr: ErrRequired("email")},
		{name: "invalid email", email: "acbcd@abc", passwd: "abcd@1234", wantErr: ErrInvalid("email")},
		{name: "missing password", email: "abcd@abcd.com", passwd: "", wantErr: ErrRequired("password")},
		{name: "invalid password", email: "abcd@abcd.com", passwd: "abcd", wantErr: ErrInvalid("password is too short")},
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

func TestRegisterReq_Validate(t *testing.T) {
	email := "sumit@kumar.com"
	passwd := "abcd@abcd"

	tests := []struct {
		name     string
		userName string
		loginReq *LoginReq
		wantErr  error
	}{
		{name: "valid info", userName: "sumit kumar", loginReq: &LoginReq{Email: email, Password: passwd}},
		{name: "missing user name", userName: "", loginReq: &LoginReq{Email: email, Password: passwd}, wantErr: ErrRequired("name")},
		{name: "invalid user name", userName: "a", loginReq: &LoginReq{Email: email, Password: passwd}, wantErr: ErrInvalid("name is too short")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &RegisterReq{
				Name:     tt.userName,
				LoginReq: tt.loginReq,
			}

			if err := r.Validate(); !errors.Is(err, tt.wantErr) {
				t.Errorf("RegisterReq.Validate()::\nGOT:\t%v\nWant:\t%v", err, tt.wantErr)
			}
		})
	}
}
