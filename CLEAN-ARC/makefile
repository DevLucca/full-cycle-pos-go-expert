.PHONY: wire grpc gql run build run-bin

wire:
	cd cmd/ordersystem && wire

grpc:
	protoc --go_out=. \
    --go-grpc_out=. \
    internal/infra/grpc/protofiles/order.proto

gql:
	go run github.com/99designs/gqlgen generate

run:
	go run -C cmd/ordersystem main.go wire_gen.go

build:
	go build -C cmd/ordersystem -o server -ldflags="-w -s" main.go wire_gen.go

run-bin: build
	cd cmd/ordersystem && ./server