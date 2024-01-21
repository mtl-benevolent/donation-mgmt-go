package organizations

import (
	"donation-mgmt/src/apperrors"
	"donation-mgmt/src/system/validation"
	"reflect"
	"time"

	ozzo "github.com/go-ozzo/ozzo-validation"
)

type OrganizationDTO struct {
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	CreatedAt time.Time `json:"createdAt"`
}

type CreateOrganizationRequest struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}

func (dto CreateOrganizationRequest) Validate() error {
	err := ozzo.ValidateStruct(
		&dto,
		ozzo.Field(&dto.Name, ozzo.Required, ozzo.Length(1, 255)),
		ozzo.Field(&dto.Slug, ozzo.Required, ozzo.Length(1, 32), ozzo.Match(validation.SlugRegex)),
	)

	if err != nil {
		return &apperrors.ValidationError{
			EntityName: reflect.TypeOf(dto).Name(),
			InnerError: err,
		}
	}

	return nil
}
