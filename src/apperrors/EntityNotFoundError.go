package apperrors

import "fmt"

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
		Status:   404,
		Detail:   fmt.Sprintf("%s entity was not found", e.EntityName),
		Instance: "",
	}
}
