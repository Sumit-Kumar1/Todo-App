package models

import (
	"regexp"
	"strings"
	"time"

	"todoapp/internal/constant"
	"todoapp/internal/errors"

	"github.com/google/uuid"
)

var (
	emailRegex = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
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

type SessionData struct {
	ID     uuid.UUID `json:"id"`
	UserID uuid.UUID `json:"userId"`
	Token  string    `json:"token"`
	Expiry time.Time `json:"expiry"`
}

func (l *LoginReq) Validate() error {
	email := strings.ToLower(strings.TrimSpace(l.Email))
	passwd := strings.TrimSpace(l.Password)

	if email == "" {
		return errors.ErrRequired(constant.Email)
	}

	if !emailRegex.MatchString(email) {
		return errors.ErrInvalid(constant.Email)
	}

	if passwd == "" {
		return errors.ErrRequired(constant.Password)
	}

	if len(passwd) < 8 {
		return errors.ErrInvalid("password is too short")
	}

	return nil
}
