## Cart
Cart is a simple cart implementation, using redis to provide high performance update. Based on experience, cart item change is pretty frequent operation. Redis is a better solution than SQL DB in this case.

The trade-off is redis needs additiontional configuration to enable data persistence. Even with persistence enabled, there's still chance that data could be lost.

Redis uses hash to store customer cart data

Hash key is the customer id with prefix `cart:`, `cart:{customer_id}`

Hash key's field is the product sku id, and key's field value is the sku quantity. It looks like `{sku_id}: {quantity}`.

Cart service uses grpc to communicate with product service in following scenarios:

1. Retrieve product sku stock when user adds item to cart. This is to check the sku stock is available for the item quantity to be added to cart.
2. Retrieve product title & sku detail(image, price, stock, attributes) when get cart items for a customer.

See `ProductProxy` from `internal/cart/proxy.go`, and `type product struct` from `cmd/main.go`.
