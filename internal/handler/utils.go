package handler

import (
	"encoding/json"
	"net/http"
)

type Request http.Request

func bind(r *http.Request, dataModel any) error {
	if r == nil {
		return nil
	}

	err := json.NewDecoder(r.Body).Decode(dataModel)
	if err != nil {
		return err
	}

	return nil
}
