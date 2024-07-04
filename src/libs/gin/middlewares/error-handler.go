package middlewares

import (
	"donation-mgmt/src/apperrors"
	"donation-mgmt/src/libs/logger"
	"donation-mgmt/src/system/contextual"
	"errors"
	"fmt"
	"io"
	"log/slog"

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
	l := logger.ForComponent("ErrorHandler")
	l = contextual.LoggerWithContextData(c, l)

	c.Next()

	if len(c.Errors) == 0 {
		return
	}

	reqErr := c.Errors.Last().Err

	var rfcErr apperrors.RFC7807Error

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

	if loggable, ok := reqErr.(apperrors.Loggable); ok {
		loggable.Log(l)
	} else {
		defaultErrorLogger(l, reqErr, rfcErr)
	}

	c.JSON(rfcErr.Status, rfcErr)
}

func defaultErrorLogger(l *slog.Logger, err error, rfcError apperrors.RFC7807Error) {
	l.Error(
		"An error occurred with the HTTP request",
		slog.String("error", err.Error()),
		slog.String("type", fmt.Sprintf("%T", err)),
		slog.String("title", rfcError.Title),
		slog.Int("status", rfcError.Status),
		slog.String("detail", rfcError.Detail),
	)
}
