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
	headers := map[string]string{models.HeaderCorrelation: models.GetCorrelationID(ctx)}

	resp, err := c.postWithHeaders(ctx, "/signup", headers, req)
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
	headers := map[string]string{models.HeaderCorrelation: models.GetCorrelationID(ctx)}

	resp, err := c.postWithHeaders(ctx, "/signin", headers, req)
	if err != nil {
		return nil, err
	}

	return handleResponse[models.AuthUserResp](resp)
}

func (c *Client) Refresh(ctx context.Context, auth string) (*string, error) {
	headers := map[string]string{
		models.HeaderCorrelation: models.GetCorrelationID(ctx),
		"Authorization":          "Bearer " + auth,
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.url+"/refresh", http.NoBody)
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

	var token struct {
		RefToken *string `json:"refreshToken,omitempty"`
	}

	if resp == nil {
		return nil, errAuthNilResp
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
	// TODO: complete revoke functionality
	return nil
}

func (c *Client) postWithHeaders(ctx context.Context, endpoint string, headers map[string]string, reqModel models.AuthUserReq) (*http.Response, error) {
	data, err := json.Marshal(reqModel)
	if err != nil {
		return nil, err
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

func handleResponse[T any](resp *http.Response) (*T, error) {
	if resp == nil {
		return nil, errAuthNilResp
	}

	var res struct {
		Data T `json:"data,omitempty"`
	}

	switch resp.StatusCode {
	case http.StatusNotFound:
		return nil, errors.ErrUserNotFound
	case http.StatusConflict:
		return nil, errors.ErrUserAlreadyExists
	case http.StatusOK, http.StatusCreated:
		if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
			return nil, err
		}

		defer resp.Body.Close()
	default:
		return nil, errInvalidStatus
	}

	return &res.Data, nil
}
