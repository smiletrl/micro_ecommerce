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
	- docker tag micro_ecommerce/service_cart:dev docker.io/smiletrl/micro_ecommerce_cart:dev
	- docker push docker.io/smiletrl/micro_ecommerce_cart:dev

deploy-cart: build-cart
	- kubectl rollout restart deployment/cart --namespace=dev

build-customer:
	- docker build -t micro_ecommerce/service_customer:dev . -f service.customer/Dockerfile
	- docker login
	- docker tag micro_ecommerce/service_customer:dev docker.io/smiletrl/micro_ecommerce_customer:dev
	- docker push docker.io/smiletrl/micro_ecommerce_customer:dev

build-product:
	- docker build -t micro_ecommerce/service_product:dev . -f service.product/Dockerfile
	- docker login
	- docker tag micro_ecommerce/service_product:dev docker.io/smiletrl/micro_ecommerce_product:dev
	- docker push docker.io/smiletrl/micro_ecommerce_product:dev

local-build: build-cart build-customer build-product

terraform:
	- cd infrastructure && terraform init
	- cd infrastructure && terraform apply -target=null_resource.dashboard_download -auto-approve
	- cd infrastructure && terraform apply -auto-approve
	- cd infrastructure/envs/dev && terraform init
	- cd infrastructure/envs/dev && terraform apply -auto-approve

# Start kubernetes at https://minikube.sigs.k8s.io/docs/start/
# brew install minikube // minikube version: v1.19.0
# minikube start
# Download & install istio https://istio.io/latest/docs/setup/getting-started/

# Start dashboard
minikube-dashboard:
	- minikube dashboard

# Start tunnel, i.e, make istio ingressgateway loadbalancer working at IP `127.0.0.1`, instead of the gateway ip.
# see more at https://github.com/istio/istio.io/issues/9340
minikube-loadbalancer:
	- minikube tunnel

# gRPC compile command at directory `service.product/internal/rpc/proto`.
#protoc --go_out=. --go_opt=paths=source_relative \
#    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
#    proto/product.proto
