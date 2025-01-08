package ginutils

import "github.com/gin-gonic/gin"

func DeserializeJSON[T any](ctx *gin.Context) (T, error) {
	var target T
	if err := ctx.BindJSON(&target); err != nil {
		return target, err
	}

	return target, nil
}
