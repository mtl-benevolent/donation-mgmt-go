package apperrors

import (
	"fmt"
	"net/http"
)

type EntityAlreadyExistsError struct {
	EntityID EntityIdentifier
}

func (e *EntityAlreadyExistsError) Error() string {
	return fmt.Sprintf("%s already exists", e.EntityID.String())
}

func (e *EntityAlreadyExistsError) ToRFC7807Error() RFC7807Error {
	return RFC7807Error{
		Type:     "AlreadyExists",
		Title:    "Entity already exists",
		Status:   http.StatusConflict,
		Detail:   fmt.Sprintf("%s entity already exists", e.EntityID.EntityType),
		Instance: "",
	}
}
