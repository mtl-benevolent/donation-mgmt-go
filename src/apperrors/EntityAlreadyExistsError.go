package apperrors

import "fmt"

type EntityAlreadyExistsError struct {
	EntityName string
	EntityID   string
	Extra      map[string]any
}

func (e *EntityAlreadyExistsError) Error() string {
	withID := formatID(e.EntityID)
	extras := formatExtras(e.Extra)

	return fmt.Sprintf("%s entity %s%s already exists", e.EntityName, withID, extras)
}
