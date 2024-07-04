package apperrors

import (
	"donation-mgmt/src/config"
	"fmt"
	"log/slog"
	"net/http"
)

type OperationForbiddenError struct {
	EntityID   EntityIdentifier
	Capability string
}

func (e *OperationForbiddenError) Error() string {
	return fmt.Sprintf("Forbidden to perform '%s' operation on %s", e.Capability, e.EntityID.String())
}

func (e *OperationForbiddenError) ToRFC7807Error() RFC7807Error {
	if config.AppConfig().RewriteForbiddenErrors {
		notFoundErr := &EntityNotFoundError{
			EntityID: e.EntityID,
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
	l.Error("operation forbidden", e.EntityID.LoggableFields()...)
}
