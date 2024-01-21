package organizations

import (
	"donation-mgmt/src/data_access"
	"donation-mgmt/src/libs/db"
	"donation-mgmt/src/pagination"
	"net/http"

	"github.com/gin-gonic/gin"
)

func registerRoutes(router *gin.Engine) {
	orgRouter := router.Group("/v1/organizations")

	orgRouter.GET("", GetOrganizationsV1)
	orgRouter.POST("", CreateOrganizationV1)
	orgRouter.GET("/:slug", GetOrganizationBySlugV1)
	orgRouter.PUT("/:slug", UpdateOrganizationV1)
}

func GetOrganizationsV1(c *gin.Context) {
	page := pagination.ParsePaginationOptions(c)

	results, err := GetOrgService().GetOrganizations(c, page)
	if err != nil {
		c.Error(err)
		return
	}

	resultDtos := make([]OrganizationDTO, len(results.Results))
	for i, org := range results.Results {
		resultDtos[i] = mapOrgToDTO(org)
	}

	dto := pagination.PaginatedDTO[OrganizationDTO]{
		Results: resultDtos,
		Total:   results.Total,
		Offset:  page.Offset,
		Limit:   page.Limit,
	}

	c.JSON(http.StatusOK, dto)
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

func UpdateOrganizationV1(c *gin.Context) {
	slug := c.Params.ByName("slug")

	reqDTO := UpdateOrganizationRequest{}
	if err := c.ShouldBindJSON(&reqDTO); err != nil {
		c.Error(err)
		return
	}

	if err := reqDTO.Validate(); err != nil {
		c.Error(err)
		return
	}

	org, err := GetOrgService().UpdateOrganization(c, data_access.UpdateOrganizationBySlugParams{
		Slug: slug,
		Name: reqDTO.Name,
	})

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
