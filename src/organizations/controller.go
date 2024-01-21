package organizations

import (
	"donation-mgmt/src/data_access"
	"donation-mgmt/src/libs/db"
	"net/http"

	"github.com/gin-gonic/gin"
)

func registerRoutes(router *gin.Engine) {
	orgRouter := router.Group("/v1/organizations")

	orgRouter.GET("", GetOrganizationsV1)
	orgRouter.POST("", CreateOrganizationV1)
	orgRouter.GET("/:slug", GetOrganizationBySlugV1)
}

func GetOrganizationsV1(c *gin.Context) {
	orgs, err := GetOrgService().GetOrganizations(c)
	if err != nil {
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
	uow := db.GetUnitOfWorkFromCtx(c)
	uow.UseTransaction()

	reqDTO := CreateOrganizationRequest{}
	if err := c.ShouldBindJSON(&reqDTO); err != nil {
		c.Error(err)
		return
	}

	if err := reqDTO.Validate(); err != nil {
		c.Error(err)
		return
	}

	org, err := GetOrgService().CreateOrganization(c, data_access.InsertOrganizationParams{
		Name: reqDTO.Name,
		Slug: reqDTO.Slug,
	})

	if err != nil {
		c.Error(err)
		return
	}

	dto := mapOrgToDTO(org)
	c.JSON(http.StatusCreated, dto)
}

func GetOrganizationBySlugV1(c *gin.Context) {
	slug := c.Params.ByName("slug")

	org, err := GetOrgService().GetOrganizationBySlug(c, slug)
	if err != nil {
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
