package webapi

import (
	"github.com/Slyngshot-Team/packages/log"
	"github.com/labstack/echo/v4"
	echomw "github.com/labstack/echo/v4/middleware"
)

const defaultStackSize = 4 << 10 // 4KB

func NewRecovery() echo.MiddlewareFunc {
	return echomw.RecoverWithConfig(echomw.RecoverConfig{
		Skipper:           echomw.DefaultSkipper,
		StackSize:         defaultStackSize,
		DisableStackAll:   false,
		DisablePrintStack: false,
		LogLevel:          0,
		LogErrorFunc: func(c echo.Context, err error, stack []byte) error {
			log.Error(c.Request().Context(), err, "recovery", "stack_trace", string(stack))
			return err
		},
		DisableErrorHandler: false,
	})
}
