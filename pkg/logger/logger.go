package logger

import (
	"context"
	"io"
	"log/slog"
	"os"

	"github.com/go-jedi/lingvogramm_backend/config"
	"github.com/natefinch/lumberjack"
)

const fileNameDefault = "logs/app.log"

// ILogger defines the interface for the logger.
//
//go:generate mockery --name=ILogger --output=mocks --case=underscore
type ILogger interface {
	Debug(msg string, args ...any)
	DebugContext(ctx context.Context, msg string, args ...any)
	Info(msg string, args ...any)
	InfoContext(ctx context.Context, msg string, args ...any)
	Warn(msg string, args ...any)
	WarnContext(ctx context.Context, msg string, args ...any)
	Error(msg string, args ...any)
	ErrorContext(ctx context.Context, msg string, args ...any)
}

// Logger is a wrapper around slog.Logger that implements ILogger.
type Logger struct {
	*slog.Logger
}

// Make sure Logger implements ILogger.
var _ ILogger = (*Logger)(nil)

// New creates a new Logger instance with the given configuration.
func New(cfg config.LoggerConfig) *Logger {
	ho := &slog.HandlerOptions{}

	levelMapping := map[string]slog.Level{
		"debug": slog.LevelDebug,
		"info":  slog.LevelInfo,
		"warn":  slog.LevelWarn,
		"error": slog.LevelError,
	}

	if v, ok := levelMapping[cfg.Level]; ok {
		ho.Level = v
	} else {
		ho.Level = slog.LevelInfo
	}

	ho.AddSource = cfg.AddSource

	var output io.Writer = os.Stdout
	if cfg.SetFile {
		if cfg.FileName == "" {
			cfg.FileName = fileNameDefault
		}

		output = io.MultiWriter(
			os.Stdout,
			&lumberjack.Logger{
				Filename:   cfg.FileName,
				MaxSize:    cfg.MaxSize,
				MaxBackups: cfg.MaxBackups,
				MaxAge:     cfg.MaxAge,
			},
		)
	}

	var h slog.Handler
	if cfg.IsJSON {
		h = slog.NewJSONHandler(output, ho)
	} else {
		h = slog.NewTextHandler(output, ho)
	}

	return &Logger{
		Logger: slog.New(h),
	}
}
