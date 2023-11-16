GOPATH:=$(shell go env GOPATH)
VERSION=$(shell git describe --tags --always)
INTERNAL_PROTO_FILES=$(shell find internal -name *.proto)
API_PROTO_FILES=$(shell find api -name *.proto)
IMG_VERSION=latest

.PHONY: init
# init env
init:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/go-kratos/kratos/cmd/kratos/v2@latest
	go install github.com/go-kratos/kratos/cmd/protoc-gen-go-http/v2@latest
	go install github.com/go-kratos/kratos/cmd/protoc-gen-go-errors/v2@latest
	go install github.com/google/gnostic/cmd/protoc-gen-openapi@v0.6.1

.PHONY: errors
# generate errors code
errors:
	protoc --proto_path=. \
               --proto_path=./third_party \
               --go_out=paths=source_relative:. \
               --go-errors_out=paths=source_relative:. \
               $(API_PROTO_FILES)

.PHONY: config
# generate internal proto
config:
	protoc --proto_path=. \
	       --proto_path=./third_party \
 	       --go_out=paths=source_relative:. \
	       $(INTERNAL_PROTO_FILES)

.PHONY: configdev
# generate internal proto
configdev:
	protoc --proto_path=. \
	       --proto_path=../../../global-proto \
	       --proto_path=../../../global-proto/client \
	       --proto_path=../../../global-proto/client/account.proto \
 	       --go_out=paths=source_relative:. \
	       $(INTERNAL_PROTO_FILES)

update:
	go mod edit -replace xy3-proto=gitlab.com/firerocksg/xy3-proto/xy3-proto@develop

.PHONY: api
# generate api proto
api:
	protoc --proto_path=. \
	       --proto_path=./third_party \
 	       --go_out=paths=source_relative:. \
 	       --go-http_out=paths=source_relative:. \
 	       --go-grpc_out=paths=source_relative:. \
 	       --openapi_out==paths=source_relative:. \
	       $(API_PROTO_FILES)

.PHONY: build
# build
build:
	go mod tidy
	mkdir -p bin/ && go build -ldflags "-X main.Version=$(VERSION) -X main.Name="logger"" -o ./bin/ ./...

.PHONY: build-run-m1-test
build-run-m1-test:
	go mod edit -replace xy3-proto=gitlab.com/firerocksg/xy3-proto/xy3-proto@develop
	docker buildx build --platform linux/amd64 -t firerocksg-registry.ap-southeast-1.cr.aliyuncs.com/firerocksg/account-test:${IMG_VERSION} --push -f Dockerfile ..

# eg: make build-run-m1-prd IMG_VERSION=v1.0.x
.PHONY: build-run-m1-prd
build-run-m1-prd:
	go mod edit -replace xy3-proto=gitlab.com/firerocksg/xy3-proto/xy3-proto@develop
	docker buildx build --platform linux/amd64 -t firerocksg-registry.ap-southeast-1.cr.aliyuncs.com/firerocksg/account-prd:v1.1.4 --push -f Dockerfile ..

.PHONY: generate
# generate
generate:
	go generate ./...

.PHONY: all
# generate all
all:
	make api;
	make errors;
	make config;
	make generate;

# show help
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
