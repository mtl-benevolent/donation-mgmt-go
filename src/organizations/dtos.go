package organizations

import (
	"donation-mgmt/src/apperrors"
	"donation-mgmt/src/system/validation"
	"reflect"
	"time"

	ozzo "github.com/go-ozzo/ozzo-validation"
)

type OrganizationDTOV1 struct {
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	CreatedAt time.Time `json:"createdAt"`
}

type CreateOrganizationRequestV1 struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}

func (dto CreateOrganizationRequestV1) Validate() error {
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

type UpdateOrganizationRequestV1 struct {
	Name string `json:"name"`
}

func (dto UpdateOrganizationRequestV1) Validate() error {
	return ozzo.ValidateStruct(
		&dto,
		ozzo.Field(&dto.Name, ozzo.Required, ozzo.Length(1, 255)),
	)
}
