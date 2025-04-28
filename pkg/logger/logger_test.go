package logger

import (
	"bytes"
	"context"
	"log/slog"
	"strings"
	"testing"

	"github.com/go-jedi/lingvogramm_backend/config"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	cfg := config.LoggerConfig{
		IsJSON:     true,
		AddSource:  false,
		Level:      "debug",
		SetFile:    false,
		FileName:   "logs/app.log",
		MaxSize:    10,
		MaxBackups: 3,
		MaxAge:     7,
	}

	_ = New(cfg)
}

func TestNewIsText(t *testing.T) {
	cfg := config.LoggerConfig{
		IsJSON: false,
	}

	_ = New(cfg)
}

func TestNewIsJSON(t *testing.T) {
	cfg := config.LoggerConfig{
		IsJSON: true,
	}

	_ = New(cfg)
}

func TestNewAddSource(t *testing.T) {
	cfg := config.LoggerConfig{
		AddSource: true,
	}

	_ = New(cfg)
}

func TestNewLevelEmpty(t *testing.T) {
	cfg := config.LoggerConfig{
		Level: "",
	}

	_ = New(cfg)
}

func TestNewSetFileAndFileNameEmpty(t *testing.T) {
	cfg := config.LoggerConfig{
		SetFile:    true,
		FileName:   "",
		MaxSize:    10,
		MaxBackups: 3,
		MaxAge:     7,
	}

	_ = New(cfg)
}

func TestNewSetFile(t *testing.T) {
	cfg := config.LoggerConfig{
		SetFile:    true,
		FileName:   "logs/app.log",
		MaxSize:    10,
		MaxBackups: 3,
		MaxAge:     7,
	}

	_ = New(cfg)
}

func TestLoggerMethods(t *testing.T) {
	var buf bytes.Buffer

	ho := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}
	h := slog.NewTextHandler(&buf, ho)
	logger := &Logger{
		Logger: slog.New(h),
	}

	tests := []struct {
		name     string
		logFunc  func()
		expected string
	}{
		{
			name: "Debug",
			logFunc: func() {
				logger.Debug("test debug message", "key", "value")
			},
			expected: "level=DEBUG msg=\"test debug message\" key=value",
		},
		{
			name: "Info",
			logFunc: func() {
				logger.Info("test info message", "key", "value")
			},
			expected: "level=INFO msg=\"test info message\" key=value",
		},
		{
			name: "Warn",
			logFunc: func() {
				logger.Warn("test warn message", "key", "value")
			},
			expected: "level=WARN msg=\"test warn message\" key=value",
		},
		{
			name: "Error",
			logFunc: func() {
				logger.Error("test error message", "key", "value")
			},
			expected: "level=ERROR msg=\"test error message\" key=value",
		},
		{
			name: "DebugContext",
			logFunc: func() {
				ctx := context.TODO()
				logger.DebugContext(ctx, "test debug context", "key", "value")
			},
			expected: "level=DEBUG msg=\"test debug context\" key=value",
		},
		{
			name: "InfoContext",
			logFunc: func() {
				ctx := context.TODO()
				logger.InfoContext(ctx, "test info context", "key", "value")
			},
			expected: "level=INFO msg=\"test info context\" key=value",
		},
		{
			name: "WarnContext",
			logFunc: func() {
				ctx := context.TODO()
				logger.WarnContext(ctx, "test warn context", "key", "value")
			},
			expected: "level=WARN msg=\"test warn context\" key=value",
		},
		{
			name: "ErrorContext",
			logFunc: func() {
				ctx := context.TODO()
				logger.ErrorContext(ctx, "test error context", "key", "value")
			},
			expected: "level=ERROR msg=\"test error context\" key=value",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			buf.Reset()
			test.logFunc()
			output := strings.TrimSpace(buf.String())
			assert.Contains(t, output, test.expected)
		})
	}
}

func TestLoggerLevels(t *testing.T) {
	var buf bytes.Buffer

	ho := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}
	h := slog.NewTextHandler(&buf, ho)
	logger := &Logger{
		Logger: slog.New(h),
	}

	logger.Debug("this should not appear")
	assert.Empty(t, buf.String())

	logger.Info("this should appear")
	assert.Contains(t, buf.String(), "level=INFO")
}
