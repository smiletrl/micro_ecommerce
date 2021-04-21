package cart

import (
	"fmt"
	"github.com/labstack/echo/v4"
	errorsd "github.com/smiletrl/micro_ecommerce/pkg/errors"
	"net/http"
	"strconv"
)

// RegisterHandlers for handlers
func RegisterHandlers(r *echo.Group, service Service) {
	res := &resource{service}

	group := r.Group("/cart_item")

	group.GET("", res.Get)

	group.POST("", res.Create)

	group.DELETE("/:sku_id", res.Delete)
}

type resource struct {
	service Service
}

type createRequest struct {
	Quantity int   `json:"quantity"`
	SKUID    int64 `json:"sku_id"`
}

type createResponse struct {
	cartItem
}

func (r resource) Get(c echo.Context) error {
	return c.String(http.StatusOK, "succeed!")
}

func (r resource) Create(c echo.Context) error {
	req := new(createRequest)
	if err := c.Bind(req); err != nil {
		return errorsd.BadRequest(c, err)
	}
	fmt.Printf("req is: %+v\n", req)

	// @todo get customer id from middleware/context.
	customerID := int64(12)

	// RPC call to service product to get the product sku title, price, stock
	cart, err := r.service.Create(c, customerID, req.SKUID, req.Quantity)
	if err != nil {
		return errorsd.Abort(c, err)
	}
	return c.JSON(http.StatusOK, createResponse{
		cartItem: cart,
	})
}

type deleteResponse struct {
	Data string `json:"data"`
}

func (r resource) Delete(c echo.Context) error {
	id := c.Param("sku_id")
	idInt64, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return errorsd.BadRequest(c, err)
	}

	err = r.service.Delete(c, idInt64)
	if err != nil {
		return errorsd.Abort(c, err)
	}
	return c.JSON(http.StatusOK, deleteResponse{
		Data: "ok",
	})
}
