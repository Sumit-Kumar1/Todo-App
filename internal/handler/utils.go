package handler

import (
	"encoding/json"
	"io"
	"net/http"
)

type Request http.Request

func bind(r *http.Request, dataModel any) error {
	body, err := r.GetBody()
	if err != nil {
		return err
	}

	data, err := io.ReadAll(body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, dataModel)
	if err != nil {
		return err
	}

	return nil
}
