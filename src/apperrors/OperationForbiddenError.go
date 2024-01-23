package apperrors

import (
	"donation-mgmt/src/config"
	"fmt"
	"net/http"
)

type OperationForbiddenError struct {
	EntityName string
	EntityID   string
	Extra      map[string]any
	Capability string
}

func (e *OperationForbiddenError) Error() string {
	withID := formatID(e.EntityID)
	extras := formatExtras(e.Extra)

	return fmt.Sprintf("Forbidden to perform '%s' operation on %s entity %s%s.", e.Capability, e.EntityName, withID, extras)
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
