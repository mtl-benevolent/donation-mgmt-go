package apperrors

import (
	"fmt"
	"log/slog"
	"net/http"
)

type EntityNotFoundError struct {
	EntityID EntityIdentifier
}

func (e *EntityNotFoundError) Error() string {
	return fmt.Sprintf("%s was not found", e.EntityID.String())
}

func (e *EntityNotFoundError) ToRFC7807Error() RFC7807Error {
	return RFC7807Error{
		Type:     "NotFound",
		Title:    "Entity not found",
		Status:   http.StatusNotFound,
		Detail:   fmt.Sprintf("%s entity was not found", e.EntityID.EntityType),
		Instance: "",
	}
}

func (e *EntityNotFoundError) Log(l *slog.Logger) {
	l.Error("entity was not found", e.EntityID.LoggableFields()...)
}
