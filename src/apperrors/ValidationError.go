package apperrors

import (
	"fmt"
	"log/slog"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation"
)

type ValidationError struct {
	EntityName string
	InnerError error
}

func (e *ValidationError) Error() string {
	if e.InnerError != nil {
		return fmt.Sprintf("Error validating %s entity: %v", e.EntityName, e.InnerError)
	}

	return fmt.Sprintf("Error validating %s entity", e.EntityName)
}

func (e *ValidationError) Unwrap() error {
	return e.InnerError
}

func (e *ValidationError) ToRFC7807Error() RFC7807Error {
	fieldErrors := make(map[string]any)
	if e.InnerError != nil {
		if err, ok := e.InnerError.(validation.Errors); ok {
			for field, err := range err {
				fieldErrors[field] = err.Error()
			}
		}
	}

	return RFC7807Error{
		Type:     "ValidationError",
		Title:    "Validation Error",
		Status:   http.StatusBadRequest,
		Detail:   fmt.Sprintf("Error validating %s entity", e.EntityName),
		Details:  fieldErrors,
		Instance: "",
	}
}

func (e *ValidationError) Log(l *slog.Logger) {
	fieldErrors := make(map[string]any)
	if e.InnerError != nil {
		if err, ok := e.InnerError.(validation.Errors); ok {
			for field, err := range err {
				fieldErrors[field] = err.Error()
			}
		}
	}

	l.Warn("validation error", slog.String("entity_name", e.EntityName), slog.Any("field_errors", fieldErrors))
}
