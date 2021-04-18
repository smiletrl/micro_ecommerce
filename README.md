### Reference
- [Go Standard Layout](https://github.com/golang-standards/project-layout)
- [Micro Service Concept](https://docs.microsoft.com/en-us/azure/service-fabric/service-fabric-overview-microservices)

### Technical Stack
- [Golang >= 1.14](https://golang.org/)
- [PostgreSQL](https://www.postgresql.org/)
- [Pgx](https://github.com/JackC/pgx)
- [gRPC](https://github.com/grpc/grpc-go)
- [MiniKube](https://minikube.sigs.k8s.io/docs/start/)
- [Kubernetes](https://kubernetes.io/)
- [Istio](https://istio.io/)
- [RocketMQ](https://rocketmq.apache.org/)

### Background
This repo shows a simple ecommerce app running based on micro service architecture. It could run in a local environment with miniKube for local development. The primary language is Golang for these micro services. Kubernetes & istio are used for service register & discovery.

For people who are new but interested in micro serice, this repo is supposed to provide basic usage for items:

1. Kubernetes & Istio for service register, discovery and more mico service features.
2. gRPC for sync communication between micro services to provide low latency internal service response.
3. RocketMQ for async communication to decouple services to focus on individual service development.
4. Terraform configuration for local development & deployment. It's supposed to be easily to be migrated to an online cloud environment, such as AWS.
5. Detailed documentation is on the way^

### Services
Services include

- [Customer]()
- Cart
- Product
- Order
- Payment

### Local Installment

### Project Structure

```
github.com/smiletrl/micro_service
|-- .github
|-- build
|-- config
|   |-- github.yaml
|   |-- k8s.yaml
|   |-- local.yaml
|   |-- prod.yaml
|-- infrastructure
|   |-- host.tf
|   |-- main.tf
|   |-- envs
|   |   |-- dev
|   |       |-- cluster.tf
|   |       |-- gateway.yaml
|   |       |-- main.tf
|   |       |-- terraform.tfvars.example
|   |       |-- variables.tf
|   |       |-- virtual_services.yaml
|-- migrations
|   |-- 2021xxxx.down.sql
|   |-- 2021xxxx.up.sql
|-- pkg
|   |-- cache
|   |-- config
|   |-- constants
|   |-- context
|   |-- entity
|   |-- errors
|   |-- healthcheck
|   |-- migration
|   |-- rocketmq
|   |-- test
|-- redis
|   |-- config
|   |   |-- redis.conf
|-- service.cart
|-- service.customer
|-- service.order
|-- service.payment
|-- service.product
|-- testdata
|-- vendor
|-- .gitignore
|-- database.yml
|-- go.mod
|-- go.sum
|-- Makefile
|-- README.md
```

## Note

This repository is still in progress. Keep an eye on it if you are interested ^
