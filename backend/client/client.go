// Package client helps us to communicate with auth-rest-api
package client

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"
	"sync"
	"time"

	"todoapp/internal/errors"
	"todoapp/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	healthEndpoint = "/.well-known/health"
	healthTimeout  = 5 * time.Second
)

const (
	errInvalidStatus errors.ConstError = "invalid status found"
	errAuthNilResp   errors.ConstError = "nil response from auth-rest-api"
)

var (
	once           sync.Once
	clientInstance *Client
)

type Client struct {
	http.Client
	url string
}

func New(url string) *Client {
	if url == "" {
		slog.Error("auth-rest-api URL cannot be empty")
		return nil
	}

	url = strings.Trim(url, `"`)

	once.Do(func() {
		slog.LogAttrs(context.Background(), slog.LevelInfo, "creating auth‑rest‑api client")
		clientInstance = &Client{
			url: url,
			Client: http.Client{
				Transport: http.DefaultTransport,
				Timeout:   10 * time.Second, // generic timeout for all calls
			},
		}

		if err := clientInstance.health(context.Background()); err != nil {
			slog.Error("auth‑rest‑api health check failed during client creation", "err", err)
			panic(err)
		}
	})

	return clientInstance
}

func (c *Client) health(ctx context.Context) error {
	hcCtx, cancel := context.WithTimeout(ctx, healthTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(hcCtx, http.MethodGet, c.url+healthEndpoint, nil)
	if err != nil {
		return err
	}

	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	defer ensureBodyClosed(resp)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return errors.Invalid("auth service health check failed: " + resp.Status)
	}

	return nil
}

func (c *Client) SignUp(ctx *gin.Context, email, password string) error {
	req := models.AuthUserReq{Email: email, Password: password}
	headers := prepareAuthAPIHeaders(ctx, "")

	resp, err := c.postWithHeaders(ctx, "/signup", headers, &req)
	if err != nil {
		return err
	}

	defer ensureBodyClosed(resp)

	_, err = handleResponse[string](resp)
	return err
}

func (c *Client) SignIn(ctx *gin.Context, email, password string) (*models.AuthUserResp, error) {
	req := models.AuthUserReq{Email: email, Password: password}
	headers := prepareAuthAPIHeaders(ctx, "")

	resp, err := c.postWithHeaders(ctx, "/signin", headers, &req)
	if err != nil {
		return nil, err
	}

	return handleResponse[models.AuthUserResp](resp)
}

func (c *Client) Refresh(ctx *gin.Context, auth string) (*string, error) {
	if auth == "" {
		return nil, errors.Required("auth token")
	}

	headers := prepareAuthAPIHeaders(ctx, auth)

	resp, err := c.postWithHeaders(ctx, "/refresh", headers, nil)
	if err != nil {
		return nil, err
	}

	defer ensureBodyClosed(resp)

	if resp == nil {
		return nil, errors.Required("response for refresh")
	}

	if resp.StatusCode != http.StatusCreated {
		return nil, errInvalidStatus
	}

	var token struct {
		RefToken *string `json:"refreshToken,omitempty"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
		return nil, err
	}

	return token.RefToken, nil
}

func (c *Client) Revoke(ctx *gin.Context, token string) error {
	if token == "" {
		return errors.Required("token")
	}

	headers := prepareAuthAPIHeaders(ctx, token)

	resp, err := c.postWithHeaders(ctx, "/revoke", headers, nil)
	if err != nil {
		return err
	}

	defer ensureBodyClosed(resp)

	if resp == nil {
		return errors.Required("response from revoke")
	}

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return errors.Invalid("revoke status code: " + resp.Status)
	}

	return nil
}

func (c *Client) Validate(ctx *gin.Context, token string) (*uuid.UUID, error) {
	slog.LogAttrs(ctx, slog.LevelInfo, "validation called !!")

	if token == "" {
		return nil, errors.Required("token")
	}

	headers := prepareAuthAPIHeaders(ctx, token)

	resp, err := c.postWithHeaders(ctx, "/validate", headers, nil)
	if err != nil {
		return nil, err
	}

	valResp, err := handleResponse[models.ValidationResp](resp)
	if err != nil {
		return nil, err
	}

	return extractUserID(&valResp.UserID)
}

func (c *Client) postWithHeaders(ctx *gin.Context, endpoint string, headers map[string]string, reqModel any) (*http.Response, error) {
	var data []byte
	if reqModel != nil {
		var err error
		data, err = json.Marshal(reqModel)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.url+endpoint, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	for key, val := range headers {
		req.Header.Set(key, val)
	}

	return c.Do(req)
}

// ensureBodyClosed safely closes the response body if it exists
func ensureBodyClosed(resp *http.Response) {
	if resp != nil && resp.Body != nil {
		_ = resp.Body.Close()
	}
}

func extractUserID(claimSubject *string) (*uuid.UUID, error) {
	if claimSubject == nil {
		return nil, errors.Invalid("nil userID")
	}

	uid, err := uuid.Parse(*claimSubject)
	if err != nil {
		return nil, err
	}

	if uid == uuid.Nil {
		return nil, errors.ErrInvalidCookie
	}

	return &uid, nil
}
