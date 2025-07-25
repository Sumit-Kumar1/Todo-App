package models

import (
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	emailReg = `^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`
)

type UserData struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	Password string    `json:"password"`
}

type LoginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterReq struct {
	Name string `json:"name"`
	*LoginReq
}

type SessionData struct {
	ID     uuid.UUID `json:"id"`
	UserID uuid.UUID `json:"userId"`
	Token  string    `json:"token"`
	Expiry time.Time `json:"expiry"`
}

func (l *LoginReq) Validate() error {
	email := strings.ToLower(strings.TrimSpace(l.Email))
	passwd := strings.TrimSpace(l.Password)
	emailRegex := regexp.MustCompile(emailReg)

	if email == "" {
		return ErrRequired("email")
	}

	if !emailRegex.MatchString(email) {
		return ErrInvalid("email")
	}

	if passwd == "" {
		return ErrRequired("password")
	}

	if len(passwd) < 8 {
		return ErrInvalid("password is too short")
	}

	return nil
}

func (r *RegisterReq) Validate() error {
	name := strings.TrimSpace(r.Name)
	if name == "" {
		return ErrRequired("name")
	}

	if len(name) < 3 {
		return ErrInvalid("name is too short")
	}

	return r.LoginReq.Validate()
}
