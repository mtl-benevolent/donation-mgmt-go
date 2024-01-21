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
