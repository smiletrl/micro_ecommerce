## Customer
Customer includes internal `customer` & `balance` pkg. It uses database PostgresSQL.

- `customer` provides CRUD to customer entity
- `balance` provides customer's balance. Balance could be used to buy different products or other services. Customer could buy balance with other payment, such as wechat payment.

Customer `balance` subscribes to RocketMQ queue to reduce the balance value when a payment is paid with customer's balance method successfully.

See more at `internal/balance/message.go`.
