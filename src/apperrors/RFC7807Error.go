package apperrors

type RFC7807Error struct {
	Type     string `json:"type,omitempty"`
	Title    string `json:"title"`
	Status   int    `json:"status"`
	Detail   string `json:"detail"`
	Instance string `json:"instance,omitempty"`
}

type DetailedError interface {
	ToRFC7807Error() RFC7807Error
}
