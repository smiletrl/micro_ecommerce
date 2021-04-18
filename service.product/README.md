## Product
Product provides both Rest Server & gRPC Server. Rest server is for external http request. gRPC server is for internal service sync communication.

It uses database mongoDB to build product catalog (to be implemented).

- gRPC server is registered at `rpcserver.Register()` from `cmd/main.go`. See more at `internal/rpc`.
- Rest server is like other services.

Product subscribes to RocketMQ queue to decrease product stock when a payment is successful.

See more at `func Consume() error` from `internal/product/message.go`.
