package organizations

import "github.com/google/uuid"

type Organization struct {
	ID      uuid.UUID `json:"id"`
	Name    string    `json:"name"`
	Slug    string    `json:"slug"`
	LogoUrl *string   `json:"logoUrl,omitempty"`
}
