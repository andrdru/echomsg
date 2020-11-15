package echomsg

import (
	"github.com/labstack/echo/v4"
)

type (
	Message interface {
		Return(c echo.Context) error

		// Deprecated: use SetErrorMap instead
		SetError(code int, message string)
		SetErrorMap(code int, message string, field string)
	}
)
