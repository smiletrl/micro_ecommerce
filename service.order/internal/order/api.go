package order

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/smiletrl/micro_ecommerce/pkg/constants"
	errorsd "github.com/smiletrl/micro_ecommerce/pkg/errors"
)

// RegisterHandlers for customer
func RegisterHandlers(r *echo.Group, service Service) {
	res := &resource{service}

	group := r.Group("/order")

	group.POST("/cart", res.CreateFromCart)
	group.POST("/product", res.CreateFromCart)
}

type resource struct {
	service Service
}

type createFromCartRequest struct {
	PaymentMethod string  `json:"payment_method"`
	CartItemIDs   []int64 `json:"cart_item_ids"`
}

type createResponse struct {
	ID int64 `json:"id"`
}

func (r resource) CreateFromCart(c echo.Context) error {
	customerID := c.Get(constants.AuthCustomerID).(int64)
	ctx := c.Request().Context()
	req := new(createFromCartRequest)
	if err := c.Bind(req); err != nil {
		return errorsd.BadRequest(c, err)
	}
	// @todo validate the logged in customer owns these cart items.
	id, err := r.service.CreateFromCart(ctx, customerID, *req)
	if err != nil {
		return errorsd.Abort(c, err)
	}
	return c.JSON(http.StatusOK, createResponse{
		ID: id,
	})
}

type createFromProductRequest struct {
	PaymentMethod string `json:"payment_method"`
	ProductID     string `json:"product_id"`
	Quantity      int    `json:"quantity"`
}

func (r resource) CreateFromProduct(c echo.Context) error {
	customerID := c.Get(constants.AuthCustomerID).(int64)
	ctx := c.Request().Context()
	req := new(createFromProductRequest)
	if err := c.Bind(req); err != nil {
		return errorsd.BadRequest(c, err)
	}
	// @todo validate the logged in customer owns these cart items.
	id, err := r.service.CreateFromProduct(ctx, customerID, *req)
	if err != nil {
		return errorsd.Abort(c, err)
	}
	return c.JSON(http.StatusOK, createResponse{
		ID: id,
	})
}
