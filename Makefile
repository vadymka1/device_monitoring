
.PHONY: all proto build run migrate-up migrate-down

all: proto build

proto:
	protoc \
	  --go_out=. \
	  --go-grpc_out=. \
	  --go_opt=paths=source_relative \
	  --go-grpc_opt=paths=source_relative \
	  internal/proto/monitoring.proto

build:
	go build -o bin/device-monitoring main.go

run: build
	./bin/device-monitoring

migrate-up:
	migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/monitor?sslmode=disable" up

migrate-down:
	migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/monitor?sslmode=disable" down
