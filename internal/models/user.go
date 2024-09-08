package models

import (
	"errors"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
)

type UserData struct {
	ID       uuid.UUID `json:"id" db:"id"`
	Name     string    `json:"name" db:"name"`
	Email    string    `json:"email" db:"email"`
	Password string    `json:"password" db:"password"`
}

type LoginReq struct {
	Email    string `json:"email" db:"email"`
	Password string `json:"password" db:"password"`
}

type RegisterReq struct {
	Name string `json:"name" db:"name"`
	*LoginReq
}

type UserSession struct {
	ID     uuid.UUID `json:"id" db:"id"`
	UserID uuid.UUID `json:"userId" db:"user_id"`
	Token  string    `json:"token" db:"token"`
	Expiry time.Time `json:"expiry" db:"expiry"`
}

func (l *LoginReq) Validate() error {
	email := strings.ToLower(strings.TrimSpace(l.Email))
	passwd := strings.TrimSpace(l.Password)

	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)

	if email == "" {
		return errors.New("email is required")
	}

	if !emailRegex.MatchString(email) {
		return errors.New("invalid email")
	}

	if passwd == "" {
		return errors.New("password is required")
	}

	if len(passwd) < 8 {
		return errors.New("password is too short")
	}

	return nil
}

func (r *RegisterReq) Validate() error {
	name := strings.TrimSpace(r.Name)
	if name == "" {
		return errors.New("name is required")
	}

	if len(name) < 3 {
		return errors.New("name is too short")
	}

	return r.LoginReq.Validate()
}
