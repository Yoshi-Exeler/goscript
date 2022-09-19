# Constants
GO_GEN_DIR = ./proto/gen/

# 'encoding' builds the protocol buffer's used to encode goscript programs
encoding:
	mkdir -p $(GO_GEN_DIR) 
	protoc ./proto/goscript.proto \
	--go_out=$(GO_GEN_DIR) \
	--proto_path=$(shell pwd)
	cp -R ./proto/gen/goscript/pkg/encoding ./src/pkg/

gsc-prod:
	CGO_ENABLED=0 go build -o ./dist/gsc --trimpath --ldflags="-s -w" ./src/cmd/gsc/gsc.go

gsc-dev:
	CGO_ENABLED=0 go build -o ./dist/gsc-dev ./src/cmd/gsc/gsc.go

gsr-prod:
	CGO_ENABLED=0 go build -o ./dist/gsr --trimpath --ldflags="-s -w" ./src/cmd/gsr/gsr.go

gsr-dev:
	CGO_ENABLED=0 go build -o ./dist/gsr-dev ./src/cmd/gsc/gsr.go

test:
	go test --test.v ./src/pkg/goscript/

lint:
	golangci-lint run

prod: lint test encoding gsc-dev gsr-dev

dev: gsc-dev gsr-dev