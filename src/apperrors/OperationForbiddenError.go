package apperrors

import (
	"donation-mgmt/src/config"
	"fmt"
	"log/slog"
	"net/http"
)

type OperationForbiddenError struct {
	EntityName string
	EntityID   string
	Extra      map[string]any
}

func (e *OperationForbiddenError) Error() string {
	withID := formatID(e.EntityID)
	extras := formatExtras(e.Extra)

	return fmt.Sprintf("Operation forbidden on %s entity %s%s.", e.EntityName, withID, extras)
}

func (e *OperationForbiddenError) ToRFC7807Error() RFC7807Error {
	if config.AppConfig().RewriteForbiddenErrors {
		notFoundErr := &EntityNotFoundError{
			EntityName: e.EntityName,
			EntityID:   e.EntityID,
			Extra:      e.Extra,
		}

		return notFoundErr.ToRFC7807Error()
	}

	return RFC7807Error{
		Type:     "Forbidden",
		Title:    "Operation forbidden",
		Status:   http.StatusForbidden,
		Detail:   e.Error(),
		Instance: "",
	}
}

func (e *OperationForbiddenError) Log(l *slog.Logger) {
	l.Error("operation forbidden", slog.String("entity_name", e.EntityName), slog.String("entity_id", e.EntityID), slog.Any("extra", e.Extra))
}
