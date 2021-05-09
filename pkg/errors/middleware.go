package errors

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"runtime"

	"go.uber.org/zap"
)

func Recover(logger *zap.SugaredLogger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			defer func() {
				if r := recover(); r != nil {
					err, ok := r.(error)
					if !ok {
						err = fmt.Errorf("%v", r)
					}
					stack := make([]byte, 4<<10)
					length := runtime.Stack(stack, true)
					msg := fmt.Sprintf("[PANIC RECOVER] %v %s\n", err, stack[:length])
					logger.Error(msg)
					c.Error(err)
				}
			}()
			return next(c)
		}
	}
}
