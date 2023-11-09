package organizations

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type OrganizationsController struct {
}

func RegisterOrganizationsController(router fiber.Router) {
	controller := &OrganizationsController{}

	group := router.Group("/api/v1/organizations")
	group.Post(":idOrSlug", controller.GetOrganizationV1)
}

func (controller *OrganizationsController) GetOrganizationV1(c *fiber.Ctx) error {
	idOrSlug := c.Params("idOrSlug")

	id, err := uuid.Parse(idOrSlug)
	if err != nil {
		return err
	}

	org, err := FindOrganizationById(c.Context(), id)
	if err != nil {
		return err
	}

	c.JSON(org)
	return nil
}
