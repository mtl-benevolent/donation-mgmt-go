package apperrors

import (
	"fmt"
	"log/slog"
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

func (e *AuthorizationError) Log(l *slog.Logger) {
	innerErrorDetails := ""
	if e.InnerError != nil {
		innerErrorDetails = e.InnerError.Error()
	}

	l.Error("could not authorize request", slog.String("message", e.Message), slog.String("inner_error", innerErrorDetails))
}
