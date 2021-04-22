package healthcheck

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

// RegisterHandlers get health check api
func RegisterHandlers(r *echo.Group) {
	r.GET("/health", func(c echo.Context) error {
		// @todo, maybe add external service connection ping, such as pg, redis
		return c.String(http.StatusOK, "ok")
	})
}
