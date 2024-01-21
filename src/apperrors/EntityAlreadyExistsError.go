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

func (e *EntityAlreadyExistsError) ToRFC7807Error() RFC7807Error {
	return RFC7807Error{
		Type:     "AlreadyExists",
		Title:    "Entity already exists",
		Status:   409,
		Detail:   fmt.Sprintf("%s entity already exists", e.EntityName),
		Instance: "",
	}
}
