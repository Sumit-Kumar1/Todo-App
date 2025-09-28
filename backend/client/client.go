package client

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"sync"

	"todoapp/internal/errors"
	"todoapp/internal/models"
)

var (
	errInvalidStatus = errors.NewConstError("invalid status found")
	errAuthNilResp   = errors.NewConstError("nil response from auth-rest-api")
	once             sync.Once
	clientInstance   *Client
)

type Client struct {
	http.Client
	url string
}

func New(url string) *Client {
	if clientInstance == nil {
		slog.LogAttrs(context.Background(), slog.LevelInfo, "creating auth-rest-api client")
		once.Do(func() {
			clientInstance = &Client{url: url}
			clientInstance.Transport = http.DefaultTransport
		})

		return clientInstance
	}

	slog.LogAttrs(context.Background(), slog.LevelInfo, "using existing auth-rest-api client")

	return clientInstance
}

func (c *Client) SignUp(ctx context.Context, email, password string) error {
	req := models.AuthUserReq{Email: email, Password: password}
	headers := prepareAuthAPIHeaders(ctx, "")

	resp, err := c.postWithHeaders(ctx, "/signup", headers, &req)
	if err != nil {
		return err
	}

	_, err = handleResponse[string](resp)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) SignIn(ctx context.Context, email, password string) (*models.AuthUserResp, error) {
	req := models.AuthUserReq{Email: email, Password: password}
	headers := prepareAuthAPIHeaders(ctx, "")

	resp, err := c.postWithHeaders(ctx, "/signin", headers, &req)
	if err != nil {
		return nil, err
	}

	return handleResponse[models.AuthUserResp](resp)
}

func (c *Client) Refresh(ctx context.Context, auth string) (*string, error) {
	var token struct {
		RefToken *string `json:"refreshToken,omitempty"`
	}

	headers := prepareAuthAPIHeaders(ctx, auth)

	resp, err := c.postWithHeaders(ctx, "/refresh", headers, nil)
	if err != nil {
		return nil, err
	}

	if resp == nil {
		return nil, errors.ErrRequired("response for refresh")
	}

	if resp.StatusCode != http.StatusCreated {
		return nil, errInvalidStatus
	}

	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
		return nil, err
	}

	return token.RefToken, nil
}

func (c *Client) Revoke(ctx context.Context, token string) error {
	headers := prepareAuthAPIHeaders(ctx, token)

	resp, err := c.postWithHeaders(ctx, "/revoke", headers, nil)
	if err != nil {
		return err
	}

	if resp == nil {
		return errors.ErrRequired("response from revoke")
	}

	if resp.StatusCode != http.StatusNoContent {
		return errors.ErrInvalid("revoke status code: " + resp.Status)
	}

	return nil
}

func (c *Client) postWithHeaders(ctx context.Context, endpoint string, headers map[string]string, reqModel *models.AuthUserReq) (*http.Response, error) {
	var (
		data []byte
		err  error
	)

	if reqModel != nil {
		data, err = json.Marshal(*reqModel)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.url+endpoint, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	for key, val := range headers {
		req.Header.Add(key, val)
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
