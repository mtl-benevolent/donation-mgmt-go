package apperrors

import (
	"fmt"
	"log/slog"
	"net/http"
)

type EntityNotFoundError struct {
	EntityName string
	EntityID   string
	Extra      map[string]any
}

func (e *EntityNotFoundError) Error() string {
	withID := formatID(e.EntityID)
	extras := formatExtras(e.Extra)

	return fmt.Sprintf("%s entity %s%s was not found", e.EntityName, withID, extras)
}

func (e *EntityNotFoundError) ToRFC7807Error() RFC7807Error {
	return RFC7807Error{
		Type:     "NotFound",
		Title:    "Entity not found",
		Status:   http.StatusNotFound,
		Detail:   fmt.Sprintf("%s entity was not found", e.EntityName),
		Instance: "",
	}
}

func (e *EntityNotFoundError) Log(l *slog.Logger) {
	l.Error("entity was not found", slog.String("entity_name", e.EntityName), slog.String("entity_id", e.EntityID), slog.Any("extra", e.Extra))
}
