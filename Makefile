default: all

GO_PACKAGES = $$(go list ./... | grep -v vendor)
GO_FILES = $$(find . -name "*.go" | grep -v vendor | uniq)

uaa-crud-linux:
	GOOS=linux GOARCH=amd64 go build -o uaa-crud.linux ./main.go

uaa-crud-darwin:
	GOOS=darwin GOARCH=amd64 go build -o uaa-crud.darwin ./main.go

build: uaa-crud-darwin uaa-crud-linux

unit-test:
	@go test ${GO_PACKAGES}

fmt:
	goimports -l -w $(GO_FILES)

vet:
	@go vet ${GO_PACKAGES}

test: generate unit-test vet

generate:
	go generate ./...

cleandep:
	rm -rf vendor
	rm -f Gopkg.lock

HAS_DEP := $(shell command -v dep;)
HAS_GOIMPORTS := $(shell command -v dep;)

boostrap: bootstrap

.PHONY: bootstrap
bootstrap:
ifndef HAS_GOIMPORTS
	go get -u golang.org/x/tools/cmd/goimports
endif
	go mod vendor


all: fmt test build

