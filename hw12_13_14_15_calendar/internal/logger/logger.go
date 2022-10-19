package logger

import (
	"errors"
	"fmt"
	"os"
	"syscall"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	levelDebug = "debug"
	levelInfo  = "info"
	levelWarn  = "warn"
	levelError = "error"
)

type ILogger interface {
	Info(string, ...interface{})
	Error(string, ...interface{})
	Debug(string, ...interface{})
	Warn(string, ...interface{})
	GetInstance() *zap.Logger
}

type Logger struct {
	instance *zap.Logger
}

func (l *Logger) GetInstance() *zap.Logger {
	return l.instance
}

func getLoggerLevel(level string) zapcore.Level {
	switch level {
	case levelInfo:
		return zap.InfoLevel
	case levelDebug:
		return zap.DebugLevel
	case levelError:
		return zap.ErrorLevel
	case levelWarn:
		return zap.WarnLevel
	}
	return 0
}

func New(level string) *Logger {
	// First, define our level-handling logic.
	levelEnabler := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= getLoggerLevel(level)
	})

	// High-priority output should also go to standard error, and low-priority
	// output should also go to standard out.
	consoleDebugging := zapcore.Lock(os.Stdout)

	consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())

	// Join the outputs, encoders, and level-handling functions into
	// zapcore.Cores, then tee the four cores together.
	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, consoleDebugging, levelEnabler),
	)

	logger := zap.New(core)
	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil && !errors.Is(err, syscall.ENOTTY) {
			fmt.Println("could not sync", err)
		}
	}(logger)

	return &Logger{logger}
}

func (l *Logger) Info(s string, i ...interface{}) {
	l.instance.Sugar().Infow(s, i...)
}

func (l *Logger) Error(s string, i ...interface{}) {
	l.instance.Sugar().Errorw(s, i...)
}

func (l *Logger) Debug(s string, i ...interface{}) {
	l.instance.Sugar().Debugw(s, i...)
}

func (l *Logger) Warn(s string, i ...interface{}) {
	l.instance.Sugar().Warnw(s, i...)
}
