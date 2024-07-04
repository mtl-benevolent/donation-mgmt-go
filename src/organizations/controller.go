package organizations

import (
	"donation-mgmt/src/apperrors"
	"donation-mgmt/src/data_access"
	"donation-mgmt/src/libs/db"
	"donation-mgmt/src/libs/gin/middlewares"
	"donation-mgmt/src/pagination"
	p "donation-mgmt/src/permissions"
	"donation-mgmt/src/system/contextual"
	"net/http"

	"github.com/gin-gonic/gin"
)

func registerRoutes(router *gin.Engine) {
	orgRouter := router.Group("/v1/organizations")

	orgRouter.GET("", ListOrganizationsV1) // Permissions are handled as part of the query
	orgRouter.POST("", middlewares.WithGlobalAuthorization(p.Organization, p.Organization.Capability(p.Create)), CreateOrganizationV1)
	orgRouter.GET("/:slug", middlewares.WithOrgAuthorization("slug", p.Organization.Capability(p.Read)), GetOrganizationBySlugV1)
	orgRouter.PUT("/:slug", middlewares.WithOrgAuthorization("slug", p.Organization.Capability(p.Update)), UpdateOrganizationV1)
}

func ListOrganizationsV1(c *gin.Context) {
	subject := contextual.GetSubject(c)
	if subject == "" {
		c.Error(&apperrors.AuthorizationError{
			Message: "User is not authenticated",
		})
		return
	}

	page := pagination.ParsePaginationOptions(c)

	hasGlobalOrgRead, err := p.GetPermissionsService().HasCapabilities(c, p.HasRequiredPermissionsParams{
		Subject:      subject,
		Capabilities: []string{p.Organization.Capability(p.Read)},
		MustBeGlobal: true,
	})
	if err != nil {
		c.Error(err)
		return
	}

	scopeQueryBySubject := subject
	if hasGlobalOrgRead {
		scopeQueryBySubject = ""
	}

	results, err := GetOrgService().GetOrganizations(c, ListOrganizationsParams{
		Subject:     scopeQueryBySubject,
		PageOptions: page,
	})
	if err != nil {
		c.Error(err)
		return
	}

	resultDtos := make([]OrganizationDTOV1, len(results.Results))
	for i, org := range results.Results {
		resultDtos[i] = mapOrgToDTO(org)
	}

	dto := pagination.PaginatedDTO[OrganizationDTOV1]{
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

	reqDTO := CreateOrganizationRequestV1{}
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

	reqDTO := UpdateOrganizationRequestV1{}
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

func mapOrgToDTO(org data_access.Organization) OrganizationDTOV1 {
	return OrganizationDTOV1{
		Name:      org.Name,
		Slug:      org.Slug,
		CreatedAt: org.CreatedAt,
	}
}
