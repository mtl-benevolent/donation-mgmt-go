package organizations

import (
	"donation-mgmt/src/apperrors"
	"donation-mgmt/src/dal"
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
	uow := db.NewUnitOfWork()
	defer uow.Finalize(c)

	querier, err := uow.GetQuerier(c)
	if err != nil {
		_ = c.Error(err)
		return
	}

	subject := contextual.GetSubject(c)
	if subject == "" {
		_ = c.Error(&apperrors.AuthorizationError{
			Message: "User is not authenticated",
		})
		return
	}

	page := pagination.ParsePaginationOptions(c)

	hasGlobalOrgRead, err := p.GetPermissionsService().HasCapabilities(c, querier, p.HasRequiredPermissionsParams{
		Subject:      subject,
		Capabilities: []string{p.Organization.Capability(p.Read)},
		MustBeGlobal: true,
	})
	if err != nil {
		_ = c.Error(err)
		return
	}

	scopeQueryBySubject := subject
	if hasGlobalOrgRead {
		scopeQueryBySubject = ""
	}

	results, err := GetOrgService().GetOrganizations(c, querier, ListOrganizationsParams{
		Subject:     scopeQueryBySubject,
		PageOptions: page,
	})
	if err != nil {
		_ = c.Error(err)
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
	uow := db.NewUnitOfWork()
	defer uow.Finalize(c)

	querier, err := uow.GetQuerier(c)
	if err != nil {
		_ = c.Error(err)
		return
	}

	reqDTO := CreateOrganizationRequestV1{}
	if err := c.ShouldBindJSON(&reqDTO); err != nil {
		_ = c.Error(err)
		return
	}

	if err := reqDTO.Validate(); err != nil {
		_ = c.Error(err)
		return
	}

	org, err := GetOrgService().CreateOrganization(c, querier, dal.InsertOrganizationParams{
		Name:     reqDTO.Name,
		Slug:     reqDTO.Slug,
		TimeZone: reqDTO.TimeZone,
	})

	if err != nil {
		_ = c.Error(err)
		return
	}

	if err = uow.Commit(c); err != nil {
		_ = c.Error(err)
		return
	}

	dto := mapOrgToDTO(org)
	c.JSON(http.StatusCreated, dto)
}

func GetOrganizationBySlugV1(c *gin.Context) {
	uow := db.NewUnitOfWork()
	defer uow.Finalize(c)

	querier, err := uow.GetQuerier(c)
	if err != nil {
		_ = c.Error(err)
		return
	}

	slug := c.Params.ByName("slug")

	org, err := GetOrgService().GetOrganizationBySlug(c, querier, slug)
	if err != nil {
		_ = c.Error(err)
		return
	}

	dto := mapOrgToDTO(org)

	c.JSON(http.StatusOK, dto)
}

func UpdateOrganizationV1(c *gin.Context) {
	uow := db.NewUnitOfWorkWithTx()
	defer uow.Finalize(c)

	querier, err := uow.GetQuerier(c)
	if err != nil {
		_ = c.Error(err)
		return
	}

	slug := c.Params.ByName("slug")

	reqDTO := UpdateOrganizationRequestV1{}
	if err := c.ShouldBindJSON(&reqDTO); err != nil {
		_ = c.Error(err)
		return
	}

	if err := reqDTO.Validate(); err != nil {
		_ = c.Error(err)
		return
	}

	org, err := GetOrgService().UpdateOrganization(c, querier, dal.UpdateOrganizationBySlugParams{
		Slug:     slug,
		Name:     reqDTO.Name,
		TimeZone: reqDTO.Timezone,
	})

	if err != nil {
		_ = c.Error(err)
		return
	}

	if err = uow.Commit(c); err != nil {
		_ = c.Error(err)
		return
	}

	dto := mapOrgToDTO(org)

	c.JSON(http.StatusOK, dto)
}

func mapOrgToDTO(org dal.Organization) OrganizationDTOV1 {
	return OrganizationDTOV1{
		Name:      org.Name,
		Slug:      org.Slug,
		TimeZone:  org.Timezone,
		CreatedAt: org.CreatedAt,
	}
}
