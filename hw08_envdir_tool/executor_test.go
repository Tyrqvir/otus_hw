package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	type args struct {
		cmd []string
		env Environment
	}
	tests := []struct {
		name           string
		args           args
		wantEnvName    string
		wantEnvValue   string
		wantReturnCode int
	}{
		{
			name: "valid data",
			args: args{
				cmd: []string{"echo"},
				env: map[string]EnvValue{
					"HELLO": {
						Value:      "world",
						NeedRemove: false,
					},
				},
			},
			wantEnvName:    "HELLO",
			wantEnvValue:   "world",
			wantReturnCode: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotReturnCode := RunCmd(tt.args.cmd, tt.args.env)

			os.LookupEnv(tt.wantEnvName)
			require.Equal(t, os.Getenv(tt.wantEnvName), tt.wantEnvValue)

			if gotReturnCode != tt.wantReturnCode {
				t.Errorf("RunCmd() = %v, want %v", gotReturnCode, tt.wantReturnCode)
			}
		})
	}
}
