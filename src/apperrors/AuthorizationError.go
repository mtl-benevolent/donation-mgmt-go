package apperrors

import (
	"fmt"
	"net/http"
	"strings"
)

type AuthorizationError struct {
	Message    string
	InnerError error
}

func (e *AuthorizationError) Error() string {
	sBuilder := strings.Builder{}
	sBuilder.WriteString("Could not authorize request")

	if e.Message != "" {
		sBuilder.WriteString(fmt.Sprintf(": %s", e.Message))
	}

	if e.InnerError != nil {
		sBuilder.WriteString(fmt.Sprintf(": %s", e.InnerError.Error()))
	}

	return sBuilder.String()
}

func (e *AuthorizationError) ToRFC7807Error() RFC7807Error {
	return RFC7807Error{
		Type:     "Authorization",
		Title:    "Authorization error",
		Status:   http.StatusUnauthorized,
		Detail:   e.Error(),
		Instance: "",
	}
}
