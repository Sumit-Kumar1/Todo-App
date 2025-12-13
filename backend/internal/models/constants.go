package models

type ContextKey string

const (
	Email                        = "email"
	User                         = "user"
	Password                     = "password"
	CtxKeyUserID      ContextKey = "user_id"
	Logger            ContextKey = "logger"
	CorrelationID     ContextKey = "correlationId"
	HeaderCorrelation            = "X-Correlation-ID"
)
