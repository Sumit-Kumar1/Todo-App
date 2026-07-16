package client

import (
	"context"
	"testing"

	"todoapp/internal/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_prepareAuthAPIHeaders(t *testing.T) {
	uid := uuid.NewString()
	ctx := context.WithValue(context.Background(), models.CorrelationID, uid)
	auth := "Bearer some-gibberish-value.token.id"

	tests := []struct {
		name string
		auth string
		want map[string]string
	}{
		{name: "valid headers", auth: "some-gibberish-value.token.id", want: map[string]string{"Authorization": auth, "X-Correlation-ID": uid}},
		{name: "no auth headers", auth: "", want: map[string]string{"X-Correlation-ID": uid}},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := prepareAuthAPIHeaders(ctx, tt.auth)

			assert.Equalf(t, tt.want, got, "Test failed - %s", i, tt.name)
		})
	}
}
