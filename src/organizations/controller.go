package organizations

import (
	"donation-mgmt/src/data_access"
	"donation-mgmt/src/libs/db"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

func RegisterRoutes(router *gin.Engine) {
	orgRouter := router.Group("/v1/organizations")

	orgRouter.GET("", GetOrganizationsV1)
	orgRouter.POST("", CreateOrganizationV1)
	orgRouter.GET("/:slug", GetOrganizationBySlugV1)
}

func GetOrganizationsV1(c *gin.Context) {
	db := data_access.New(db.DBPool())

	orgs, err := db.GetOrganizations(c)
	if err != nil {
		if err == pgx.ErrNoRows {
			c.JSON(http.StatusOK, []string{})
			return
		}

		c.Error(err)
		return
	}

	dtos := make([]OrganizationDTO, len(orgs))
	for i, org := range orgs {
		dtos[i] = mapOrgToDTO(org)
	}

	c.JSON(http.StatusOK, dtos)
}

func CreateOrganizationV1(c *gin.Context) {
	c.String(http.StatusNotFound, "Not implemented")
}

func GetOrganizationBySlugV1(c *gin.Context) {
	slug := c.Params.ByName("slug")

	db := data_access.New(db.DBPool())

	org, err := db.GetOrganizationBySlug(c, slug)
	if err != nil {
		if err == pgx.ErrNoRows {
			c.Status(http.StatusNotFound)
			return
		}

		c.Error(err)
		return
	}

	dto := mapOrgToDTO(org)

	c.JSON(http.StatusOK, dto)
}

func mapOrgToDTO(org data_access.Organization) OrganizationDTO {
	return OrganizationDTO{
		Name:      org.Name,
		Slug:      org.Slug,
		CreatedAt: org.CreatedAt,
	}
}
