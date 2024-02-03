export GOOSE_DRIVER=postgres
export GOOSE_DBSTRING="host=localhost port=5432 user=postgres password=postgresql dbname=ipam sslmode=disable"
export APP_PORT=8080
export APP_DSN="host=localhost port=5432 user=postgres password=postgresql dbname=ipam sslmode=disable"
export APP_DBPOOLCOUNT=100
export APP_WORKERPOOLCOUNT=100

generate-main:
	@swagger generate server -f ./api/api.yml --server-package=internal/handler -A ipamService --exclude-main

# generate-wo-main:
	# @swagger generate server -f ./api/api.yml -t internal/handler --exclude-main

grpc-gen:
	@protoc --plugin=/root/go/bin/protoc-gen-go --plugin=/root/go/bin/protoc-gen-go-grpc --go-rpc_out=./ --go_out=./ protoc/notificationService.proto

migrate-up:
	@GOOSE_DRIVER=postgres GOOSE_DBSTRING="host=localhost port=5432 user=postgres password=postgresql dbname=ipam sslmode=disable" goose -dir=./migrations up

migrate-down:
	@GOOSE_DRIVER=postgres GOOSE_DBSTRING="host=localhost port=5432 user=postgres password=postgresql dbname=ipam sslmode=disable" goose -dir=./migrations down

add-user:
	@curl -X POST  http://127.0.0.1:8090/user -H "fingerprint: fp-123" -H 'Content-Type: application/json' -d '{"payload":{"phone":"79296771805","tags":["1","2"], "timezone":"UTC-3"}}'

test:
	@go test -cover pkg/hasher/*.go
	@go test -cover internal/cache/*.go
	@go test -cover internal/config/*.go
	@go test -cover internal/handlers/*.go
	@go test -cover internal/metrics/*.go
	@go test -cover cmd/*.go

run:
	@go run cmd/notify-service-server/main.go
