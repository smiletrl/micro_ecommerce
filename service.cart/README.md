## Cart
Cart is a simple cart implementation, using database redis to provide high performance update (to be implemented).

It only inclues below fields at this moment. 

- customer_id
- sku_id
- quantity

It uses grpc to communicate with product service in following scenarios:
1. Retrieve product sku stock when user adds item to cart. This is to check the sku stock is available for the item quantity to be added to cart.
2. Retrieve product title & sku detail(image, price, stock, attributes) when get cart items for a customer.

See `ProductProxy` from `internal/cart/proxy.go`, and `type product struct` from `cmd/main.go`.

