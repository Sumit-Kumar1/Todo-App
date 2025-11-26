package userhttp

import (
	"net/http"
	"time"

	"todoapp/internal/models"

	"github.com/gin-gonic/gin"
)

const (
	authCookieName    = "auth"
	refreshCookieName = "refresh"
)

type Handler struct {
	svc UserServicer
}

func New(svc UserServicer) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Register(c *gin.Context) {
	var user models.LoginReq

	if err := c.BindJSON(&user); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if err := h.svc.Register(c, &user); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	resp := "user created successfully"

	c.IndentedJSON(http.StatusCreated, resp)
}

func (h *Handler) Login(c *gin.Context) {
	var user models.LoginReq

	if err := c.BindJSON(&user); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	resp, err := h.svc.Login(c, &user)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	authCookie, refCookie := getLoginCookie(resp)
	c.SetCookieData(&authCookie)
	c.SetCookieData(&refCookie)

	success := "user login successfully"

	c.IndentedJSON(http.StatusOK, success)
}

func (h *Handler) Logout(c *gin.Context) {
	authCookie, err := c.Cookie(authCookieName)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if err := h.svc.Logout(c, authCookie); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.SetCookieData(&http.Cookie{
		Name:     authCookieName,
		Value:    "",
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
		MaxAge:   -1,
		SameSite: http.SameSiteLaxMode,
	})

	success := "user logout successfully"

	c.IndentedJSON(http.StatusOK, success)
}

func getLoginCookie(resp *models.AuthUserResp) (auth, refresh http.Cookie) {
	if resp == nil {
		return
	}

	auth = http.Cookie{
		Name:     authCookieName,
		Value:    resp.AccessToken,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
		MaxAge:   int(time.Minute * 10),
		SameSite: http.SameSiteLaxMode,
	}

	refresh = http.Cookie{
		Name:     refreshCookieName,
		Value:    resp.RefreshToken,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
		MaxAge:   int(time.Hour * 6),
		SameSite: http.SameSiteLaxMode,
	}

	return
}
