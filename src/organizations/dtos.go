package organizations

import "time"

type OrganizationDTO struct {
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	CreatedAt time.Time `json:"created_at"`
}
