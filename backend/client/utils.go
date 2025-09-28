package client

import (
	"context"
	"encoding/json"
	"net/http"
	"todoapp/internal/errors"
	"todoapp/internal/models"
)

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
	case http.StatusUnauthorized, http.StatusForbidden:
		return nil, errors.ErrPsswdNotMatch
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

func prepareAuthAPIHeaders(ctx context.Context, auth string) map[string]string {
	headers := make(map[string]string)

	headers[models.HeaderCorrelation] = models.GetCorrelationID(ctx)

	if auth != "" {
		headers["Authorization"] = "Bearer " + auth
	}

	return headers
}
