package middlewares

import (
	"donation-mgmt/src/libs/db"

	"github.com/gin-gonic/gin"
)

func UnitOfWork(c *gin.Context) {
	unitOfWork := db.NewUnitOfWork()
	c.Set(db.UnitOfWorkCtxKey, unitOfWork)

	// We execute our API Handler
	c.Next()

	hasErrors := len(c.Errors) > 0

	err := unitOfWork.Finalize(c, !hasErrors)
	if err != nil {
		_ = c.Error(err)
		return
	}
}
