TEST?=$$(go list ./... | grep -v 'vendor')
TESTNAME?=
NAME?=buddy
BINARY?=terraform-provider-${NAME}
TF_LOG?=INFO
GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)
BUDDY_TOKEN?=1234567890
BUDDY_BASE_URL?=https://api.buddy.works
BUDDY_INSECURE?=false
BUDDY_GET_TOKEN?=curl
BUDDY_GH_PROJECT?=
BUDDY_GH_TOKEN?=

default: build

build: fmt
	go build -o ./bin/${BINARY}

test_dev:
	$(eval BUDDY_TOKEN=$(shell sh -c "${BUDDY_GET_TOKEN}"))
	go clean -testcache
	TF_ACC=1 TF_LOG=${TF_LOG} BUDDY_TOKEN=${BUDDY_TOKEN} BUDDY_GH_PROJECT=${BUDDY_GH_PROJECT} BUDDY_GH_TOKEN=${BUDDY_GH_TOKEN} BUDDY_BASE_URL=https://api.dev.io BUDDY_INSECURE=true go test $(TEST) -v ${TESTNAME} -timeout 60m

test:
	go clean -testcache
	TF_ACC=1 TF_LOG=${TF_LOG} BUDDY_TOKEN=${BUDDY_TOKEN} BUDDY_GH_PROJECT=${BUDDY_GH_PROJECT} BUDDY_GH_TOKEN=${BUDDY_GH_TOKEN} BUDDY_BASE_URL=${BUDDY_BASE_URL} BUDDY_INSECURE=${BUDDY_INSECURE} go test $(TEST) -v ${TESTNAME} -timeout 60m

fmt:
	gofmt -w $(GOFMT_FILES)

lint: fmt tfprovider golangci

docs:
	go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs

golangci:
	go run github.com/golangci/golangci-lint/cmd/golangci-lint run ./...

tfprovider:
	go run github.com/bflad/tfproviderlint/cmd/tfproviderlintx \
	-XS001=false \
	-XS002=false \
	-AT003=false \
	-XAT001=false \
	-V012=false \
	-R018=false \
 	./...

.PHONY: default build test_dev test fmt lint docs golangci tfprovider