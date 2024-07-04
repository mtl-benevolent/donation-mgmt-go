package apperrors

import "log/slog"

type RFC7807Error struct {
	Type     string         `json:"type,omitempty"`
	Title    string         `json:"title"`
	Status   int            `json:"status"`
	Detail   string         `json:"detail"`
	Details  map[string]any `json:"details,omitempty"`
	Instance string         `json:"instance,omitempty"`
}

type DetailedError interface {
	ToRFC7807Error() RFC7807Error
}

type Loggable interface {
	Log(l *slog.Logger)
}
