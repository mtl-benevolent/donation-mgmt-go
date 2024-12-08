package middlewares

import (
	"donation-mgmt/src/apperrors"
	"donation-mgmt/src/libs/db"
	"donation-mgmt/src/permissions"
	"donation-mgmt/src/system/contextual"
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
)

var (
	ErrMissingOrgSlug = errors.New("missing org slug parameter")
)

func WithOrgAuthorization(orgSlugParam string, capabilities ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		subject := contextual.GetSubject(c)
		if subject == "" {
			c.Error(&apperrors.AuthorizationError{
				Message: "User is not authenticated",
			})
			c.Abort()
			return
		}

		slug, hasOrgSlug := c.Params.Get(orgSlugParam)
		if !hasOrgSlug {
			c.Error(ErrMissingOrgSlug)
			c.Abort()
			return
		}

		params := permissions.HasRequiredPermissionsParams{
			Subject:          subject,
			Capabilities:     capabilities,
			OrganizationSlug: slug,
		}

		uow := db.NewUnitOfWork()
		defer uow.Finalize(c)

		querier, err := uow.GetQuerier(c)
		if err != nil {
			c.Error(err)
			c.Abort()
			return
		}

		canDo, err := permissions.GetPermissionsService().HasCapabilities(c, querier, params)
		if err != nil {
			c.Error(err)
			c.Abort()
			return
		}

		if !canDo {
			c.Error(&apperrors.OperationForbiddenError{
				EntityID: apperrors.EntityIdentifier{
					EntityType: permissions.Organization.String(),
					IDField:    "id",
					EntityID:   fmt.Sprintf("%d", params.OrganizationID),
					Extras: map[string]any{
						"slug":         params.OrganizationSlug,
						"capabilities": capabilities,
					},
				},
			})

			c.Abort()
			return
		}

		c.Next()
	}
}
