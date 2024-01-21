package pagination

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

const DefaultLimit = 25
const DefaultOffset = 0

type PaginationOptions struct {
	Offset int
	Limit  int
}

type PaginatedDTO[T any] struct {
	Results []T `json:"items"`
	Total   int `json:"total"`
	Offset  int `json:"offset"`
	Limit   int `json:"limit"`
}

type PaginatedResult[T any] struct {
	Results []T
	Total   int
}

func ParsePaginationOptions(c *gin.Context) PaginationOptions {
	rawOffset := c.DefaultQuery("offset", fmt.Sprintf("%d", DefaultOffset))
	offset, err := strconv.ParseUint(rawOffset, 10, 32)
	if err != nil {
		offset = DefaultOffset
	}

	rawLimit := c.DefaultQuery("limit", fmt.Sprintf("%d", DefaultLimit))
	limit, err := strconv.ParseUint(rawLimit, 10, 32)
	if err != nil {
		limit = DefaultLimit
	}

	return PaginationOptions{
		Offset: int(offset),
		Limit:  int(limit),
	}
}
