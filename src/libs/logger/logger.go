package logger

import (
	"log/slog"
	"os"
)

var logger *slog.Logger

func BootstrapLogger() *slog.Logger {
	// TODO: Come back at some point and format correctly to adhere
	// to GCP's logging format: https://cloud.google.com/logging/docs/structured-logging#structured_logging_special_fields
	logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

	return logger
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
