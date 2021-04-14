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

build-product:
	- docker build -t micro_ecommerce/service_product:dev . -f service.product/Dockerfile

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
minikube-dashboard:
	- minikube dashboard

minikube-loadbalancer:
	- minikube tunnel
