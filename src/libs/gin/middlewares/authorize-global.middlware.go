package middlewares

import (
	"donation-mgmt/src/apperrors"
	"donation-mgmt/src/libs/db"
	"donation-mgmt/src/permissions"
	"donation-mgmt/src/system/contextual"

	"github.com/gin-gonic/gin"
)

func WithGlobalAuthorization(entityType permissions.Entity, capabilities ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		subject := contextual.GetSubject(c)
		if subject == "" {
			_ = c.Error(&apperrors.AuthorizationError{
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

		// TODO: Consider making this Unit of Work global and sharing it across all requests
		uow := db.NewUnitOfWork()
		defer uow.Finalize(c)

		querier, err := uow.GetQuerier(c)
		if err != nil {
			_ = c.Error(err)
			c.Abort()
			return
		}

		canDo, err := permissions.GetPermissionsService().HasCapabilities(c, querier, params)
		if err != nil {
			_ = c.Error(err)
			c.Abort()
			return
		}

		if !canDo {
			_ = c.Error(&apperrors.OperationForbiddenError{
				EntityID: apperrors.EntityIdentifier{
					EntityType: string(entityType),
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
