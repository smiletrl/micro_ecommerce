OURPATH = $(GOPATH)/src/github.com/smiletrl/micro_ecommerce

run-test:
	- STAGE=local go test ./...

db-start:
	- cd build && docker-compose up -d

db-restart:
	- cd build && docker-compose down && docker-compose up -d

db-stop:
	- docker-compose down

db-reset:
	- soda reset

build-cart:
	- docker build -t micro_ecommerce/service_cart:dev . -f service.cart/Dockerfile
	- docker login
	- docker tag micro_ecommerce/service_cart:dev smiletrl/micro_ecommerce_cart:dev
	- docker push smiletrl/micro_ecommerce_cart:dev

deploy-cart: build-cart
	- kubectl rollout restart deployment/cart --namespace=dev

build-customer:
	- docker build -t micro_ecommerce/service_customer:dev . -f service.customer/Dockerfile

build-product:
	- docker build -t micro_ecommerce/service_product:dev . -f service.product/Dockerfile

local-build: build-cart build-customer build-product
