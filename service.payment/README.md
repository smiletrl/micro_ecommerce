## Payment
Payment produces message to RocketMQ queue to notify other services. It uses database postgresSQL.

It includes following scenarios at this moment.
- tell order service to update order status.
- tell product service to decrease the sku stock.
- tell customer service to decrease the customer balance value if the payment is balance method.

See more at `WechatPayCallback(c echo.Context)` from `internal/payment/api.go`.

The RocketMQ consumers should be responsible for handling the [Consumption idempotence](https://partners-intl.aliyun.com/help/doc-detail/44397.htm).