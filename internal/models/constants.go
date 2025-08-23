package models

type ContextKey string

const (
	Email                   = "email"
	User                    = "user"
	Password                = "password"
	CtxKeyUserID ContextKey = "user_id"
)
