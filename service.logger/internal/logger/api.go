package logger

import (
	"github.com/labstack/echo/v4"
	"github.com/smiletrl/micro_ecommerce/pkg/auth"
	errorsd "github.com/smiletrl/micro_ecommerce/pkg/errors"
	"github.com/smiletrl/micro_ecommerce/pkg/jwt"
	"net/http"
)

// RegisterHandlers for handlers
func RegisterHandlers(r *echo.Group, service Service, jwt jwt.Service) {
	res := &resource{service}

	group := r.Group("/logger")
	group.Use(auth.CustomerMiddleware(jwt))

	group.POST("", res.Create)
}

type resource struct {
	service Service
}

type createRequest struct {
	Type      string `json:"type"`       // e.g, `click`, `view`, `apply`
	ProductID string `json:"product_id"` // maybe entity_id depending on the log model.
	Category  string `json:"category"`   // e.g, `/order/recommended/product/`
}

type createResponse struct {
	Data string `json:"data"`
}

func (r resource) Create(c echo.Context) error {
	req := new(createRequest)
	if err := c.Bind(req); err != nil {
		return errorsd.BadRequest(c, err)
	}

	// RPC call to service product to get the product sku title, price, stock
	err := r.service.Create(c, *req)
	if err != nil {
		return errorsd.Abort(c, err)
	}
	return c.JSON(http.StatusOK, createResponse{
		Data: "ok",
	})
}
