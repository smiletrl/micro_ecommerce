package cart

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

	group := r.Group("/cart")
	group.Use(auth.CustomerMiddleware(jwt))

	group.GET("", res.Get)

	group.POST("", res.Create)

	group.DELETE("/:sku_id", res.Delete)
}

type resource struct {
	service Service
}

type createRequest struct {
	Quantity int    `json:"quantity"`
	SkuID    string `json:"sku_id"`
}

type createResponse struct {
	Data string `json:"data"`
}

type getResponse struct {
	Data []cartItem `json:"data"`
}

func (r resource) Get(c echo.Context) error {
	items, err := r.service.Get(c)
	if err != nil {
		return errorsd.Abort(c, err)
	}
	return c.JSON(http.StatusOK, getResponse{
		Data: items,
	})
}

func (r resource) Create(c echo.Context) error {
	req := new(createRequest)
	if err := c.Bind(req); err != nil {
		return errorsd.BadRequest(c, err)
	}

	// RPC call to service product to get the product sku title, price, stock
	err := r.service.Create(c, req.SkuID, req.Quantity)
	if err != nil {
		return errorsd.Abort(c, err)
	}
	return c.JSON(http.StatusOK, createResponse{
		Data: "ok",
	})
}

type deleteResponse struct {
	Data string `json:"data"`
}

func (r resource) Delete(c echo.Context) error {
	skuID := c.Param("sku_id")
	err := r.service.Delete(c, skuID)
	if err != nil {
		return errorsd.Abort(c, err)
	}
	return c.JSON(http.StatusOK, deleteResponse{
		Data: "ok",
	})
}
