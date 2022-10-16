package logger

import (
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Test_getLoggerLevel(t *testing.T) {
	type args struct {
		level string
	}
	tests := []struct {
		name string
		args args
		want zapcore.Level
	}{
		{
			name: "With debug level",
			args: args{
				level: levelDebug,
			},
			want: zap.DebugLevel,
		},
		{
			name: "With info level",
			args: args{
				level: levelInfo,
			},
			want: zap.InfoLevel,
		},
		{
			name: "With error level",
			args: args{
				level: levelError,
			},
			want: zap.ErrorLevel,
		},
		{
			name: "With warn level",
			args: args{
				level: levelWarn,
			},
			want: zap.WarnLevel,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getLoggerLevel(tt.args.level); got != tt.want {
				t.Errorf("getLoggerLevel() = %v, want %v", got, tt.want)
			}
		})
	}
}
