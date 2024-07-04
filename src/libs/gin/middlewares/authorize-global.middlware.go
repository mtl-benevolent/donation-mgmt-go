package middlewares

import (
	"donation-mgmt/src/apperrors"
	"donation-mgmt/src/permissions"
	"donation-mgmt/src/system/contextual"

	"github.com/gin-gonic/gin"
)

func WithGlobalAuthorization(entityType string, capabilities ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		subject := contextual.GetSubject(c)
		if subject == "" {
			c.Error(&apperrors.AuthorizationError{
				Message: "User is not authenticated",
			})
			c.Abort()
			return
		}

		params := permissions.HasRequiredPermissionsParams{
			Subject:      subject,
			Capabilities: capabilities,
			MustBeGlobal: true,
		}

		canDo, err := permissions.GetPermissionsService().HasCapabilities(c, params)
		if err != nil {
			c.Error(err)
			c.Abort()
			return
		}

		if !canDo {
			c.Error(&apperrors.OperationForbiddenError{
				EntityID: apperrors.EntityIdentifier{
					EntityType: entityType,
					Extras: map[string]any{
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
