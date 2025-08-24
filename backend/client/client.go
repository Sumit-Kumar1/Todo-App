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

// TODO: enhance these functionality
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

	if resp == nil {
		return errors.NewConstError("auth service - nil response for signup")
	}

	if resp.StatusCode != http.StatusCreated {
		return invalidStatus
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

	if resp == nil {
		return nil, errors.NewConstError("auth service - nil response for signin")
	}

	if resp.StatusCode != http.StatusCreated {
		return nil, invalidStatus
	}

	var authUser models.AuthUserResp

	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&authUser); err != nil {
		return nil, err
	}

	return &authUser, nil
}

func (c *Client) Refresh(ctx *gofr.Context) (*string, error) {
	headers := make(map[string]string)

	resp, err := c.service.PostWithHeaders(ctx, "refresh", nil, nil, headers)
	if err != nil {
		return nil, err
	}

	if resp == nil {
		return nil, errors.NewConstError("auth service - nil response for refresh")
	}

	if resp.StatusCode != http.StatusCreated {
		return nil, invalidStatus
	}

	var token struct {
		RefToten *string `json:"refreshToken,omitempty"`
	}

	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
		return nil, err
	}

	return token.RefToten, nil
}

func (c *Client) Revoke(ctx *gofr.Context) error {
	return nil
}
