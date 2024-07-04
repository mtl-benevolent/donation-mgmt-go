package organizations

import (
	"donation-mgmt/src/apperrors"
	"donation-mgmt/src/data_access"
	"donation-mgmt/src/libs/db"
	"donation-mgmt/src/libs/gin/ginutils"
	"donation-mgmt/src/libs/gin/middlewares"
	"donation-mgmt/src/pagination"
	p "donation-mgmt/src/permissions"
	"donation-mgmt/src/system/contextual"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	orgSlugParam = "orgSlug"
)

type ControllerV1 struct {
	permissionsService *p.PermissionsService
	orgService         *OrganizationService
}

func NewControllerV1() *ControllerV1 {
	return &ControllerV1{
		permissionsService: p.GetPermissionsService(),
		orgService:         GetOrgService(),
	}
}

func (c *ControllerV1) RegisterRoutes(router *gin.Engine) {
	orgRouter := router.Group("/v1/organizations")

	orgCreate := p.Organization.Capability(p.Create)
	orgRead := p.Organization.Capability(p.Read)
	orgUpdate := p.Organization.Capability(p.Update)

	orgRouter.GET("", c.ListOrganizationsV1)
	orgRouter.POST("", middlewares.WithGlobalAuthorization(p.Organization.String(), orgCreate), c.CreateOrganizationV1)
	orgRouter.GET("/:orgSlug", middlewares.WithOrgAuthorization(orgSlugParam, orgRead), c.GetOrganizationBySlugV1)
	orgRouter.PUT("/:orgSlug", middlewares.WithOrgAuthorization(orgSlugParam, orgUpdate), c.UpdateOrganizationV1)
}

func (c *ControllerV1) ListOrganizationsV1(ctx *gin.Context) {
	subject := contextual.GetSubject(ctx)
	if subject == "" {
		ctx.Error(&apperrors.AuthorizationError{
			Message: "User is not authenticated",
		})
		return
	}

	page := pagination.ParsePaginationOptions(ctx)

	hasGlobalOrgRead, err := p.GetPermissionsService().HasCapabilities(ctx, p.HasRequiredPermissionsParams{
		Subject:      subject,
		Capabilities: []string{p.Organization.Capability(p.Read)},
		MustBeGlobal: true,
	})
	if err != nil {
		ctx.Error(err)
		return
	}

	scopeQueryBySubject := subject
	if hasGlobalOrgRead {
		scopeQueryBySubject = ""
	}

	results, err := GetOrgService().GetOrganizations(ctx, ListOrganizationsParams{
		Subject:     scopeQueryBySubject,
		PageOptions: page,
	})
	if err != nil {
		ctx.Error(err)
		return
	}

	resultDtos := make([]OrganizationDTOV1, len(results.Results))
	for i, org := range results.Results {
		resultDtos[i] = mapOrgToDTOV1(org)
	}

	dto := pagination.PaginatedDTO[OrganizationDTOV1]{
		Results: resultDtos,
		Total:   results.Total,
		Offset:  page.Offset,
		Limit:   page.Limit,
	}

	ctx.JSON(http.StatusOK, dto)
}

func (c *ControllerV1) CreateOrganizationV1(ctx *gin.Context) {
	uow := db.GetUnitOfWorkFromCtx(ctx)
	uow.UseTransaction()

	reqDTO, err := ginutils.DeserializeJSON[CreateOrganizationRequestV1](ctx)
	if err != nil {
		ctx.Error(err)
		return
	}

	if err := reqDTO.Validate(); err != nil {
		ctx.Error(err)
		return
	}

	org, err := GetOrgService().CreateOrganization(ctx, data_access.InsertOrganizationParams{
		Name: reqDTO.Name,
		Slug: reqDTO.Slug,
	})

	if err != nil {
		ctx.Error(err)
		return
	}

	dto := mapOrgToDTOV1(org)
	ctx.JSON(http.StatusCreated, dto)
}

func (c *ControllerV1) GetOrganizationBySlugV1(ctx *gin.Context) {
	slug := ctx.Params.ByName(orgSlugParam)

	org, err := c.orgService.GetOrganizationBySlug(ctx, slug)
	if err != nil {
		ctx.Error(err)
		return
	}

	dto := mapOrgToDTOV1(org)

	ctx.JSON(http.StatusOK, dto)
}

func (c *ControllerV1) UpdateOrganizationV1(ctx *gin.Context) {
	slug := ctx.Params.ByName(orgSlugParam)

	reqDTO, err := ginutils.DeserializeJSON[UpdateOrganizationRequestV1](ctx)
	if err != nil {
		ctx.Error(err)
		return
	}

	if err := reqDTO.Validate(); err != nil {
		ctx.Error(err)
		return
	}

	org, err := c.orgService.UpdateOrganization(ctx, data_access.UpdateOrganizationBySlugParams{
		Slug: slug,
		Name: reqDTO.Name,
	})

	if err != nil {
		ctx.Error(err)
		return
	}

	dto := mapOrgToDTOV1(org)

	ctx.JSON(http.StatusOK, dto)
}

func mapOrgToDTOV1(org data_access.Organization) OrganizationDTOV1 {
	return OrganizationDTOV1{
		Name:      org.Name,
		Slug:      org.Slug,
		CreatedAt: org.CreatedAt,
	}
}
