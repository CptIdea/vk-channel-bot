package cfg

import (
	"os"
	"reflect"
	"testing"
)

func TestParseEnv(t *testing.T) {
	tests := []struct {
		name      string
		envToken  string
		envAdmins string
		want      *Config
		wantErr   bool
	}{
		{
			"normal",
			"test",
			"123",
			&Config{
				Token:  "test",
				Admins: []int{123},
			},
			false,
		},
		{
			"no token",
			"",
			"123",
			nil,
			true,
		},
		{
			"no one",
			"",
			"",
			nil,
			true,
		},
		{
			"admin parsing error",
			"test",
			"test",
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := os.Setenv("TOKEN", tt.envToken)
			if err != nil {
				t.Fatal(err)
			}
			err = os.Setenv("ADMINS", tt.envAdmins)
			if err != nil {
				t.Fatal(err)
			}

			got, err := ParseEnv()
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseEnv() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseEnv() got = %v, want %v", got, tt.want)
			}
		})
	}
}
