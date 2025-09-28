package todocookie

import (
	"encoding/base64"
	"net/http"

	"todoapp/internal/errors"
)

func WriteCookie(w http.ResponseWriter, c http.Cookie) error {
	// Encode the cookie value using base64, as to accept char apart from utf-16
	c.Value = base64.URLEncoding.EncodeToString([]byte(c.Value))

	if len(c.String()) > 4096 {
		return errors.ErrCookieValTooLong
	}

	http.SetCookie(w, &c)

	return nil
}

func ReadCookie(r *http.Request, cookieName string) (string, error) {
	cookie, err := r.Cookie(cookieName)
	if err != nil {
		return "", err
	}

	val, err := base64.URLEncoding.DecodeString(cookie.Value)
	if err != nil {
		return "", errors.ErrInvalidCookie
	}

	return string(val), nil
}
