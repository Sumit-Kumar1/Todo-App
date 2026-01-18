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
		Data    T       `json:"data,omitempty"`
		Code    *int    `json:"code,omitempty"`
		Message *string `json:"message,omitempty"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusNotFound:
		return nil, errors.ErrUserNotFound
	case http.StatusConflict:
		return nil, errors.ErrUserAlreadyExists
	case http.StatusBadRequest:
		return nil, errors.NewHTTPError(*res.Code, *res.Message, "")
	case http.StatusOK, http.StatusCreated:
		return &res.Data, nil
	default:
		return nil, errInvalidStatus
	}
}

func prepareAuthAPIHeaders(ctx context.Context, auth string) map[string]string {
	headers := make(map[string]string)

	headers[models.HeaderCorrelation] = models.GetCorrelationID(ctx)

	if auth != "" {
		headers["Authorization"] = "Bearer " + auth
	}

	return headers
}
