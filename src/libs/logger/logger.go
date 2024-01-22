package logger

import (
	"donation-mgmt/src/config"
	"log/slog"
	"os"
)

var logger *slog.Logger

func BootstrapLogger(appConfig *config.AppConfiguration) *slog.Logger {
	levelVar := slog.LevelVar{}
	err := levelVar.UnmarshalText([]byte(appConfig.LogLevel))

	var leveler slog.Leveler = &levelVar
	if err != nil {
		leveler = slog.LevelInfo
	}

	// TODO: Come back at some point and format correctly to adhere
	// to GCP's logging format: https://cloud.google.com/logging/docs/structured-logging#structured_logging_special_fields
	logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     leveler,
		AddSource: appConfig.LogAddSource,
	}))

	if err != nil {
		logger.Error("Error initializing logger. Defaulting to INFO level", slog.Any("error", err))
	}

	return logger
}

func ForceSetLogger(newLogger *slog.Logger) {
	logger = newLogger
}

func Logger() *slog.Logger {
	if logger == nil {
		panic("Logger was not bootstrapped")
	}

	return logger
}

func ForComponent(component string) *slog.Logger {
	return Logger().With("component", component)
}
