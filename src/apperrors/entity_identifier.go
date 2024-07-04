package apperrors

import (
	"fmt"
	"strings"
)

type EntityIdentifier struct {
	EntityType string
	IDField    string
	EntityID   string

	Extras map[string]any
}

func (e EntityIdentifier) String() string {
	return fmt.Sprintf("%s%s", e.formatID(), e.formatExtras())
}

func (e EntityIdentifier) formatID() string {
	if e.IDField == "" {
		return fmt.Sprintf("%s entity", e.EntityType)
	} else {
		return fmt.Sprintf("%s entity with %s \"%s\"", e.EntityType, e.IDField, e.EntityID)
	}
}

func (e EntityIdentifier) formatExtras() string {
	strB := strings.Builder{}
	for key, value := range e.Extras {
		strB.WriteString(fmt.Sprintf("%s=%v, ", key, value))
	}

	if strB.Len() == 0 {
		return ""
	}

	return fmt.Sprintf(" {%s}", strB.String())
}

func (e EntityIdentifier) LoggableFields() []any {
	return []any{
		"entity_type", e.EntityType,
		"id_field", e.IDField,
		"entity_id", e.EntityID,
		"extras", e.Extras,
	}
}
