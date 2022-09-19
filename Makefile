# Constants
GO_GEN_DIR = ./proto/gen/

# 'encoding' builds the protocol buffer's used to encode goscript programs
encoding:
	mkdir -p $(GO_GEN_DIR) 
	protoc ./proto/goscript.proto \
	--go_out=$(GO_GEN_DIR) \
	--proto_path=$(shell pwd)
	cp -R ./proto/gen/goscript/pkg/encoding ./src/pkg/

# 'gsc-prod' builds the goscript compiler in production mode (stripped)
gsc-prod:
	CGO_ENABLED=0 go build -o ./dist/gsc --trimpath --ldflags="-s -w" ./src/cmd/gsc/gsc.go

# 'gsc-dev' builds the goscript compiler in developement mode
gsc-dev:
	CGO_ENABLED=0 go build -o ./dist/gsc-dev ./src/cmd/gsc/gsc.go

# 'gsr-prod' build the goscript runtime in production mode (stripped)
gsr-prod:
	CGO_ENABLED=0 go build -o ./dist/gsr --trimpath --ldflags="-s -w" ./src/cmd/gsr/gsr.go

# 'gsr-dev' builds the goscript runtime in developement mode
gsr-dev:
	CGO_ENABLED=0 go build -o ./dist/gsr-dev ./src/cmd/gsc/gsr.go

# 'test' runs all unit tests defined in the goscript package which powers the compiler and the runtime
test:
	go test --test.v ./src/pkg/goscript/

# 'lint' runs static code analysis on the entire workspace according to our .golagci.yml
lint:
	golangci-lint run

# 'prod' builds a full distribution production release
prod: lint test encoding gsc-dev gsr-dev

# 'dev' compiles the entire distribution
dev: gsc-dev gsr-dev