package handler

import (
	"encoding/json"
	"net/http"
)

func WriteResponse[T any](w http.ResponseWriter, status int, data *T) error {
	if data == nil {
		w.WriteHeader(status)
		return nil
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	wrapData := struct {
		Data T `json:"data"`
	}{
		Data: *data,
	}

	rawData, err := json.Marshal(wrapData)
	if err != nil {
		return err
	}

	_, err = w.Write(rawData)
	return err
}
