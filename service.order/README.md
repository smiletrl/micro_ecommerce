## Order
Order subscribes to RocketMQ queue to update order status when a payment is successful.

See more at `func Consume() error` from `internal/order/message.go`.

It used database postgresSQL.
