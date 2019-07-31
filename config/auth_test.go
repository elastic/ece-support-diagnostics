package config

import (
	"testing"
)

func TestConfig_checkForPassword(t *testing.T) {

	type fields struct {
		Auth Auth
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name:   "plain username",
			fields: fields{Auth{User: "test_username:password"}},
			want:   "test_username",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				Auth: tt.fields.Auth,
			}
			c.checkForPassword()
			if c.User != tt.want {
				t.Errorf("fp() = %v, want %v", c.Auth, tt.want)
			}
		})
	}
}
