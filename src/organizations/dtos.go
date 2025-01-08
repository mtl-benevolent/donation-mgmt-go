package organizations

import (
	"donation-mgmt/src/apperrors"
	"donation-mgmt/src/system/validation"
	"fmt"
	"reflect"
	"strings"
	"time"

	ozzo "github.com/go-ozzo/ozzo-validation/v4"
)

type OrganizationDTOV1 struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	TimeZone  string    `json:"timezone"`
	CreatedAt time.Time `json:"createdAt"`
}

type CreateOrganizationRequestV1 struct {
	Name     string `json:"name"`
	Slug     string `json:"slug"`
	TimeZone string `json:"timezone"`
}

func (dto CreateOrganizationRequestV1) Validate() error {
	err := ozzo.ValidateStruct(
		&dto,
		ozzo.Field(&dto.Name, ozzo.Required, ozzo.Length(1, 255)),
		ozzo.Field(&dto.Slug, ozzo.Required, ozzo.Length(1, 32), ozzo.Match(validation.SlugRegex)),
		ozzo.Field(&dto.TimeZone, ozzo.Required, ozzo.By(validateTimezone)),
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
	Name     string `json:"name"`
	Timezone string `json:"timezone"`
}

func (dto UpdateOrganizationRequestV1) Validate() error {
	return ozzo.ValidateStruct(
		&dto,
		ozzo.Field(&dto.Name, ozzo.Required, ozzo.Length(1, 255)),
		ozzo.Field(&dto.Timezone, ozzo.Required, ozzo.By(validateTimezone)),
	)
}

func validateTimezone(value any) error {
	v, ok := value.(string)
	if !ok {
		return fmt.Errorf("value is not a string")
	}

	if strings.ToLower(v) == "local" {
		return fmt.Errorf("local timezone is not supported. Specify an explicit timezone")
	}

	if _, err := time.LoadLocation(v); err != nil {
		return fmt.Errorf("invalid timezone. Make sure you use a IANA timezone name")
	}

	return nil
}
