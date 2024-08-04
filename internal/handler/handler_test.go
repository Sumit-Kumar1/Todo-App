package handler

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestHandler_AddTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	s := NewMockServicer(ctrl)
	h := New(s)

	tests := []struct {
		name     string
		method   string
		target   string
		body     io.Reader
		header   string
		respCode int
		respBody string
	}{
		{name: "valid req", method: http.MethodPost, target: "/add", header: "true", respCode: http.StatusOK},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(tt.method, tt.target, tt.body)
			r.Header.Add(hxRequest, tt.header)
			r.Form.Set("", "") //TODO: put the form values here

			w := httptest.NewRecorder()

			h.AddTask(w, r)

			assert.Equalf(t, tt.respCode, w.Result().StatusCode, "TEST[%d] Failed - %s", i, tt.name)

			w.Result().Body.Close()
		})
	}
}
