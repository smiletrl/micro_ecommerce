## Payment
Payment produces message to RocketMQ queue to notify other services. It used database postgresSQL.

It includes following scenarios at this moment.
- tell order service to update order status.
- tell product service to decrease the sku stock.
- tell customer service to decrease the customer balance value if the payment is balance method.

See more at `WechatPayCallback(c echo.Context)` from `internal/payment/api.go`.