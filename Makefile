GOPATH:=$(shell go env GOPATH)
VERSION=$(shell git describe --tags --always)
INTERNAL_PROTO_DIR="./"
API_PROTO_DIR="./proto/api"
PKG_CONF_DIR="./proto/conf"
THIRD_PARTY_PROTO_DIR="./proto/third_party"
INTERNAL_PROTO_FILES=$(shell find $(INTERNAL_PROTO_DIR) -name "*.proto")
API_PROTO_FILES=$(shell find $(API_PROTO_DIR) -name "*.proto" -type f ! -name "error_reason.proto")
ERROR_PROTO_FILES=$(shell find $(API_PROTO_DIR) -name "error_reason.proto")

WIRE_FILES=$(shell find $(INTERNAL_PROTO_DIR) -name "wire.go")

.PHONY: init
# 初始化项目并下载和更新依赖项
init:
	git submodule init
	git submodule update
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/go-kratos/kratos/cmd/kratos/v2@latest
	go install github.com/go-kratos/kratos/cmd/protoc-gen-go-http/v2@latest
	go install github.com/go-kratos/kratos/cmd/protoc-gen-go-errors/v2@latest
	go install github.com/google/gnostic/cmd/protoc-gen-openapi@v0.6.1
	go install github.com/google/wire/cmd/wire@latest
	go install github.com/envoyproxy/protoc-gen-validate@latest
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
	go mod tidy

.PHONY: wire
# 依赖注入
wire:
	@echo '依赖注入...'
	@for file in $(WIRE_FILES) ; do \
		path=$$( dirname $$file ) ; \
		wire $$path ; \
	done

.PHONY: error
# 生成错误文件代码
error:
	@echo '生成错误文件代码...'
	@protoc --proto_path=. \
               --proto_path=$(THIRD_PARTY_PROTO_DIR) \
               --go_out=paths=source_relative:. \
               --go-errors_out=paths=source_relative:. \
               $(ERROR_PROTO_FILES)

.PHONY: config
# 生成配置文件代码
config:
	@echo '生成配置文件代码...'
	@for file in $(INTERNAL_PROTO_FILES) ; do \
		protoc --proto_path=. \
                --proto_path=$(THIRD_PARTY_PROTO_DIR) \
                --proto_path=$(PKG_CONF_DIR) \
                --go_out=paths=source_relative:. \
                $$file ; \
	done

.PHONY: api
# 生成api文件代码
api:
	@echo '生成api文件代码...'
	@protoc --proto_path=. \
	       --proto_path=$(THIRD_PARTY_PROTO_DIR) \
 	       --go_out=paths=source_relative:. \
 	       --go-http_out=paths=source_relative:. \
 	       --go-grpc_out=paths=source_relative:. \
 	       --validate_out=paths=source_relative,lang=go:. \
 	       --openapiv2_out . \
           --openapiv2_opt logtostderr=true \
           --openapiv2_opt json_names_for_fields=false \
	       $(API_PROTO_FILES)

.PHONY: build
# 构建
build:
	mkdir -p bin/ && go build -ldflags "-X main.Version=$(VERSION)" -o ./bin/ ./...

.PHONY: generate
# 代码生成
generate:
	go get github.com/google/wire/cmd/wire
	go install github.com/google/wire/cmd/wire@latest
	go generate ./...

.PHONY: all
# 生成所有代码
all:
	@make api;
	@make error;
	@make config;
	@make wire;

# 显示帮助
help:
	@echo ''
	@echo 'Usage:'
	@echo ' make [target]'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-\_0-9]+:/ { \
	helpMessage = match(lastLine, /^# (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")-1); \
			helpMessage = substr(lastLine, RSTART + 2, RLENGTH); \
			printf "\033[36m%-22s\033[0m %s\n", helpCommand,helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help
