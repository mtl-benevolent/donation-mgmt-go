package organizations

import (
	"donation-mgmt/src/apperrors"
	"donation-mgmt/src/data_access"
	"donation-mgmt/src/libs/db"
	"donation-mgmt/src/pagination"
	p "donation-mgmt/src/permissions"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// TODO: Remove me and implement a proper middleware
const mockSubject = "ginette"

func registerRoutes(router *gin.Engine) {
	orgRouter := router.Group("/v1/organizations")

	orgRouter.GET("", ListOrganizationsV1)
	orgRouter.POST("", authorize(p.EntityOrganization.Capability(p.ActionCreate)), CreateOrganizationV1)
	orgRouter.GET("/:slug", authorize(p.EntityOrganization.Capability(p.ActionRead)), GetOrganizationBySlugV1)
	orgRouter.PUT("/:slug", authorize(p.EntityOrganization.Capability(p.ActionUpdate)), UpdateOrganizationV1)
}

func ListOrganizationsV1(c *gin.Context) {
	page := pagination.ParsePaginationOptions(c)

	hasGlobalOrgRead, err := p.GetPermissionsService().HasCapabilities(c, p.HasRequiredPermissionsParams{
		Subject:      mockSubject,
		Capabilities: []string{p.EntityOrganization.Capability(p.ActionRead)},
		MustBeGlobal: true,
	})
	if err != nil {
		c.Error(err)
		return
	}

	scopeQueryBySubject := mockSubject
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

func authorize(capability string) gin.HandlerFunc {
	return func(c *gin.Context) {
		slug, found := c.Params.Get("slug")

		var params p.HasRequiredPermissionsParams
		if !found {
			params = p.HasRequiredPermissionsParams{
				Subject:      mockSubject,
				Capabilities: []string{capability},
				MustBeGlobal: true,
			}
		} else {
			params = p.HasRequiredPermissionsParams{
				Subject:          mockSubject,
				Capabilities:     []string{capability},
				OrganizationSlug: slug,
			}
		}

		canDo, err := p.GetPermissionsService().HasCapabilities(c, params)
		if err != nil {
			c.Error(err)
			c.Abort()
			return
		}

		if !canDo {
			c.Error(&apperrors.OperationForbiddenError{
				EntityName: p.EntityOrganization.String(),
				EntityID:   fmt.Sprintf("%d", params.OrganizationID),
				Extra: map[string]any{
					"slug": params.OrganizationSlug,
				},
				Capability: capability,
			})

			c.Abort()
			return
		}

		c.Next()
	}
}
