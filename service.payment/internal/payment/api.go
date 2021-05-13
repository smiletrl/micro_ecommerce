package payment

import (
	"context"
	mq "github.com/apache/rocketmq-client-go/v2"
	"github.com/labstack/echo/v4"
	"github.com/smiletrl/micro_ecommerce/pkg/constants"
	errorsd "github.com/smiletrl/micro_ecommerce/pkg/errors"
	"github.com/smiletrl/micro_ecommerce/pkg/rocketmq"

	wePayment "github.com/medivhzhan/weapp/payment"
	"net/http"
)

// RegisterHandlers for routes
func RegisterHandlers(r *echo.Group, rocketMQ rocketmq.Provider) {
	res := newResource(rocketMQ)

	productGroup := r.Group("/payment")

	productGroup.POST("", res.Create)
	productGroup.POST("/wechatpay/callback", res.WechatPayCallback)
}

type resource struct {
	rocketProducer mq.Producer
}

func newResource(rocketMQ rocketmq.Provider) resource {
	p, err := rocketMQ.CreateProducer(context.Background(), constants.RocketMQGroupPayment)
	if err != nil {
		panic(err)
	}
	r := resource{
		rocketProducer: p,
	}
	return r
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
	err := wePayment.HandlePaidNotify(w, req, func(ntf wePayment.PaidNotify) (bool, string) {
		// Pay successfully, notify other services via rocketMQ.

		// @todo move these tags to constants
		// order service will subscribe to this.
		message := rocketmq.CreateMessage(constants.RocketMQTopicPayment, constants.RocketMQTag("Pay Succeed||order"), "order_id:xxx")
		_, err := r.rocketProducer.SendSync(ctx, message)
		if err != nil {
			panic(err)
		}

		// product service will subscribe to this.
		// product will reduce the stock value.
		message = rocketmq.CreateMessage(constants.RocketMQTopicPayment, constants.RocketMQTag("Pay Succeed||product||sku||quantity"), "sku_id:xxx||quantity:xxx")
		_, err = r.rocketProducer.SendSync(ctx, message)
		if err != nil {
			panic(err)
		}

		// customer service will subscribe to this.
		// if this payment type/method is `balance`, customer's balance should be reduced.
		message = rocketmq.CreateMessage(constants.RocketMQTopicPayment, constants.RocketMQTag("Pay Succeed||method||customer||balance"), "customer_id:xxx||amount:xxx")
		_, err = r.rocketProducer.SendSync(ctx, message)
		if err != nil {
			panic(err)
		}
		return true, ""
	})

	if err != nil {
		return err
	}

	return nil
}
