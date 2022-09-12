# Constants
GO_GEN_DIR = ./gen/
PROTO_FILES = goscript.proto
PROTO_PATH = $(shell pwd)

proto:
	mkdir -p $(GO_GEN_DIR) 
	protoc $(PROTO_FILES) \
	--go_out=$(GO_GEN_DIR) \
	--proto_path=$(PROTO_PATH)
	cp -R ./gen/goscript/pkg/encoding ./pkg/
