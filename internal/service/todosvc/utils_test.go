package todosvc

import (
	"strings"
	"testing"

	"github.com/google/uuid"
)

func Test_generateID(t *testing.T) {
	uid := uuid.NewString()
	tests := []struct {
		name string
		want string
	}{
		{name: "valid ID", want: prefixTask + uid},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := generateID()

			if !strings.Contains(got, prefixTask) || len(got) != len(uid)+5 {
				t.Errorf("Test[%d] Failed - %s\nGot:\t%+v\nWant:\t%+v", i, tt.name, got, tt.want)
			}
		})
	}
}

func Test_validateTask(t *testing.T) {
	type args struct {
		id    string
		title string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateTask(tt.args.id, tt.args.title); (err != nil) != tt.wantErr {
				t.Errorf("validateTask() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_validateID(t *testing.T) {
	type args struct {
		id string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateID(tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("validateID() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
