package usersvc

import (
	"errors"
	"testing"
)

func Test_encryptedPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  error
	}{
		{name: "valid case", password: "abcd"},
		{name: "valid case 2", password: "abcd124@adgbalje"},
	}
	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := encryptedPassword(tt.password)

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Test[%d] Failed - %s\nGOT:\t%v\nWANT:\t%v", i, tt.name, err, tt.wantErr)
			}

			if tt.password == got {
				t.Errorf("Test[%d] Failed - %s\nGOT:%v\nWANT:%v", i, tt.name, got, tt.password)
			}
		})
	}
}
