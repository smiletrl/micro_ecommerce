## Product
This service is a simple implementation for product.

### DB schema
It uses database mongoDB to build product catalog. MongoDB is selected for its dynamic schema to provide great flexiblity for product skus/attributes.

One example is this app might have over 10 thousands types/categories of products. Each type will have its own specific fields/properties. In a SQL database, we might need to create 10 thousands customized tables to hold each product type's customized fields. With MongoDB, one collection/document is enough to hold all product types.

While the benefit of NO-SQL DB for this product service shines, the disadvantage is also pretty clear. Need to be careful with the unsupported transaction operations. The future potential of SQL relation requirement for products data may grow too. I have seen a few companies switched to NO-SQL, and then switched back to SQL again.

Product doc
```
{
    "_id": "123abc",
    "title": "Pretty mac",
    "body": "This mac was made in USA and it was for children."
    "category": "/computer/mac", // 
    "assets": {
        "ppt": [ // different asset types.
            {
                "src": "xxx.png"
            },
            {
                "src": "xxx.png"
            }
        ],
        "thumbnail": {
            "src": "xxx.png"
        }
    },
    "variants": { // holds the product variants configuration.
        "attrs": [ // each variant should have these attrs value.
            {
                "name": "Color"
            },
            {
                "name": "Memory"
            }
        ]
    },
    "createdAt": 123231212,
    "updatedAt": 131223482
}
```

Sku/Variant doc
```
// attr value can be manually added at edit screen, or be selected from predefined attr list.
// by making use of mongoDB, we want to allow user to manually add new attr, instead of selecting from a predefined list.
{
    "_id": "12231212",
    "productId": "123abc",
    "assets": {
        "ppt":[ // different asset types.
            {
                "src": "xxx.png"
            },
            {
                "src": "xxx.png"
            }
        ],
    },
    "attrs": [
        {
            "name": "Color",
            "value": "Red" // value could be added manually at edit screen.
        },
        {
            "name": "Memory",
            "value": "32GB"
        }
    ],
    "price": "8900", // actual value is 89.00. In real env, price could deserve its own doc.
    "stock": "69" // stock number for this sku. In real env, stock could deserve its own doc.
}
```

Category taxonomy doc
```
{
    "_id": "123bbb",
    "name": "computer",
    "parent": "1212",
}

{
    "_id": "123ccc",
    "name": "mac",
    "parent": "1212",
}
```

### Rest & gRPC server
Product provides both Rest Server & gRPC Server. Rest server is for external http request. gRPC server is for internal service sync communication.

- gRPC server is registered at `rpcserver.Register()` from `cmd/main.go`. See more at `internal/rpc`.
- Rest server is like other services.

Product subscribes to RocketMQ queue to decrease product stock when a payment is successful.

See more at `func Consume() error` from `internal/product/message.go`.
