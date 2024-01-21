package middlewares

import (
	"donation-mgmt/src/apperrors"
	"errors"
	"fmt"
	"io"

	"github.com/gin-gonic/gin"
)

func PanicHandler(c *gin.Context, panicReason any) {
	var err error
	var ok bool

	if err, ok = panicReason.(error); !ok {
		err = fmt.Errorf("endpoint panicked: %v", panicReason)
	}

	c.Error(err)
}

func ErrorHandler(c *gin.Context) {
	c.Next()

	if len(c.Errors) == 0 {
		return
	}

	reqErr := c.Errors.Last().Err

	rfcErr := apperrors.RFC7807Error{}

	if errors.Is(reqErr, io.EOF) {
		rfcErr = apperrors.RFC7807Error{
			Status:   400,
			Title:    "Bad Request",
			Detail:   "Request body is required",
			Instance: "",
		}
	} else if detailedErr, ok := reqErr.(apperrors.DetailedError); ok {
		rfcErr = detailedErr.ToRFC7807Error()
	} else {
		rfcErr = apperrors.RFC7807Error{
			Status:   500,
			Title:    "Unknown error",
			Detail:   "An error occurred while performing the request",
			Instance: "",
		}
	}

	c.JSON(rfcErr.Status, rfcErr)
}
