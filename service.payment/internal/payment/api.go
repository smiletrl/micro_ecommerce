package payment

import (
	"github.com/labstack/echo/v4"
	errorsd "github.com/smiletrl/micro_ecommerce/pkg/errors"

	"net/http"
)

// RegisterHandlers for routes
func RegisterHandlers(r *echo.Group, service Service) {
	res := resource{service}

	productGroup := r.Group("/payment")

	productGroup.POST("", res.Create)
	productGroup.POST("/wechatpay/callback", res.WechatPayCallback)
}

type resource struct {
	service Service
}

type createRequest struct {
	Type       string `json:"type"`
	OrderID    int64  `json:"order_id"`
	Amount     int    `json:"amount"`
	CustomerID int64  `json:"customer_id"`
}

type createResponse struct {
	Data string `json:"data"`
}

func (r resource) Create(c echo.Context) error {
	req := new(createRequest)
	if err := c.Bind(req); err != nil {
		return errorsd.BadRequest(c, err)
	}
	return c.JSON(http.StatusOK, createResponse{
		Data: "ok",
	})
}

func (r resource) WechatPayCallback(c echo.Context) error {
	ctx := c.Request().Context()
	w := c.Response().Writer
	req := c.Request()

	if err := r.service.PaySucceed(ctx, w, req); err != nil {
		return errorsd.Abort(c, err)
	}

	return c.String(200, "ok")
}
