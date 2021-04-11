package healthcheck

import (
	"github.com/labstack/echo/v4"
	"github.com/smiletrl/micro_ecommerce/pkg/dbcontext"
	"net/http"
)

// RegisterHandlers get health check api
func RegisterHandlers(r *echo.Group, db dbcontext.DB) {
	r.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})
}
