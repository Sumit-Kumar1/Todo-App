package client

import (
	"encoding/json"
	"net/http"
	"todoapp/internal/errors"
	"todoapp/internal/models"

	"gofr.dev/pkg/gofr"
	"gofr.dev/pkg/gofr/service"
)

var (
	invalidStatus = errors.NewConstError("invalid status found")
)

type Client struct {
	service service.HTTP
}

func New(svc service.HTTP) *Client {
	return &Client{service: svc}
}

func (c *Client) SignUp(ctx *gofr.Context, email, password string) error {
	var req = models.AuthUserReq{Email: email, Password: password}

	data, err := json.Marshal(req)
	if err != nil {
		return err
	}

	resp, err := c.service.Post(ctx, "signup", nil, data)
	if err != nil {
		return err
	}

	if _, err = handleResponse[models.AuthUserResp](resp); err != nil {
		return err
	}

	return nil
}

func (c *Client) SignIn(ctx *gofr.Context, email, password string) (*models.AuthUserResp, error) {
	var req = models.AuthUserReq{Email: email, Password: password}

	data, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	resp, err := c.service.Post(ctx, "signin", nil, data)
	if err != nil {
		return nil, err
	}

	authUser, err := handleResponse[models.AuthUserResp](resp)
	if err != nil {
		return nil, err
	}

	return authUser, nil
}

func (c *Client) Refresh(ctx *gofr.Context) (*string, error) {
	headers := make(map[string]string)

	resp, err := c.service.PostWithHeaders(ctx, "refresh", nil, nil, headers)
	if err != nil {
		return nil, err
	}

	var token struct {
		RefToken *string `json:"refreshToken,omitempty"`
	}

	if resp == nil {
		return nil, errors.NewConstError("auth service - nil response for refresh")
	}

	if resp.StatusCode != http.StatusCreated {
		return nil, invalidStatus
	}

	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
		return nil, err
	}

	return token.RefToken, nil
}

func (c *Client) Revoke(ctx *gofr.Context, token string) error {
	// TODO: complete revoke functionality
	return nil
}

func handleResponse[T any](resp *http.Response) (*T, error) {
	if resp == nil {
		return nil, errors.ErrRequired("auth-rest-api : response body")
	}

	var res T

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
		return nil, errors.ErrInvalid("auth-response : status code")
	}

	return &res, nil
}
