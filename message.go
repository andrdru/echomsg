package echomsg

import (
	"github.com/labstack/echo/v4"
)

type (
	Message interface {
		Return(c echo.Context) error
		SetError(options ...Option)
		Raw() interface{}
	}
)
