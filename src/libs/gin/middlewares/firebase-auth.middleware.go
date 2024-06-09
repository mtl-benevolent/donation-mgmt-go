package middlewares

import (
	"donation-mgmt/src/apperrors"
	firebaseadmin "donation-mgmt/src/libs/firebase-admin"
	"donation-mgmt/src/libs/logger"
	"donation-mgmt/src/system/contextual"
	"log/slog"
	"strings"

	"github.com/gin-gonic/gin"
)

func FirebaseAuthMiddleware() gin.HandlerFunc {
	firebaseAuth := firebaseadmin.AuthClient()

	return func(c *gin.Context) {
		l := logger.ForComponent("FirebaseAuthMiddleware")
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			err := &apperrors.AuthorizationError{
				Message: "No Authorization header provided",
			}

			c.Error(err)
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			err := &apperrors.AuthorizationError{
				Message: "Invalid Authorization header provided. Make sure you are sending a valid Bearer token",
			}

			c.Error(err)
			c.Abort()
			return
		}

		jwt := parts[1]
		token, err := firebaseAuth.VerifyIDToken(c, jwt)
		if err != nil {
			appErr := &apperrors.AuthorizationError{
				Message:    "could not verify token",
				InnerError: err,
			}

			c.Error(appErr)
			c.Abort()
			return
		}

		subject := token.Subject
		c.Set(string(contextual.SubjectCtxKey), subject)

		l.Info("Successfully authorized user using Firebase", slog.String("subject", subject))
		c.Next()
	}
}
