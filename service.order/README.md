## Order
Order subscribes to RocketMQ queue to update order status when a payment is successful.

See more at `func Consume() error` from `internal/order/message.go`.

It uses database postgresSQL. Except the message part, other order business logic is not implemented in this project yet.
