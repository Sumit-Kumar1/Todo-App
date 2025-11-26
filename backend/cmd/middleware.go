package cmd

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"todoapp/internal/errors"
	"todoapp/internal/models"

	"github.com/google/uuid"
)

const (
	cookieName     = "auth"
	allowedOrigin  = "http://localhost:4173"
	allowedMethods = "POST, GET, PUT, DELETE, PATCH, OPTIONS"
	allowedHeaders = "Accept, Content-Type, Content-Length, Accept-Encoding"
	corsMaxAge     = "8640"
)

type Claims struct {
	Email    string `json:"email"`
	ClaimUID string `json:"claimID"`
	// registeredClaim's subject is userID from db
	jwt.RegisteredClaims
}

type Middleware func(http.HandlerFunc) http.HandlerFunc

func Chain(f http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
	for _, m := range middlewares {
		f = m(f)
	}

	return f
}

// AuthMiddleware check for valid cookie and extracts user id for the user
// TODO: Decide wether if token is expired had to refresh and add a new auth cookie to disrupt the flow ?
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		uid, err := validateCookie(c)
		if err != nil {
			errors.HandleHTTPError(c, err)
			return
		}

		c.Request.WithContext(context.WithValue(c, models.CtxKeyUserID, *uid))
	}

}

func validateCookie(c *gin.Context) (*uuid.UUID, error) {
	val, err := c.Cookie(cookieName)
	if err != nil {
		return nil, err
	}

	token, err := jwt.ParseWithClaims(val, &Claims{}, func(token *jwt.Token) (any, error) {
		secVal := "33cea8f88c5c8ad73b1700af7d72891fe3097297e59fb6cbe5fd8b545a8316d0"
		return []byte(secVal), nil
	}, jwt.WithExpirationRequired())
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, errors.ErrInvalid("claim")
	}

	if !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}

	return extractUserID(claims.Subject)
}

func extractUserID(claimSubject string) (*uuid.UUID, error) {
	uid, err := uuid.Parse(claimSubject)
	if err != nil {
		return nil, err
	}

	if uid == uuid.Nil {
		return nil, errors.ErrInvalid("nil user id in claims")
	}

	return &uid, nil
}
