OURPATH = $(GOPATH)/src/github.com/smiletrl/micro_ecommerce

run-staff:
	go run ./cmd/staff-admin/main.go

run-miniprogram:
	go run ./cmd/wechat-miniprogram/main.go

run-docs-staff:
	# http://localhost:8080/smiletrl/micro_ecommerce/staff/
	- swag init -g ./cmd/staff-admin/main.go -o ./cmd/staff-admin/docs/ --exclude *_wechat_miniprogram.go

run-docs-miniprogram:
	# http://localhost:8080/smiletrl/micro_ecommerce/miniprogram/
	- swag init -g ./cmd/wechat-miniprogram/main.go -o ./cmd/wechat-miniprogram/docs/ --exclude *_staff.go

run-test:
	- STAGE=local go test ./...

# Database setup and management for localhost & testing
POSTGRESQL_URL := 'postgres://postgres:password@localhost:5433/eshop_api?sslmode=disable'
db-create:
	- migrate -database ${POSTGRESQL_URL} -path ./migrations drop

db-start:
	- docker-compose up -d

db-restart:
	- cd build && docker-compose down && docker-compose up -d

db-stop:
	- docker-compose down

migrate-up:
	- migrate -database ${POSTGRESQL_URL} -path ./migrations up

migrate-down:
	- migrate -database ${POSTGRESQL_URL} -path ./migrations down 1

migrate-set5:
	- migrate -database ${POSTGRESQL_URL} -path ./migrations force 5

# Need to run this once on the DB to setup migrations before testing
migrate-reset:
	- migrate -database ${POSTGRESQL_URL} -path ./migrations drop -f